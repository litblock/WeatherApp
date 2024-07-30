package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Region  string `json:"region"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		AirQuality struct {
			CO           float64 `json:"co"`
			NO2          float64 `json:"no2"`
			O3           float64 `json:"o3"`
			SO2          float64 `json:"so2"`
			PM2_5        float64 `json:"pm2_5"`
			PM10         float64 `json:"pm10"`
			USEPAIndex   int     `json:"us-epa-index"`
			GBDEFRAIndex int     `json:"gb-defra-index"`
		} `json:"air_quality"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date  string `json:"date"`
			Astro struct {
				Sunrise string `json:"sunrise"`
				Sunset  string `json:"sunset"`
			} `json:"astro"`
			Day struct {
				MaxtempC  float64 `json:"maxtemp_c"`
				MintempC  float64 `json:"mintemp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
	Alerts struct {
		Alert []struct {
			Headline  string `json:"headline"`
			Desc      string `json:"desc"`
			Effective string `json:"effective"`
			Expires   string `json:"expires"`
		} `json:"alert"`
	} `json:"alerts"`
}

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	weatherKey := os.Getenv("WEATHER_TOKEN")
	if weatherKey == "" {
		log.Fatal("WEATHER_TOKEN not set in .env file")
	}

	location := promptUser("Enter location: ")
	days := promptUser("Enter number of forecast days (1-10, default 1): ")
	aqi := promptUser("Include air quality data? (yes/no, default no): ")
	alerts := promptUser("Include weather alerts? (yes/no, default no): ")

	params := url.Values{}
	params.Add("key", weatherKey)
	params.Add("q", location)

	numDays := 1
	if days != "" {
		if d, err := strconv.Atoi(days); err == nil && d >= 1 && d <= 10 {
			numDays = d
			params.Add("days", days)
		} else {
			fmt.Println("Invalid number of days. Using default (1).")
		}
	}

	includeAQI := false
	if strings.ToLower(aqi) == "yes" {
		params.Add("aqi", "yes")
		includeAQI = true
	}

	includeAlerts := false
	if strings.ToLower(alerts) == "yes" {
		params.Add("alerts", "yes")
		includeAlerts = true
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?%s", params.Encode())

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("HTTP request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	fmt.Printf("Weather for %s, %s, %s\n\n", weather.Location.Name, weather.Location.Region, weather.Location.Country)

	fmt.Println("Current weather:")
	fmt.Printf("Temperature: %.2f°C\n", weather.Current.TempC)
	fmt.Printf("Condition: %s\n", weather.Current.Condition.Text)

	if includeAQI {
		fmt.Println("\nAir Quality:")
		fmt.Printf("US EPA Index: %d\n", weather.Current.AirQuality.USEPAIndex)
		fmt.Printf("UK DEFRA Index: %d\n", weather.Current.AirQuality.GBDEFRAIndex)
		fmt.Printf("CO: %.2f, NO2: %.2f, O3: %.2f, SO2: %.2f, PM2.5: %.2f, PM10: %.2f\n",
			weather.Current.AirQuality.CO, weather.Current.AirQuality.NO2,
			weather.Current.AirQuality.O3, weather.Current.AirQuality.SO2,
			weather.Current.AirQuality.PM2_5, weather.Current.AirQuality.PM10)
	}

	if numDays > 0 && len(weather.Forecast.Forecastday) > 0 {
		fmt.Println("\nForecast:")
		for _, forecast := range weather.Forecast.Forecastday {
			date, _ := time.Parse("2006-01-02", forecast.Date)
			fmt.Printf("%s:\n", date.Format("Monday, January 2"))
			fmt.Printf("  Max: %.2f°C, Min: %.2f°C\n", forecast.Day.MaxtempC, forecast.Day.MintempC)
			fmt.Printf("  Condition: %s\n", forecast.Day.Condition.Text)
			fmt.Printf("  Sunrise: %s\n", forecast.Astro.Sunrise)
			fmt.Printf("  Sunset: %s\n", forecast.Astro.Sunset)
		}
	}

	if includeAlerts && len(weather.Alerts.Alert) > 0 {
		fmt.Println("\nWeather Alerts:")
		for _, alert := range weather.Alerts.Alert {
			fmt.Printf("Headline: %s\n", alert.Headline)
			fmt.Printf("Description: %s\n", alert.Desc)
			fmt.Printf("Effective: %s\n", alert.Effective)
			fmt.Printf("Expires: %s\n\n", alert.Expires)
		}
	}
}
