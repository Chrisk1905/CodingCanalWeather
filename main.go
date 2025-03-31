package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type city struct {
	name string
	lat  float32
	lon  float32
}

type weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
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

func getWeather(city city) error {
	baseUrl := "https://api.openweathermap.org/data/2.5/weather"
	apiKey := os.Getenv("WEATHER_API_KEY")
	urlToCall, _ := url.Parse(baseUrl)
	// Define query parameters
	params := url.Values{}
	params.Add("lat", fmt.Sprintf("%f", city.lat))
	params.Add("lon", fmt.Sprintf("%f", city.lon))
	params.Add("appid", apiKey)
	urlToCall.RawQuery = params.Encode()

	//make get call, and unmarshall JSON
	weather := weather{}
	res, err := http.Get(urlToCall.String())
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return err
	}
	prettyJSON, err := json.MarshalIndent(weather, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(city.name)
	fmt.Println(string(prettyJSON))
	return nil
}

func getCities() []city {
	Seattle := city{name: "Seattle", lat: 47.608013, lon: -122.335167}
	LA := city{name: "Los Angeles", lat: 34.052235, lon: -118.243683}
	NewYork := city{name: "New York", lat: 40.730610, lon: -73.935242}
	Seoul := city{name: "Seoul", lat: 37.532600, lon: 127.024612}
	Vancouver := city{name: "Vancouver", lat: 49.246292, lon: -123.116226}
	return []city{Seattle, LA, NewYork, Seoul, Vancouver}
}

func main() {
	var cities []city = getCities()

	for _, city := range cities {
		if err := getWeather(city); err != nil {
			fmt.Println(err)
		}
	}

}
