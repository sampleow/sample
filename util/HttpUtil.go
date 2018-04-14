package util

import (
	"errors"
	"github.com/microdevs/missy/log"
	"github.com/shopspring/decimal"
	"net/http"
)

func GetDecimal(w http.ResponseWriter, r *http.Request, paramName string) (quantity decimal.Decimal, err error) {
	query := r.URL.Query()
	qtyParam := query.Get(paramName)
	quantity, err = GetDecimalFromString(qtyParam)
	if err != nil || quantity.LessThan(decimal.NewFromFloat(0.0)) {
		s := "value: " + qtyParam + ", invalid for param: " + paramName
		log.Print(s)
		err = errors.New(s)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(s))
	}
	return quantity, err
}

func GetDecimalFromString(qtyParam string) (price decimal.Decimal, err error) {
	price, err = decimal.NewFromString(qtyParam)
	if err != nil {
		log.Print("value: ", qtyParam, ", invalid for conversion to decimal: ", err)
		err = errors.New("Invalid decimal value: " + qtyParam)
	}
	return
}
