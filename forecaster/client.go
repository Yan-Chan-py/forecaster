package forecaster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Yan-Chan-py/forecaster/config"
)

type WeatherClient struct {
	client  HTTPClient
	baseURL string
	key     string
    cfg *config.APIconfig
}

func NewWeatherClient(cfg *config.APIconfig) (*WeatherClient, error) {
	if cfg == nil {
		slog.Error("invalid config")
		return nil, errors.New("invalid config")
	}
	return &WeatherClient{
		client:  &http.Client{
            Timeout: cfg.Timeout,
        },
		baseURL: cfg.BaseUrl,
		key:     cfg.APIkey,
        cfg: cfg,
	}, nil
}

func (c *WeatherClient) GetWeather(lat float64, lon float64,opts *RequestOptions) (*WeatherResponse, *ErrorResponse) {
	ctx := context.Background()
    var requestUrl string
    if opts == nil {
	requestUrl = fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s", c.baseURL, lat, lon,c.key)
    } else {
	requestUrl = fmt.Sprintf("%s?lat=%f&lon=%f&units=%s&lang=%s&appid=%s", c.baseURL, lat, lon,opts.Unit,opts.Language ,c.key)
    }

	fmt.Println(requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
    if err != nil {
        slog.Error("cannot create request")
    }
    switch c.cfg.Timeout {
    case   0:
        req = req.WithContext(ctx)
    default:
        ctxWithTimeout,cancel := context.WithTimeout(ctx,c.cfg.Timeout)
        req = req.WithContext(ctxWithTimeout)
        defer cancel()
    }
	
    weather := &WeatherResponse{} 
	errResp := c.doRequest(req, weather)
	if errResp != nil {
		return nil, errResp
	}
    if len(weather.Weather) == 0 && errResp == nil  {
        return nil, &ErrorResponse{
            Cod:"511",
            Message: " unknown error",
        }

    }
    return weather, errResp
}

func (c *WeatherClient) doRequest(req *http.Request, w *WeatherResponse) *ErrorResponse {
	res, err := c.client.Do(req)
    select {
    case <-req.Context().Done():
    		return &ErrorResponse{
			Message: "cannot do request, deadline expired",
			Cod:     "504",
		}
    default: 
        if err != nil {
        return &ErrorResponse {
            Message:"cannot proceed request",
            Cod:"502",
        }
    }
    	defer func() {
		_ = res.Body.Close()
	}()
    }
	fmt.Println("HTTP Status:", res.StatusCode)

	errorResponse := &ErrorResponse{}
	if res.StatusCode != http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(errorResponse); err != nil {
			fmt.Println("Failed to decode error response:", err)
			return &ErrorResponse{
				Cod:     "488",
				Message: "unknown error",
			}
		}
		return errorResponse
	}

	if w != nil { // Ensure w is not nil before decoding
		if err := json.NewDecoder(res.Body).Decode(w); err != nil {
			return &ErrorResponse{
				Cod:     "476",
				Message: "cannot decode response",
			}
		}
	}
	return nil
}


