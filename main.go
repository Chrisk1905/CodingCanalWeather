package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Chrisk1905/CodingCanalWeather/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type App struct {
	dbQueries *database.Queries
}

type City struct {
	name string
	lat  float32
	lon  float32
}

type WeatherCondition struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Conditions []WeatherCondition `json:"weather"`
	Base       string             `json:"base"`
	Main       struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func getWeather(city City) (Weather, error) {
	baseUrl := "https://api.openweathermap.org/data/2.5/weather"
	apiKey := os.Getenv("WEATHER_API_KEY")
	urlToCall, _ := url.Parse(baseUrl)
	// Define query parameters
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", city.lat))
	params.Add("lon", fmt.Sprintf("%f", city.lon))
	params.Add("appid", apiKey)
	urlToCall.RawQuery = params.Encode()

	//make GET request
	weather := Weather{}
	res, err := http.Get(urlToCall.String())
	if err != nil {
		return weather, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return weather, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return weather, err
	}
	//Unmarshal JSON
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return weather, err
	}

	prettyWeatherJSON, err := json.MarshalIndent(weather, "", "  ")
	if err != nil {
		fmt.Printf("error printing weather: %+v", err)
	}
	fmt.Println(city.name)
	fmt.Println(string(prettyWeatherJSON))

	return weather, nil
}

func getCities() []City {
	Seattle := City{name: "Seattle", lat: 47.608013, lon: -122.335167}
	LA := City{name: "Los Angeles", lat: 34.052235, lon: -118.243683}
	NewYork := City{name: "New York", lat: 40.730610, lon: -73.935242}
	Seoul := City{name: "Seoul", lat: 37.532600, lon: 127.024612}
	Vancouver := City{name: "Vancouver", lat: 49.246292, lon: -123.116226}
	return []City{Seattle, LA, NewYork, Seoul, Vancouver}
}

func (app *App) storeWeather(w Weather) error {
	//check if the coordinate is in the database
	ctx := context.Background()
	getCoordParams := database.GetWeatherCoordParams{
		Lon: sql.NullFloat64{Float64: w.Coord.Lon, Valid: true},
		Lat: sql.NullFloat64{Float64: w.Coord.Lat, Valid: true},
	}
	coord, err := app.dbQueries.GetWeatherCoord(ctx, getCoordParams)
	//insert the coordinate if not in the database
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		coordParams := database.InsertWeatherCoordinatesParams{
			ID:  uuid.New(),
			Lon: sql.NullFloat64{Float64: w.Coord.Lon, Valid: true},
			Lat: sql.NullFloat64{Float64: w.Coord.Lat, Valid: true},
		}
		coord, err = app.dbQueries.InsertWeatherCoordinates(ctx, coordParams)
		if err != nil {
			return err
		}
	}
	//check if WeatherCondition is in the database
	conditions := make([]database.WeatherCondition, 0, len(w.Conditions))
	for _, c := range w.Conditions {
		condID := sql.NullInt32{Int32: int32(c.ID), Valid: true}
		dbCondition, err := app.dbQueries.GetConditionByConditionID(ctx, condID)

		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			// Not found in database - insert new condition
			newCond := database.InsertConditionParams{
				ID:          uuid.New(),
				ConditionID: condID,
				Main:        sql.NullString{String: c.Main, Valid: c.Main != ""},
				Description: sql.NullString{String: c.Description, Valid: c.Description != ""},
				Icon:        sql.NullString{String: c.Icon, Valid: c.Icon != ""},
			}

			dbCondition, err = app.dbQueries.InsertCondition(ctx, newCond)
			if err != nil {
				return err
			}
		}
		conditions = append(conditions, dbCondition)
	}
	//insert WeatherDatum
	insertDatumParams := database.InsertWeatherDatumParams{
		ID:            uuid.New(),
		CoordinatesID: uuid.NullUUID{UUID: coord.ID, Valid: true},
		CityName:      sql.NullString{String: w.Name, Valid: w.Name != ""},
		Country:       sql.NullString{String: w.Sys.Country, Valid: w.Sys.Country != ""},
		Temperature:   sql.NullFloat64{Float64: w.Main.Temp, Valid: true},
		FeelsLike:     sql.NullFloat64{Float64: w.Main.FeelsLike, Valid: true},
		TempMin:       sql.NullFloat64{Float64: w.Main.TempMin, Valid: true},
		TempMax:       sql.NullFloat64{Float64: w.Main.TempMax, Valid: true},
		Pressure:      sql.NullInt32{Int32: int32(w.Main.Pressure), Valid: true},
		Humidity:      sql.NullInt32{Int32: int32(w.Main.Humidity), Valid: true},
		SeaLevel:      sql.NullInt32{Int32: int32(w.Main.SeaLevel), Valid: true},
		GrndLevel:     sql.NullInt32{Int32: int32(w.Main.GrndLevel), Valid: true},
		Visibility:    sql.NullInt32{Int32: int32(w.Visibility), Valid: true},
		WindSpeed:     sql.NullFloat64{Float64: w.Wind.Speed, Valid: true},
		WindDeg:       sql.NullInt32{Int32: int32(w.Wind.Deg), Valid: true},
		Cloudiness:    sql.NullInt32{Int32: int32(w.Clouds.All), Valid: true},
		Timestamp:     sql.NullTime{Time: time.Unix(int64(w.Dt), 0), Valid: true},
		Sunrise:       sql.NullTime{Time: time.Unix(int64(w.Sys.Sunrise), 0), Valid: true},
		Sunset:        sql.NullTime{Time: time.Unix(int64(w.Sys.Sunset), 0), Valid: true},
		Timezone:      sql.NullInt32{Int32: int32(w.Timezone), Valid: true},
	}
	datum, err := app.dbQueries.InsertWeatherDatum(ctx, insertDatumParams)
	if err != nil {
		return err
	}
	//insert Weather_data_conditions
	for _, condition := range conditions {

		insertDataConditionsParams := database.InsertWeatherDataConditionsParams{
			WeatherDataID: datum.ID,
			ConditionID:   condition.ConditionID.Int32,
		}

		_, err := app.dbQueries.InsertWeatherDataConditions(ctx, insertDataConditionsParams)
		if err != nil {
			return err
		}
	}

	return nil
}

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

	app := &App{
		dbQueries: database.New(connection),
	}

	//periodic querying of API
	var cities []City = getCities()
	ticker := time.NewTicker(time.Second * 6)
	defer ticker.Stop()

	for tick := range ticker.C {
		fmt.Println(tick)
		for _, city := range cities {
			weather, err := getWeather(city)
			if err != nil {
				fmt.Printf("error while in getWeather: %+v\n", err)
				continue
			}
			err = app.storeWeather(weather)
			if err != nil {
				fmt.Printf("error storing weather: %+v\n", err)
			}
		}
	}
}
