package weather

import (
	"context"
	"fmt"
	"github.com/rimo02/zeotap/assignment2/config"
	"github.com/rimo02/zeotap/assignment2/database"
	gomail "gopkg.in/mail.v2"

	"go.mongodb.org/mongo-driver/bson"
)

func TriggerAlert(city string, data WeatherAPI, threshold config.Threshold) {
	symbol := "⚠️"
	fmt.Printf("%s ALERT: Weather Threshold Breached in city: %s %s", symbol, city, symbol)
	fmt.Printf("Current Temperature: %.2f°C (Threshold: %.2f°C)\n", data.Main.Temp, threshold.MaximumTemp)
	fmt.Printf("Weather Condition: %s\n", data.Weather[0].Main)
	fmt.Printf("%s Please take necessary precautions! %s\n", symbol, symbol)

	collection := database.GetUserCollection(database.UserClient, city)
	cursor, err := collection.Find(context.TODO(), bson.M{"email": bson.M{"$exists": true}})
	if err != nil {
		fmt.Println("Error fetching users:", err.Error())
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user struct {
			Email string `bson:"email"`
		}
		if err := cursor.Decode(&user); err != nil {
			fmt.Println("Error decoding user data:", err.Error())
			continue
		}
		message := gomail.NewMessage()
		message.SetHeader("From", "weatherdepartment@gmail.co.in")
		message.SetHeader("To", user.Email)
		message.SetHeader("Subject", fmt.Sprintf("Weather Alert for %s", city))
		body := fmt.Sprintf("%s ALERT: Weather Threshold Breached in city: %s %s", symbol, city, symbol)
		body += fmt.Sprintf("Current Temperature: %.2f°C (Threshold: %.2f°C)\n", data.Main.Temp, threshold.MaximumTemp)
		body += fmt.Sprintf("Weather Condition: %s\n", data.Weather[0].Main)
		body += fmt.Sprintf("%s Please take necessary precautions! %s\n", symbol, symbol)
		message.SetBody("text/plain", body)

		dialer := gomail.NewDialer("live.smtp", 587, "api", "abcd12456")
		if err := dialer.DialAndSend(message); err != nil {
			fmt.Println("Error", err.Error())
		} else {
			fmt.Println("Email sent successfully!")
		}
	}
	if err := cursor.Err(); err != nil {
		fmt.Println("Error iterating through users:", err.Error())
	}
}