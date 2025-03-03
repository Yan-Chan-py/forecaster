package forecaster

import (
	"encoding/json"
	"fmt"
)
type Unit int
type Lang string

const (
    Metric = iota
    Standart
    Imperial
    Ukrainan = "ua"

)

type WeatherResponse struct {
	Main     Main      `json:"main"`
	Wind     Wind      `json:"wind"`
	Clouds   Clouds    `json:"clouds"`
	Weather  []Weather `json:"weather"`
	Coord    Coord     `json:"coord"`
	Timezone int       `json:"timezone"`
	Name     string    `json:"name"`
}
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Define struct for Clouds
type Clouds struct {
	All int `json:"all"`
}
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type ErrorResponse struct {
	Cod     string `json:"cod"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}
func (e *ErrorResponse) UnmarshalJSON(data []byte) error {
	// Create a temporary structure to hold raw data
	var raw struct {
		Cod     any    `json:"cod"` // Can be int or string
		Message string `json:"message"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert "cod" field to string format
	switch v := raw.Cod.(type) {
	case float64: // JSON numbers are parsed as float64
		e.Cod = fmt.Sprintf("%d", int(v))
	case string:
		e.Cod = v
	default:
		e.Cod = "unknown"
	}

	e.Message = raw.Message
	return nil
}

func (u Unit)String() string {
    return [...]string{"metric","standart","imperial"}[u]
}
type RequestOptions struct {
    unit Unit
    Language Lang

}
