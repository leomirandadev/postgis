package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db := newDB()

	db.migrate()

	places, err := db.getLocationsCloser(10.01, 20, 20)
	if err != nil {
		log.Fatal("[MAIN ERROR]", err)
	}

	fmt.Println(places)
}

type DB struct {
	dbConn *gorm.DB
}

func newDB() DB {
	conn, err := gorm.Open(
		postgres.Open("user=root password=root dbname=postgis_test host=localhost port=5432 sslmode=disable"),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("[ERROR]", err)
	}

	return DB{conn}
}

type Place struct {
	ID      int      `json:"id" gorm:"id"`
	Name    string   `json:"name" gorm:"name"`
	Geocode GeoPoint `json:"geocode" gorm:"geocode"`
}

func (d DB) getLocationsCloser(lat, lng, radiusMiles float64) ([]Place, error) {
	var rows []Place

	result := d.dbConn.Where(`
		ST_DistanceSphere(
			places.geocode,
			ST_MakePoint($1, $2)
		) <= $3
	`, lat, lng, milesToMeter(radiusMiles)).Find(&rows)

	return rows, result.Error
}

func (d DB) migrate() {

	db, _ := d.dbConn.DB()

	// turn on postgis
	_, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS postgis`)
	if err != nil {
		log.Fatal("[MIGRATE ERROR] ", err)
	}

	// migrate model
	err = d.dbConn.AutoMigrate(Place{})
	if err != nil {
		log.Fatal("[MIGRATE ERROR] ", err)
	}

	// check if the seed has already exists
	rows, _ := db.Exec(`SELECT * FROM places WHERE name = 'place 1'`)
	rowsAffected, _ := rows.RowsAffected()
	if rowsAffected > 0 {
		return
	}

	// seed
	_, err = db.Exec(`
		INSERT INTO places (name, geocode) 
		VALUES
			('place 1', 'POINT(10.01 20.0000)'),
			('place 2', 'POINT(10.20 20.1000)'),
			('place 3', 'POINT(100.0 10.0000)'),
			('place 4', 'POINT(100.11 10.002)')
	`)
	if err != nil {
		log.Fatal("[MIGRATE ERROR] ", err)
	}
}
