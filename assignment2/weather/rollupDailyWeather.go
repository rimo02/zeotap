package weather

import (
	"context"
	"fmt"
	"github.com/rimo02/zeotap/assignment2/database"
	"math"
	"time"
)

// happens at midnight. It takes the entire data collected for the day and then performs aggregation functions 
// like maximum temperature for the day or avg temprature and dominant condition for the entire day
func RollUpDailyWeatherData() {
	Mu.Lock()
	defer Mu.Unlock()
	for city, summaries := range DailyWeatherData {
		var mini, maxi, avg float64
		condn := make(map[string]int)
		for i, summary := range summaries {
			if i == 0 {
				mini = summary.MinTemp
				maxi = summary.MaxTemp
			} else {
				mini = math.Min(mini, summary.MinTemp)
				maxi = math.Max(maxi, summary.MaxTemp)
			}
			avg += summary.AvgTemp
			condn[summary.DominantCond]++
		}
		avg /= float64(len(summaries))
		dominantCond := ""
		maxCount := 0
		for condition, count := range condn {
			if count > maxCount {
				dominantCond = condition
				maxCount = count
			}
		}

		dailySummary := &Summary{
			City:         city,
			MinTemp:      mini,
			MaxTemp:      maxi,
			Date:         time.Now(),
			AvgTemp:      avg,
			DominantCond: dominantCond,
		}
		collection := database.GetCollection(database.WeatherClient, city+"_daily")
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()
		_, err := collection.InsertOne(ctx, dailySummary)
		if err != nil {
			fmt.Printf("Error storing daily summary for %s: %v\n", city, err)
		}
		DailyWeatherData = make(map[string][]Summary) // clean it after daily rollup
	}
}
