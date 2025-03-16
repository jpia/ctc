package models

import (
	"ctc/logger"
	"sync"
	"time"
)

type WeatherStatus struct {
	Status      int
	DateChecked time.Time
}

var (
	weatherStatusSingleton *WeatherStatus
	once                   sync.Once
	weatherCodeLookup      = map[int]string{
		1000: "Clear",
		1003: "Partly cloudy",
		1006: "Cloudy",
		1009: "Overcast",
		1030: "Mist",
		1063: "Patchy rain possible",
		1066: "Patchy snow possible",
		1069: "Patchy sleet possible",
		1072: "Patchy freezing drizzle possible",
		1087: "Thundery outbreaks possible",
		1114: "Blowing snow",
		1117: "Blizzard",
		1135: "Fog",
		1147: "Freezing fog",
		1150: "Patchy light drizzle",
		1153: "Light drizzle",
		1168: "Freezing drizzle",
		1171: "Heavy freezing drizzle",
		1180: "Patchy light rain",
		1183: "Light rain",
		1186: "Moderate rain at times",
		1189: "Moderate rain",
		1192: "Heavy rain at times",
		1195: "Heavy rain",
		1198: "Light freezing rain",
		1201: "Moderate or heavy freezing rain",
		1204: "Light sleet",
		1207: "Moderate or heavy sleet",
		1210: "Patchy light snow",
		1213: "Light snow",
		1216: "Patchy moderate snow",
		1219: "Moderate snow",
		1222: "Patchy heavy snow",
		1225: "Heavy snow",
		1237: "Ice pellets",
		1240: "Light rain shower",
		1243: "Moderate or heavy rain shower",
		1246: "Torrential rain shower",
		1249: "Light sleet showers",
		1252: "Moderate or heavy sleet showers",
		1255: "Light snow showers",
		1258: "Moderate or heavy snow showers",
		1261: "Light showers of ice pellets",
		1264: "Moderate or heavy showers of ice pellets",
		1273: "Patchy light rain with thunder",
		1276: "Moderate or heavy rain with thunder",
		1279: "Patchy light snow with thunder",
		1282: "Moderate or heavy snow with thunder",
		500:  "API_SICK_DAY",
	}
)

func GetWeatherStatusInstance() *WeatherStatus {
	once.Do(func() {
		weatherStatusSingleton = &WeatherStatus{}
	})
	return weatherStatusSingleton
}

// IsCheckedToday returns true if DateChecked is equal to today's date (ignoring time)
func (ws *WeatherStatus) IsCheckedToday() bool {
	const layout = "2006-01-02"
	today := time.Now().Format(layout)
	dateChecked := ws.DateChecked.Format(layout)

	logger.DebugLog("Today's date: %s", today)
	logger.DebugLog("DateChecked: %s", dateChecked)

	return dateChecked == today
}

// GetWeatherLabel returns the weather label for a given weather code
func GetWeatherLabel(code int) string {
	if label, exists := weatherCodeLookup[code]; exists {
		return label
	}
	return "Unknown"
}

// IsValidForStandardRelease returns true if the weather code is valid for standard release
func IsValidForStandardRelease(code int) bool {
	return code == 1000
}

// IsValidForApiSickDayRelease returns true if the weather code is valid for API sick day
func IsValidForApiSickDayRelease(code int) bool {
	return code == 500
}
