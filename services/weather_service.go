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
		return 0, fmt.Errorf("Weather API key not set")
	}

	encodedLocation := url.QueryEscape(location)

	var urlStr string
	if date.Before(time.Now()) {
		urlStr = fmt.Sprintf("https://api.weatherapi.com/v1/history.json?key=%s&q=%s&dt=%s", apiKey, encodedLocation, date.Format("2006-01-02"))
	} else {
		urlStr = fmt.Sprintf("https://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&dt=%s", apiKey, encodedLocation, date.Format("2006-01-02"))
	}

	fmt.Printf("Checking weather for %s at %s...\n", location, date.Format("2006-01-02"))
	fmt.Printf("URL: %s\n", urlStr)

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

	return weatherResponse.Forecast.Forecastday[0].Day.Condition.Code, nil
}

func UpdateStoreByWeather() {
	fmt.Println("Checking weather for pending CTCs...")
	for shortcode, ctc := range models.CTCStore {
		if ctc.Status == models.Pending {
			weatherCode, err := checkWeather(ctc.ReleaseDate, "New York City")
			if err != nil {
				fmt.Printf("Error checking weather for shortcode %s: %v\n", shortcode, err)
				continue
			}

			if weatherCode == 1000 {
				ctc.Status = models.Ready
				models.CTCStore[shortcode] = ctc
				fmt.Printf("Shortcode %s is now ready due to clear weather.\n", shortcode)
			}
		}
	}
}
