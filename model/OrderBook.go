package model

import (
	"bitbucket.org/luno/project-sanjiv/util"
	"encoding/json"
	"errors"
	"github.com/microdevs/missy/config"
	"github.com/microdevs/missy/log"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type BuyOrSell string

const (
	BUY  BuyOrSell = "BUY"
	SELL BuyOrSell = "SELL"
)

func (bs BuyOrSell) IsBuy() bool {
	return bs == BUY
}
func (bs BuyOrSell) IsSell() bool {
	return bs == SELL
}

type OrderBook struct {
	Timestamp int64          `json:"timestamp"`
	Bid       []*PriceVolume `json:"bids"`
	Ask       []*PriceVolume `json:"asks"`
}

type PriceVolume struct {
	Volume string `json:"volume"`
	Price  string `json:"price"`
}

func (orderBook *OrderBook) Populate() (err error) {
	url := config.Get("bitx.orderbook.url")
	log.Println("url: ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("NewRequest: ", err)
		return errors.New("Error trying to generate request to api.mybitx.com")
	}

	client := &http.Client{
		// set a 5s timeout
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Do: ", err)
		return errors.New("could not connect to api.mybitx.com")
	}

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(orderBook); err != nil {
		log.Print(err)
		return errors.New("Error unmarshalling data from api.mybitx.com")
	}

	if orderBook.Bid == nil && orderBook.Ask == nil {
		s := "could not get data from api.mybitx.com"
		log.Print(s)
		return errors.New(s)
	} else {
		log.Printf("Populated %d ASK and %d BID elements", len(orderBook.Ask), len(orderBook.Bid))
	}
	return err
}

func (orderBook *OrderBook) GenerateQuote(quantity decimal.Decimal, buyOrSell BuyOrSell) (quote *Quote, err error) {
	quotePrice := decimal.NewFromFloat(0.0)
	qtyRem := quantity
	log.Println("quantity: ", qtyRem)

	var ptrToUse []*PriceVolume
	if buyOrSell.IsBuy() {
		ptrToUse = orderBook.Ask
	} else {
		ptrToUse = orderBook.Bid
	}

	for _, pv := range ptrToUse {
		vol, err := util.GetDecimalFromString(pv.Volume)
		if err != nil {
			return nil, err
		}
		thePrice, err := util.GetDecimalFromString(pv.Price)
		if err != nil {
			return nil, err
		}

		if vol.GreaterThan(qtyRem) == true {
			quotePrice = quotePrice.Add(thePrice.Mul(qtyRem))
			log.Println("qtyRemaining: ", qtyRem, ", vol was: ", vol, ", price was: ", thePrice, ", quotePrice now: ", quotePrice.StringFixed(2))
			qtyRem = qtyRem.Sub(qtyRem)
			break
		}
		quotePrice = quotePrice.Add(vol.Mul(thePrice))
		qtyRem = qtyRem.Sub(vol)
		log.Println("qtyRemaining: ", qtyRem, ", vol was: ", vol, ", price was: ", thePrice, ", quotePrice now: ", quotePrice)
	}
	if qtyRem.Equal(decimal.NewFromFloat(0.0)) == false {
		quantity = quantity.Sub(qtyRem)
	}
	quote = &Quote{
		Quantity: quantity.String(),
		Price:    quotePrice.StringFixed(2),
	}
	return quote, err
}
