package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slackman"
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

// Alert takes care of sending the Slack message and logging the error
func Alert(err error, API string) {
	text := fmt.Sprintf("The loading batch [%v] failed due to the following error: %v \n", Batch, err)
	msg := slackman.NewMessage(API, "#log---weather", "GoLog", text, "https://img.icons8.com/cotton/2x/server.png")
	msg.Send()
	log.Fatal(err)
}

func main() {
	Batch, _ = strconv.Atoi(time.Now().Format("2006010215"))

	//path := os.Getenv("GOPATH") + "/src/weather/"
	path := "/home/kubistmi/go/src/weather/"

	conFile, err := os.Open(path + "config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer conFile.Close()

	json.NewDecoder(conFile).Decode(&Conf)

	citiesCZ := GetCities()
	weather := GetWeather(&citiesCZ, time.Second)
	UploadSQL(&weather, &citiesCZ, path)

	log.Printf("Finished the loading %v\n", Batch)
}
