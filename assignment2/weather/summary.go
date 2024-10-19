package weather

import (
	"fmt"
	"github.com/rimo02/zeotap/assignment2/config"
	"os"
	"strings"
	"sync"
	"time"
)

type Summary struct {
	City         string
	Date         time.Time
	AvgTemp      float64
	MaxTemp      float64
	MinTemp      float64
	DominantCond string
}

var (
	AllWeatherData   []interface{}                // contains fields to be displayed ebery 5 mins
	DailyWeatherData = make(map[string][]Summary) // city: summary
	AlertCounts      = make(map[string]int)       // map to store how many times city:threshold breached
	Mu               sync.RWMutex
)

func CheckThreshold(city string, data WeatherAPI, threshold config.Threshold) {
	Mu.Lock()
	defer Mu.Unlock()
	if data.Main.Temp > threshold.MaximumTemp {
		AlertCounts[city]++
		if AlertCounts[city] >= threshold.Breach {
			TriggerAlert(city, data, threshold)
		}
	}
}

func FetchWeatherData(cfg config.Config) {
	for {
		time.Sleep(cfg.Interval)
		tempData := make([]interface{}, len(cfg.Cities))

		for i, city := range cfg.Cities {
			data, err := FetchWeather(strings.ToLower(city.Name), os.Getenv("apiKey"))
			if err != nil {
				fmt.Printf("Error fetching weather for %s: %v\n", city.Name, err)
				continue
			}
			data.Main.Temp = ConvertTemp(data.Main.Temp, cfg.TempUnit)
			tempData[i] = data
			UpdateDailyWeatherdata(city.Name, data)

			CheckThreshold(city.Name, data, city.TempThreshold)
		}

		Mu.Lock()
		AllWeatherData = tempData
		Mu.Unlock()

		// if time.Now().Hour() == 0 && time.Now().Minute() == 0 {
		// }
		RollUpDailyWeatherData()
	}
}

func UpdateDailyWeatherdata(city string, data WeatherAPI) {
	Mu.Lock()
	defer Mu.Unlock()
	summary := Summary{
		City:         city,
		Date:         time.Now(),
		MaxTemp:      data.Main.Temp,
		MinTemp:      data.Main.Temp,
		AvgTemp:      data.Main.Temp,
		DominantCond: data.Weather[0].Main,
	}
	DailyWeatherData[city] = append(DailyWeatherData[city], summary)
}
