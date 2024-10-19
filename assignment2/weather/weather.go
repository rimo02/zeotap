package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherAPI struct {
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_Like float64 `json:"feels_like"`
	} `json:"main"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Dt int64 `json:"dt"`
}

func ConvertTemp(temp float64, unit string) float64 {

	switch unit {
	case "C":
		return temp - 273.15
	case "F":
		return (temp-273.15)*9/5 + 32
	default:
		return temp - 273.15
	}
}

// fetching the weather from API
func FetchWeather(city, API_KEY string) (WeatherAPI, error) {
	var data WeatherAPI
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, API_KEY)
	resp, err := http.Get(url)
	if err != nil {
		return WeatherAPI{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}
	return data, nil
}
