package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		log.Fatal("Error loading .env file")
	}

	weatherKey := os.Getenv("WEATHER_TOKEN")

	var location string
	fmt.Print("Enter location: ")
	_, err := fmt.Scanln(&location)
	if err != nil {
		log.Fatalf("Error reading location: %v", err)
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", weatherKey, location)
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(response.StatusCode)
	//fmt.Println(string(responseData))

	var result map[string]interface{}
	err = json.Unmarshal(responseData, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	locationMap, ok := result["location"].(map[string]interface{})
	if !ok {
		log.Fatal("error: 'location' field is not of type map[string]interface{}")
	}
	name, ok := locationMap["name"].(string)
	if !ok {
		log.Fatal("error: 'name' field is not of type string")
	}
	region, ok := locationMap["region"].(string)
	if !ok {
		log.Fatal("error: 'region' field is not of type string")
	}
	country, ok := locationMap["country"].(string)
	if !ok {
		log.Fatal("error: 'country' field is not of type string")
	}

	currentMap, ok := result["current"].(map[string]interface{})
	if !ok {
		log.Fatal("error: 'current' field is not of type map[string]interface{}")
	}
	tempC, ok := currentMap["temp_c"].(float64)
	if !ok {
		log.Fatal("error: 'temp_c' field is not of type float64")
	}
	conditionMap, ok := currentMap["condition"].(map[string]interface{})
	if !ok {
		log.Fatal("error: 'condition' field is not of type map[string]interface{}")
	}
	conditionText, ok := conditionMap["text"].(string)
	if !ok {
		log.Fatal("error: 'text' field is not of type string")
	}

	fmt.Printf("Location: %s, %s, %s\n", name, region, country)
	fmt.Printf("Temperature: %.2f°C\n", tempC)
	fmt.Printf("Condition: %s\n", conditionText)
}
