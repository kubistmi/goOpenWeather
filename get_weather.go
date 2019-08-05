package main

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Condition describes current weather condition
type Condition struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Conditions wraps the array of distinct conditions
type Conditions []Condition

// Value method is used to implement the driver.Valuer interface
func (c Conditions) Value() driver.Value {
	val, err := json.Marshal(c)
	if err != nil {
		Alert(err, Conf.Slack, Batch)
	}
	return val
}

// Measure describes the json schema
type Measure struct {
	Coord      `json:"coord"` // Unused
	Conditions `json:"weather"`
	Base       string `json:"base"` // Unused
	Measures   struct {
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
	CityID   int    `json:"id"`   // Unused
	CityName string `json:"name"` //? Check against Cities
	Cod      int    `json:"cod"`  // Unused
}

// GetWeather is used to collect all wetaher data for the provided cities.
func GetWeather(cities *[]City, rate time.Duration) []Measure {

	var empty []Measure
	var err error
	defer func() {
		if err != nil {
			Alert(err, Conf.Slack, Batch)
		}
	}()

	// Rate limiting
	limiter := time.Tick(rate)

	weather := make([]Measure, len(*cities))
	//weather := make([]Measure, 5) //! REMOVE AFTER TESTING

	//for ix, val := range cities {
	for ix, val := range *cities {
		//if ix < 5 { //! REMOVE AFTER TESTING
		<-limiter

		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?id=%v&units=metric&appid=%s", val.ID, Conf.Weather)

		resp, err := http.Get(url)
		if err != nil {
			err = fmt.Errorf("Error at ID:%v ----> %s", val.ID, err)
			return empty
		}

		var CityWeather Measure
		if resp.StatusCode == 200 {
			json.NewDecoder(resp.Body).Decode(&CityWeather)

			if val.Name != CityWeather.CityName {
				err = fmt.Errorf("Error at ID:%v ----> || The city names do not match between cities and weather json. This really shouldn't happen", val.ID)
				return empty
			}
		}

		weather[ix] = CityWeather
	}

	return weather
}
