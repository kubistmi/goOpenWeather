package main

import (
	sql "database/sql"
	"io/ioutil"
	"os"

	pq "github.com/lib/pq"
)

// UploadSQL takes the downloaded data and uploads them to postgreSQL database
func UploadSQL(weather *[]Measure, cities *[]City, path string) {

	var err error
	defer func() {
		if err != nil {
			Alert(err, Conf.Slack)
		}
	}()

	// DATABASE CONNECTION
	db, err := sql.Open("postgres", Conf.Psql)
	if err != nil {
		return
	}

	// BUILD TABLES IF NEEDED
	sqlDefFile, err := os.Open(path + "sql/table_definition.sql")
	if err != nil {
		return
	}
	defer sqlDefFile.Close()

	sqlDefBytes, err := ioutil.ReadAll(sqlDefFile)
	if err != nil {
		return
	}

	sqlDef := string(sqlDefBytes)

	_, err = db.Exec(sqlDef)
	if err != nil {
		return
	}

	// CITIES TRANSACTION PREPARATION
	trnc, err := db.Begin()
	if err != nil {
		return
	}

	stmc, err := trnc.Prepare(
		pq.CopyIn(
			"cities", //table
			"id", "city", "country", "lon", "lat"))
	if err != nil {
		return
	}

	for _, record := range *cities {
		_, err = stmc.Exec(
			record.ID, record.Name, record.Country, record.Coord.Lon, record.Coord.Lat)
		if err != nil {
			return
		}
	}

	// CITIES FLUSH AND COMMIT
	_, err = stmc.Exec()
	if err != nil {
		return
	}

	err = stmc.Close()
	if err != nil {
		return
	}

	err = trnc.Commit()
	if err != nil {
		return
	}

	// WEATHER TRANSACTION PREPARATION
	trn, err := db.Begin()
	if err != nil {
		return
	}

	stm, err := trn.Prepare(
		pq.CopyIn(
			"weather", //table
			"city_id", "conditions", "temperature", "pressure", "humidity", "temp_min", "temp_max", "visibility",
			"winddir", "windspeed", "clouds", "sunrise", "sunset", "timezone", "extraction_time", "batch"))
	if err != nil {
		return
	}

	for _, record := range *weather {
		_, err = stm.Exec(
			record.CityID, record.Conditions.Value(), record.Measures.Temp, record.Measures.Pressure, record.Measures.Humidity,
			record.Measures.TempMin, record.Measures.TempMax, record.Visibility, record.Wind.Deg, record.Wind.Speed, record.Clouds.All,
			record.Sys.Sunrise, record.Sys.Sunset, record.Timezone, record.Dt, Batch)
		if err != nil {
			return
		}
	}

	// WATHER FLUSH AND COMMIT
	_, err = stm.Exec()
	if err != nil {
		return
	}

	err = stm.Close()
	if err != nil {
		return
	}

	err = trn.Commit()
	if err != nil {
		return
	}

	db.Close()
}
