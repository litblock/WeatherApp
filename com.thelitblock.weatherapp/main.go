package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	weatherKey := os.Getenv("WEATHER_TOKEN")
	if weatherKey == "" {
		log.Fatal("WEATHER_TOKEN not set in .env file")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter location: ")
	location, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading location: %v", err)
	}
	location = strings.TrimSpace(location)

	encodedLocation := url.QueryEscape(location)

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherKey, encodedLocation)
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

	fmt.Printf("Location: %s, %s, %s\n", weather.Location.Name, weather.Location.Region, weather.Location.Country)
	fmt.Printf("Temperature: %.2f°C\n", weather.Current.TempC)
	fmt.Printf("Condition: %s\n", weather.Current.Condition.Text)
}
