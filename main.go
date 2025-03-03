package main

import (
	"fmt"
	"github.com/Yan-Chan-py/forecaster/config"
	"github.com/Yan-Chan-py/forecaster/forecaster"
	"log/slog"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
	client, err := forecaster.NewWeatherClient(cfg)
	if err != nil {
		slog.Error("cannot create client")
	}
	weather, errResp := client.GetWeather(31.84, 24.031,0)
	if errResp != nil {
		fmt.Println(errResp)
	}
	fmt.Printf(weather.MainInfo())

}
