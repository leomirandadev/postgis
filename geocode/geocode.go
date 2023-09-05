package geocode

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"

	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/encoding/wkbhex"
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
	point := geom.NewPoint(geom.XY)
	point.SetCoords([]float64{p.Lat, p.Lng})

	value, _ := wkbhex.Encode(point, binary.LittleEndian)

	return strings.ToUpper(value)
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
