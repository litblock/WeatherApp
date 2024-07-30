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
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempC  float64 `json:"maxtemp_c"`
				MintempC  float64 `json:"mintemp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
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

	if strings.ToLower(aqi) == "yes" {
		params.Add("aqi", "yes")
	}

	if strings.ToLower(alerts) == "yes" {
		params.Add("alerts", "yes")
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
	fmt.Printf("Condition: %s\n\n", weather.Current.Condition.Text)

	if numDays > 1 {
		fmt.Println("Forecast:")
		for _, forecast := range weather.Forecast.Forecastday {
			date, _ := time.Parse("2006-01-02", forecast.Date)
			fmt.Printf("%s:\n", date.Format("Monday, January 2"))
			fmt.Printf("  Max: %.2f°C, Min: %.2f°C\n", forecast.Day.MaxtempC, forecast.Day.MintempC)
			fmt.Printf("  Condition: %s\n\n", forecast.Day.Condition.Text)
		}
	}
}
