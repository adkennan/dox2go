package dox2go

import "testing"

func checkConv(t *testing.T, fromVal float64, from Unit, to Unit, expected float64) {
	v := ConvertUnit(fromVal, from, to)
	if v != expected {
		t.Errorf("Expected %f. Was %f", expected, v)
	}
}

func TestConvertUnit(t *testing.T) {

	checkConv(t, 1.0, U_CM, U_MM, 10.0)
	checkConv(t, 1.0, U_MM, U_CM, 0.1)
	checkConv(t, 1.0, U_PT, U_IN, 1.0/72.0)
	checkConv(t, 1.0, U_IN, U_PT, 72.0)

	checkConv(t, 5.0, U_CM, U_MM, 50.0)
	checkConv(t, 5.0, U_MM, U_CM, 0.5)
	checkConv(t, 5.0, U_PT, U_IN, (1.0/72.0)*5.0)
	checkConv(t, 5.0, U_IN, U_PT, 72.0*5.0)
}
