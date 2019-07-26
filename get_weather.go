package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

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

// GetWeather is used to collect all wetaher data for the provided cities.
func GetWeather(cities []City, APIKEY []byte, rate time.Duration) []Measure {

	// Rate limiting
	limiter := time.Tick(rate)

	weather := make([]Measure, len(cities))
	// weather := make([]Measure, 5) //! REMOVE AFTER TESTING

	for ix, val := range cities {
		//for ix, val := range cities[0:3] { //! REMOVE AFTER TESTING
		<-limiter

		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&units=metric&appid=%s", val.ID, APIKEY)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error at ID:%v ----> %s\n", val.ID, err)
		}

		var CityWeather Measure
		json.NewDecoder(resp.Body).Decode(&CityWeather)

		if val.Name != CityWeather.CityName {
			log.Fatalf("Error at ID:%v ----> || The city names do not match between cities and weather json. This really shouldn't happen. \n", val.ID)
		}

		weather[ix] = CityWeather
	}

	return weather
}
