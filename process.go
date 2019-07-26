package main

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
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
	citiesCZ := GetCities()

	var cities []Measure
	for i := 0; i < 5; i++ {
		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&units=metric&appid=%s", citiesCZ[i].ID, "2ef0af7b4735394260790c58a56f8810")

		resp, err := http.Get(url)

		if err != nil {
			log.Fatal(err)
		}

		var one Measure
		json.NewDecoder(resp.Body).Decode(&one)
		fmt.Printf("%v | %s | %s | %v \n", one.ID, citiesCZ[i].Name, one.Weather[0].Description, one.Main.Temp)

		cities = append(cities, one)
	}
}
