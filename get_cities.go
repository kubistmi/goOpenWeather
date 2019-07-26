package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// City structure
type City struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Coord   `json:"coord"`
}

// GetCities is used to collect the file of City IDs from OpenWeatherMap
func GetCities() []City {

	var empty []City
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Get("http://bulk.openweathermap.org/sample/city.list.json.gz")
	if err != nil {
		return empty
	}

	defer resp.Body.Close()

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return empty
	}
	defer gz.Close()

	jsonData, err := ioutil.ReadAll(gz)
	if err != nil {
		return empty
	}

	var cities []City
	json.Unmarshal(jsonData, &cities)

	var citiesCZ []City
	for _, val := range cities {
		if val.Country == "CZ" {
			citiesCZ = append(citiesCZ, val)
		}
	}

	return citiesCZ
}
