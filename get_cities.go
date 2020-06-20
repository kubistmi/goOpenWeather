package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
func GetCities() ([]City, error) {

	var empty []City
	var err error

	resp, err := http.Get("http://bulk.openweathermap.org/sample/city.list.json.gz")
	if err != nil {
		return empty, fmt.Errorf("Error in API request ----> %s", err)
	}

	defer resp.Body.Close()

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return empty, fmt.Errorf("Can't prepare the reader for gzipped data ----> %s", err)
	}
	defer gz.Close()

	jsonData, err := ioutil.ReadAll(gz)
	if err != nil {
		return empty, fmt.Errorf("Can't read the gzipped data ----> %s", err)
	}

	var cities []City
	json.Unmarshal(jsonData, &cities)

	var citiesCZ []City
	for _, val := range cities {
		if val.Country == "CZ" {
			citiesCZ = append(citiesCZ, val)
		}
	}

	return citiesCZ, nil
}
