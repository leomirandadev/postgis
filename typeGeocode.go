package main

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/twpayne/go-geom/encoding/wkb"
)

type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (g GeoPoint) GormDataType() string {
	return "GEOMETRY(Point)"
}

func (g GeoPoint) GormDBDataType() string {
	return "GEOMETRY(Point)"
}

func (p *GeoPoint) String() string {
	return fmt.Sprintf("SRID=4326;POINT(%v %v)", p.Lng, p.Lat)
}

func (p *GeoPoint) Scan(val interface{}) error {
	bs, err := hex.DecodeString(val.(string))
	if err != nil {
		return err
	}

	data := wkb.Geom{}
	if err = data.Scan(bs); err != nil {
		return err
	}

	geocode := data.FlatCoords()
	if len(geocode) < 2 {
		return errors.New("latitude and longitude don't match the pattern")
	}

	p.Lat = geocode[0]
	p.Lng = geocode[1]

	return nil
}

func (p GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}
