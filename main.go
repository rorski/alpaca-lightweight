package main

import (
	"log"
	"net/http"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/rorski/alpaca-lightweight/server"
)

func main() {
	conf := server.Config{
		Client: marketdata.NewClient(marketdata.ClientOpts{
			BaseURL: "https://data.alpaca.markets"},
		)}

	mux := http.NewServeMux()
	mux.HandleFunc("/chart", conf.Chart)

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal("error creating http server")
	}
}
