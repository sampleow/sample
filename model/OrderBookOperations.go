package model

//go:generate moq -out OrderBookOperations_moq_test.go . OrderBookOperations
type OrderBookOperations interface {
	Populate() error
	GenerateQuote(quantity float64, buyOrSell BuyOrSell) (*Quote, error)
}
