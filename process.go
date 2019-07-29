package main

import (
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
	weather := GetWeather(&citiesCZ, APIKEY, time.Second)
	UploadSQL(&weather, &citiesCZ)
}
