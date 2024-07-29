package main

import (
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

	fmt.Println(response.StatusCode)
	fmt.Println(string(responseData))
}
