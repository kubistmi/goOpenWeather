package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"
	"strconv"
)

// Coord structure
type Coord struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

func main() {

	batch, err := strconv.Atoi(time.Now().Format("2006010215"))
	if err != nil {
		log.Fatal(err)
	}

	path := "/home/kubistmi/go/src/weather/"

	keyFile, err := os.Open(path + "api")
	if err != nil {
		log.Fatal(err)
	}

	APIKEY, err := ioutil.ReadAll(keyFile)
	if err != nil {
		log.Fatal(err)
	}

	citiesCZ := GetCities()
	weather := GetWeather(&citiesCZ, APIKEY, time.Second)
	UploadSQL(&weather, &citiesCZ, path, batch)
}