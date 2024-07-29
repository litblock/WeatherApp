package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	weatherKey, ok := os.LookupEnv("WEATHER_TOKEN")

	if !ok {
		fmt.Println("WEATHER_TOKEN is not set")
		os.Exit(1)
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1?key=%s", weatherKey)
	response, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(responseData))
}
