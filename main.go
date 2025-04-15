package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chrisk1905/CodingCanalWeather/internal/database"
	"github.com/Chrisk1905/CodingCanalWeather/internal/services"
	_ "github.com/lib/pq"
)

func main() {
	//open connection to postgresDB
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL enviroment variable is empty")
	}
	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %+v", err)
	}
	defer connection.Close()

	app := &services.WeatherService{
		Repo:       database.New(connection),
		Client:     http.DefaultClient,
		WeatherAPI: os.Getenv("WEATHER_API_KEY"),
		WeatherURL: "https://api.openweathermap.org/data/2.5/weather",
	}

	//periodic querying of API
	var cities []services.City = services.GetCities()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for tick := range ticker.C {
		fmt.Println(tick)
		for _, city := range cities {
			weather, err := app.GetWeather(city)
			if err != nil {
				fmt.Printf("error while in getWeather: %+v\n", err)
				continue
			}
			err = app.StoreWeather(weather)
			if err != nil {
				fmt.Printf("error storing weather: %+v\n", err)
			}
		}
	}
}
