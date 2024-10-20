package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rimo02/zeotap/assignment2/config"
	"github.com/rimo02/zeotap/assignment2/database"
	"github.com/rimo02/zeotap/assignment2/weather"
	"log"
	"net/http"
	"os"
)

func init() {
	database.InitializeConnections()
	err1 := godotenv.Load()

	if err1 != nil {
		panic(err1)
	}

	_, err := weather.FetchWeather("London", os.Getenv("API_KEY"))
	if err != nil {
		log.Printf("Error fetching weather data: %v\n", err)
		log.Fatal("Unable to connect to OpenWeatherMap API. Stopping the program.\n")
		os.Exit(1)
	} else {
		fmt.Printf("Successfully connected to OpenWeatherMap API\n")
	}
}
func main() {
	cfg := config.LoadConfig()
	go weather.FetchWeatherData(cfg)

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(weather.AllWeatherData)
		if err != nil {
			http.Error(w, "Error marshalling weather data", http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
	})

	if err := http.ListenAndServe(":8070", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
