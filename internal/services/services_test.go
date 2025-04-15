package services

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetWeather_ValidResponse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{
		"coord": { "lon": -123.1162, "lat": 49.2463 },
		"weather": [ { "id": 803, "main": "Clouds", "description": "broken clouds", "icon": "04d" } ],
		"base": "stations",
		"main": {
			"temp": 289.13, "feels_like": 288.25, "temp_min": 285.62, "temp_max": 292.09,
			"pressure": 1020, "humidity": 56, "sea_level": 1020, "grnd_level": 1012
		},
		"visibility": 10000,
		"wind": { "speed": 7.15, "deg": 243 },
		"clouds": { "all": 75 },
		"dt": 1744671611,
		"sys": {
			"type": 2, "id": 2009577, "country": "CA",
			"sunrise": 1744636913, "sunset": 1744686177
		},
		"timezone": -25200, "id": 6173331, "name": "Vancouver", "cod": 200
		}`)
	}))
	defer mockServer.Close()

	service := &WeatherService{
		Client:     mockServer.Client(),
		WeatherAPI: "",
		WeatherURL: mockServer.URL,
	}

	city := City{Name: "Vancouver", Lat: 49.2463, Lon: -123.1162}
	weather, err := service.GetWeather(city)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if weather.Coord.Lon != -123.1162 {
		t.Errorf("Expected Coord.Lon to be -123.1162, got %f", weather.Coord.Lon)
	}
	if weather.Coord.Lat != 49.2463 {
		t.Errorf("Expected Coord.Lat to be 49.2463, got %f", weather.Coord.Lat)
	}

	if len(weather.Conditions) != 1 {
		t.Errorf("Expected 1 weather condition, got %d", len(weather.Conditions))
	} else {
		cond := weather.Conditions[0]
		if cond.ID != 803 {
			t.Errorf("Expected Weather[0].ID to be 803, got %d", cond.ID)
		}
		if cond.Main != "Clouds" {
			t.Errorf("Expected Weather[0].Main to be Clouds, got %s", cond.Main)
		}
		if cond.Description != "broken clouds" {
			t.Errorf("Expected Weather[0].Description to be 'broken clouds', got %s", cond.Description)
		}
		if cond.Icon != "04d" {
			t.Errorf("Expected Weather[0].Icon to be '04d', got %s", cond.Icon)
		}
	}

	if weather.Base != "stations" {
		t.Errorf("Expected Base to be 'stations', got %s", weather.Base)
	}

	if weather.Main.Temp != 289.13 {
		t.Errorf("Expected Main.Temp to be 289.13, got %f", weather.Main.Temp)
	}
	if weather.Main.FeelsLike != 288.25 {
		t.Errorf("Expected Main.FeelsLike to be 288.25, got %f", weather.Main.FeelsLike)
	}
	if weather.Main.TempMin != 285.62 {
		t.Errorf("Expected Main.TempMin to be 285.62, got %f", weather.Main.TempMin)
	}
	if weather.Main.TempMax != 292.09 {
		t.Errorf("Expected Main.TempMax to be 292.09, got %f", weather.Main.TempMax)
	}
	if weather.Main.Pressure != 1020 {
		t.Errorf("Expected Main.Pressure to be 1020, got %d", weather.Main.Pressure)
	}
	if weather.Main.Humidity != 56 {
		t.Errorf("Expected Main.Humidity to be 56, got %d", weather.Main.Humidity)
	}
	if weather.Main.SeaLevel != 1020 {
		t.Errorf("Expected Main.SeaLevel to be 1020, got %d", weather.Main.SeaLevel)
	}
	if weather.Main.GrndLevel != 1012 {
		t.Errorf("Expected Main.GrndLevel to be 1012, got %d", weather.Main.GrndLevel)
	}

	if weather.Visibility != 10000 {
		t.Errorf("Expected Visibility to be 10000, got %d", weather.Visibility)
	}

	if weather.Wind.Speed != 7.15 {
		t.Errorf("Expected Wind.Speed to be 7.15, got %f", weather.Wind.Speed)
	}
	if weather.Wind.Deg != 243 {
		t.Errorf("Expected Wind.Deg to be 243, got %d", weather.Wind.Deg)
	}

	if weather.Clouds.All != 75 {
		t.Errorf("Expected Clouds.All to be 75, got %d", weather.Clouds.All)
	}

	if weather.Dt != 1744671611 {
		t.Errorf("Expected Dt to be 1744671611, got %d", weather.Dt)
	}

	if weather.Sys.Type != 2 {
		t.Errorf("Expected Sys.Type to be 2, got %d", weather.Sys.Type)
	}
	if weather.Sys.ID != 2009577 {
		t.Errorf("Expected Sys.ID to be 2009577, got %d", weather.Sys.ID)
	}
	if weather.Sys.Country != "CA" {
		t.Errorf("Expected Sys.Country to be 'CA', got %s", weather.Sys.Country)
	}
	if weather.Sys.Sunrise != 1744636913 {
		t.Errorf("Expected Sys.Sunrise to be 1744636913, got %d", weather.Sys.Sunrise)
	}
	if weather.Sys.Sunset != 1744686177 {
		t.Errorf("Expected Sys.Sunset to be 1744686177, got %d", weather.Sys.Sunset)
	}

	if weather.Timezone != -25200 {
		t.Errorf("Expected Timezone to be -25200, got %d", weather.Timezone)
	}
	if weather.ID != 6173331 {
		t.Errorf("Expected ID to be 6173331, got %d", weather.ID)
	}
	if weather.Name != "Vancouver" {
		t.Errorf("Expected Name to be 'Vancouver', got %s", weather.Name)
	}
	if weather.Cod != 200 {
		t.Errorf("Expected Cod to be 200, got %d", weather.Cod)
	}

}

func TestGetWeather_StatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "400 Bad Request",
			statusCode:     http.StatusBadRequest,
			responseBody:   `{"cod":400, "message":"Bad request"}`,
			expectError:    true,
			expectedErrMsg: "response failed with status code: 400",
		},
		{
			name:           "403 Forbidden",
			statusCode:     http.StatusForbidden,
			responseBody:   `{"cod":403, "message":"Access denied"}`,
			expectError:    true,
			expectedErrMsg: "response failed with status code: 403",
		},
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"cod":404, "message":"City not found"}`,
			expectError:    true,
			expectedErrMsg: "response failed with status code: 404",
		},
		{
			name:           "429 Too Many Requests",
			statusCode:     http.StatusTooManyRequests,
			responseBody:   `{"cod":429, "message":"Rate limit exceeded"}`,
			expectError:    true,
			expectedErrMsg: "response failed with status code: 429",
		},
		{
			name:           "502 Bad Gateway",
			statusCode:     http.StatusBadGateway,
			responseBody:   `{"cod":502, "message":"Bad gateway"}`,
			expectError:    true,
			expectedErrMsg: "response failed with status code: 502",
		},
		{
			name:           "503 Service Unavailable",
			statusCode:     http.StatusServiceUnavailable,
			responseBody:   `{"cod":503, "message":"Service unavailable"}`,
			expectError:    true,
			expectedErrMsg: "response failed with status code: 503",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				fmt.Fprintln(w, tc.responseBody)
			}))
			defer mockServer.Close()

			service := &WeatherService{
				Client:     mockServer.Client(),
				WeatherAPI: "",
				WeatherURL: mockServer.URL,
			}

			city := City{Name: "Vancouver", Lat: 49.2463, Lon: -123.1162}
			_, err := service.GetWeather(city)

			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error for status code %d, got none", tc.statusCode)
				}
				if !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %v", tc.expectedErrMsg, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error for status code %d, got %v", tc.statusCode, err)
				}
			}
		})
	}
}
