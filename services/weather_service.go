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

func UpdateWeatherStatus() {
	weatherInstance := models.GetWeatherStatusInstance()

	if weatherInstance.IsCheckedToday() {
		DebugLog("Weather status already checked today, skipping update.")
		return
	}

	weatherCode, err := checkWeather(time.Now(), "New York City")
	if err != nil {
		ErrorLog("Error checking weather: %v\n", err)
		return
	}

	weatherInstance.Status = weatherCode
	weatherInstance.DateChecked = time.Now()
	DebugLog("Weather status updated: %d\n", weatherCode)
}
