package geocode

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := "010100000085EB51B81E0524400000000000003440"
	geocode := GeoPoint{}
	geocode.Scan(input)

	output := geocode.String()

	println(output == input)
	if output != input {
		t.Error(input, output)
	}
}
