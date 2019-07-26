package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Coord structure
type Coord struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

func main() {

	keyFile, err := os.Open("api")
	if err != nil {
		log.Fatal(err)
	}

	APIKEY, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	citiesCZ := GetCities()
	weather := GetWeather(citiesCZ, APIKEY, time.Second)

	fmt.Println(weather[0])
	fmt.Println(weather[len(weather)-1])

	weatherJSON, err := json.Marshal(weather)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("weather.json", weatherJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
