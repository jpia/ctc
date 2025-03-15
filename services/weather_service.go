package services

import (
	"ctc/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WeatherResponse struct {
	Forecast struct {
		Forecastday []struct {
			Day struct {
				Condition struct {
					Code int `json:"code"`
				} `json:"condition"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

type WeatherErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func checkWeather(date time.Time, location string) (int, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("weather API key not set")
	}

	encodedLocation := url.QueryEscape(location)

	var urlStr string
	if date.Before(time.Now()) {
		urlStr = fmt.Sprintf("https://api.weatherapi.com/v1/history.json?key=%s&q=%s&dt=%s", apiKey, encodedLocation, date.Format("2006-01-02"))
	} else {
		urlStr = fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&dt=%s", apiKey, encodedLocation, date.Format("2006-01-02"))
	}

	DebugLog("checking weather for %s at %s using %s\n", location, date.Format("2006-01-02"), urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var weatherError WeatherErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&weatherError); err != nil {
			return 0, fmt.Errorf("failed to get weather data: %s", resp.Status)
		}
		return 0, fmt.Errorf("failed to get weather data: %s", weatherError.Error.Message)
	}

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return 0, err
	}

	if len(weatherResponse.Forecast.Forecastday) == 0 {
		return 0, fmt.Errorf("no forecast data available")
	}
	conditionCode := weatherResponse.Forecast.Forecastday[0].Day.Condition.Code
	DebugLog("weather condition code: %d\n", conditionCode)
	return conditionCode, nil
}

func UpdateStoreByWeather() {

	for shortcode, ctc := range models.URLStore {
		if ctc.Status == models.PendingStatus {
			weatherCode, err := checkWeather(ctc.ReleaseDate, "New York City")
			if err != nil {
				ErrorLog("Error checking weather for shortcode %s: %v\n", shortcode, err)
				continue
			}

			if weatherCode == 1000 {
				ctc.Status = models.ReleasedStatus
				models.URLStore[shortcode] = ctc
				DebugLog("Shortcode %s is now ready due to clear weather.\n", shortcode)
			} else {
				DebugLog("Shortcode %s is still pending due to weather code %d.\n", shortcode, weatherCode)
			}
		}
	}
}
