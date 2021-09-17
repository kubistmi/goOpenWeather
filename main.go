package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
}

// Conf is prepared for JSON unmarshalling
var Conf Configuration

// Batch describes the process ID
var Batch int

// Alert takes care of logging the error
func Alert(err error, fatal bool) {
	text := fmt.Sprintf("The loading batch [%v] failed due to the following error: %v \n", Batch, err)
	if fatal {
		log.Fatal(text)
	}
}

func main() {
	Batch, _ = strconv.Atoi(time.Now().Format("2006010215"))

	path := os.Getenv("HOME") + "/go/src/weather/"

	conFile, err := os.Open(path + "config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer conFile.Close()

	err = json.NewDecoder(conFile).Decode(&Conf)
	if err != nil {
		log.Fatal(err)
	}

	citiesCZ, err := GetCities()
	if err != nil {
		Alert(err, true)
	}

	weather, err := GetWeather(&citiesCZ, time.Second)
	if err != nil {
		Alert(err, true)
	}

	err = UploadSQL(&weather, &citiesCZ, path)
	if err != nil {
		Alert(err, true)
	}

	// Log success
	finished := fmt.Sprintf("Finished the loading %v\n", Batch)
	log.Print(finished)
}
