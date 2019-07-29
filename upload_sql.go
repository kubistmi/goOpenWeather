package main

import (
	sql "database/sql"
	"io/ioutil"
	"log"
	"os"

	pq "github.com/lib/pq"
)

// UploadSQL takes the downloaded data and uploads them to postgreSQL database
func UploadSQL(weather *[]Measure, cities *[]City) {

	// DATABASE CONNECTION
	sqlFile, err := os.Open("connstr")
	if err != nil {
		log.Fatal(err)
	}

	sqlBytes, err := ioutil.ReadAll(sqlFile)
	if err != nil {
		log.Fatal(err)
	}

	SQLCONN := string(sqlBytes)

	db, err := sql.Open("postgres", SQLCONN)
	if err != nil {
		log.Fatal(err)
	}

	// CITIES TRANSACTION PREPARATION
	trnc, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmc, err := trnc.Prepare(
		pq.CopyIn(
			"cities", //table
			"id", "city", "country", "lon", "lat"))
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range *cities {
		_, err = stmc.Exec(
			record.ID, record.Name, record.Country, record.Coord.Lon, record.Coord.Lat)
		if err != nil {
			log.Fatal(err)
		}
	}

	// CITIES FLUSH AND COMMIT
	_, err = stmc.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmc.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = trnc.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// WEATHER TRANSACTION PREPARATION
	trn, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stm, err := trn.Prepare(
		pq.CopyIn(
			"weather", //table
			"city_id", "conditions", "temperature", "pressure", "humidity", "temp_min", "temp_max", "visibility",
			"winddir", "windspeed", "clouds", "sunrise", "sunset", "timezone", "extraction_time"))
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range *weather {
		_, err = stm.Exec(
			record.CityID, record.Conditions.Value(), record.Measures.Temp, record.Measures.Pressure, record.Measures.Humidity,
			record.Measures.TempMin, record.Measures.TempMax, record.Visibility, record.Wind.Deg, record.Wind.Speed, record.Clouds.All,
			record.Sys.Sunrise, record.Sys.Sunset, record.Timezone, record.Dt)
		if err != nil {
			log.Fatal(err)
		}
	}

	// WATHER FLUSH AND COMMIT
	_, err = stm.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stm.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = trn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
}
