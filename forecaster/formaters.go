package forecaster
import (
    "fmt"
)
func (w WeatherResponse) MainInfo() string {
    return fmt.Sprintf(
        "Current temperature: %.2f°C\n"+
        "Feels like: %.2f°C\n"+
        "Min temperature today: %.2f°C\n"+
        "Max temperature today: %.2f°C",
        w.Main.Temp, w.Main.FeelsLike, w.Main.TempMin, w.Main.TempMax,
    )
}





