package server

import (
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

// Config is a struct for server configuration. Currently it just embeds marketdata.Client.
type Config struct {
	*marketdata.Client
}

// Data is the data passed to the html template to be rendered.
type Data struct {
	Ticker     string
	Bars       []marketdata.Bar
	Technicals map[string]any
}

func calcTechnicals(bars []marketdata.Bar) map[string]any {
	var (
		totalVol           uint64
		maxPrice, minPrice float64
	)
	// calculate some basic technicals over the time period from the bars
	for _, bar := range bars {
		if bar.High > maxPrice {
			maxPrice = bar.High
		}
		if bar.Low < minPrice || minPrice == 0 {
			minPrice = bar.Low
		}
		totalVol += bar.Volume
	}

	technicals := map[string]any{
		"MaxPrice": maxPrice,
		"MinPrice": minPrice,
		"TotalVol": totalVol,
	}

	return technicals
}

// Chart is a handler for the /chart path, which formats the query parameters and passes it to the
// [Alpaca marketdata API] for processing. If successful, it renders the layout.html template with the data.
//
// [Alpaca marketdata API]: https://github.com/alpacahq/alpaca-trade-api-go/tree/master/marketdata
func (c *Config) Chart(w http.ResponseWriter, r *http.Request) {
	var (
		days   int
		ticker string
		err    error
	)

	ticker = r.URL.Query().Get("ticker")
	if ticker == "" {
		http.Error(w, "must specify a ticker symbol", http.StatusBadRequest)
		return
	}

	// default to 120 day look back
	if r.URL.Query().Get("days") == "" {
		days = 120
	} else {
		days, err = strconv.Atoi(r.URL.Query().Get("days"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	bars, err := c.GetBars(ticker, marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     time.Now().AddDate(0, 0, -days),
		End:       time.Now().AddDate(0, 0, -1),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Data{
		Ticker:     ticker,
		Bars:       bars,
		Technicals: calcTechnicals(bars),
	}

	tmpl := template.Must(template.New("layout.html").Funcs(
		template.FuncMap{
			"FormatBigNum": FormatBigNum,
		},
	).ParseFiles("layout.html"))

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
