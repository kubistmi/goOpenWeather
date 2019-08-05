package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"strconv"
	"time"
)

// Coord structure
type Coord struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

// Configuration defines the keys and connection strings used by the process
type Configuration struct {
	Weather string `json:"weather"`
	Psql    string `json:"psql"`
	Slack   string `json:"slack"`
}

// Conf is prepared for JSON unmarshalling
var Conf Configuration

// Batch describes the process ID
var Batch int

	batch, err := strconv.Atoi(time.Now().Format("2006010215"))
	if err != nil {
		log.Fatal(err)
	}

func main() {

	Batch, err := strconv.Atoi(time.Now().Format("2006010215"))
	if err != nil {
		log.Fatal(err)
	}

	path := os.Getenv("GOPATH") + "/src/weather/"

	conFile, err := os.Open(path + "config.json")
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(conFile).Decode(&Conf)

	citiesCZ := GetCities()
	weather := GetWeather(&citiesCZ, time.Second)
	UploadSQL(&weather, &citiesCZ, path)

	log.Printf("Finished the loading %v\n", Batch)
}
