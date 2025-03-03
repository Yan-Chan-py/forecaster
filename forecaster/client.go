package forecaster

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Yan-Chan-py/forecaster/config"
)

type WeatherClient struct {
	client  HTTPClient
	baseURL string
	key     string
}

func NewWeatherClient(cfg *config.APIconfig) (*WeatherClient, error) {
	if cfg == nil {
		slog.Error("invalid config")
		return nil, errors.New("invalid config")
	}
	return &WeatherClient{
		client:  &http.Client{},
		baseURL: cfg.BaseUrl,
		key:     cfg.APIkey,
	}, nil
}

func (c *WeatherClient) GetWeather(lat float64, lon float64,unit Unit) (*WeatherResponse, *ErrorResponse) {
	ctx := context.Background()
    fmt.Println(unit)
	requestUrl := fmt.Sprintf("%s?lat=%f&lon=%f&units=%s&appid=%s", c.baseURL, lat, lon,unit ,c.key)
	fmt.Println(requestUrl)
	req, err := http.NewRequest("GET", requestUrl, nil)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(200)*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctxWithTimeout)
	if err != nil {
		slog.Error("cannot create request")
	}
    weather := &WeatherResponse{} 
	errResp := c.doRequest(req, weather)
	if errResp != nil {
		return nil, errResp
	}
    if weather.Weather == nil && errResp == nil  {
        return nil, &ErrorResponse{
            Cod:"511",
            Message: " unknown error",
        }

    }
    return weather, errResp
}

func (c *WeatherClient) doRequest(req *http.Request, w *WeatherResponse) *ErrorResponse {
	res, err := c.client.Do(req)
    if err != nil {
        return &ErrorResponse {
            Message:"cannot proceed request",
            Cod:"502",
        }
    }
    select {
    case <-req.Context().Done():
    		return &ErrorResponse{
			Message: "cannot do request, deadline expired",
			Cod:     "504",
		}
    default: 
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
