package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Coord structure
type Coord struct {
	Lon float32 `json:"lon"`
	Lat float32 `json:"lat"`
}

// Measure describes the json schema
type Measure struct {
	Coord      `json:"coord"` // Unused
	Conditions []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"` // Unused
	} `json:"weather"`
	Base     string `json:"base"` // Unused
	Measures struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`    // Unused
		ID      int     `json:"id"`      // Unused
		Message float64 `json:"message"` // Unused
		Country string  `json:"country"` // Unused
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	CityID   int    `json:"id"`   //? Check against Cities
	CityName string `json:"name"` //? Check against Cities
	Cod      int    `json:"cod"`  // Unused
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

	// Rate limiting
	limiter := time.Tick(time.Second)

	var cities []Measure
	for i := 0; i < 5; i++ {
		<-limiter

		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&units=metric&appid=%s", citiesCZ[i].ID, APIKEY)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error at ID:%v ----> %s\n", citiesCZ[i].ID, err)
		}

		var one Measure
		json.NewDecoder(resp.Body).Decode(&one)

		cities = append(cities, one)
	}
}
