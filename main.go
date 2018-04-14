package main

import (
	"bitbucket.org/luno/project-sanjiv/model"
	"bitbucket.org/luno/project-sanjiv/util"
	"encoding/json"
	"github.com/microdevs/missy/log"
	"github.com/microdevs/missy/service"
	"net/http"
)

func main() {

	s := service.New()
	s.HandleFunc("/generateBuyQuote", generateBuyQuoteHandler).Methods("GET")
	s.HandleFunc("/generateSellQuote", generateSellQuoteHandler).Methods("GET")
	log.Print("Starting service")
	s.Start()

}

func generateBuyQuoteHandler(w http.ResponseWriter, r *http.Request) {
	generateQuoteHandler(w, r, model.SELL)
}

func generateSellQuoteHandler(w http.ResponseWriter, r *http.Request) {
	generateQuoteHandler(w, r, model.BUY)
}

func generateQuoteHandler(w http.ResponseWriter, r *http.Request, buyOrSell model.BuyOrSell) {

	quantity, err := util.GetDecimal(w, r, "quantity")
	if err != nil {
		return
	}
	log.Print("quantity got: ", quantity)

	orderBook := model.OrderBook{}
	err = orderBook.Populate()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	quote, err := orderBook.GenerateQuote(quantity, buyOrSell)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Print("About to json.encode quote with price: ", quote.Price)
	if err := json.NewEncoder(w).Encode(&quote); err != nil {
		panic(err)
	}

}
