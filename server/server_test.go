package server

import (
	"testing"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/stretchr/testify/assert"
)

func TestBigNum(t *testing.T) {
	bigNum := FormatBigNum(uint64(123456789))
	assert.Equal(t, bigNum, "123,456,789")
}

func TestSmallNum(t *testing.T) {
	bigNum := FormatBigNum(uint64(123))
	assert.Equal(t, bigNum, "123")
}

func TestCalcTechnicals(t *testing.T) {
	bars := []marketdata.Bar{
		{
			High:   12.34,
			Low:    11.20,
			Volume: 1,
		},
		{
			High:   34.45,
			Low:    34.40,
			Volume: 3,
		},
		{
			High:   56.78,
			Low:    51.27,
			Volume: 4,
		},
	}

	technicals := calcTechnicals(bars)

	assert.Equal(t, technicals["MaxPrice"], float64(56.78))
	assert.Equal(t, technicals["MinPrice"], float64(11.20))
	assert.Equal(t, technicals["TotalVol"], uint64(8))
}
