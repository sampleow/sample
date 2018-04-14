package model_test

import (
	. "bitbucket.org/luno/project-sanjiv/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	decimal "github.com/shopspring/decimal"
)

var _ = Describe("OrderBook with data", func() {
	var (
		orderBook OrderBook
	)

	BeforeEach(func() {

		bids := []*PriceVolume{
			{Volume: "1.5", Price: "178500"},
			{Volume: "1", Price: "178000"},
			{Volume: "0.5", Price: "178000"},
		}
		orderBook.Bid = bids

		asks := []*PriceVolume{
			{Volume: "0.5", Price: "178500"},
			{Volume: "1", Price: "179000"},
			{Volume: "1.5", Price: "179500"},
		}
		orderBook.Ask = asks
	})

	Describe("Getting Quote", func() {
		Context("get BUY quote", func() {
			It("should be ((178500*0.5) + (179000*0.5))", func() {
				quote, _ := orderBook.GenerateQuote(decimal.NewFromFloat(1.0), BUY)
				Expect(quote.Price).To(Equal(decimal.NewFromFloat((178500 * 0.5) + (179000 * 0.5)).StringFixed(2)))
			})
		})

		Context("get SELL quote", func() {
			It("should be ((178500*1.5) + (178000*0.5))", func() {
				quote, _ := orderBook.GenerateQuote(decimal.NewFromFloat(2.0), SELL)
				Expect(quote.Price).To(Equal(decimal.NewFromFloat((178500 * 1.5) + (178000 * 0.5)).StringFixed(2)))
			})
		})
	})

})

var _ = Describe("OrderBook with insufficient volume", func() {
	var (
		orderBook OrderBook
	)

	BeforeEach(func() {

		bids := []*PriceVolume{
			{Volume: "1.5", Price: "178500"},
			{Volume: "1", Price: "178000"},
			{Volume: "0.5", Price: "178000"},
		}
		orderBook.Bid = bids

		asks := []*PriceVolume{
			{Volume: "0.5", Price: "178500"},
			{Volume: "1", Price: "179000"},
			{Volume: "1.5", Price: "179500"},
		}
		orderBook.Ask = asks
	})

	Describe("Getting Quote", func() {
		Context("get BUY quote", func() {
			It("should be ((178500*0.5) + (179000*1*1) + (179500*1.5))", func() {
				quote, _ := orderBook.GenerateQuote(decimal.NewFromFloat(5.0), BUY)
				Expect(quote.Price).To(Equal(decimal.NewFromFloat((178500 * 0.5) + (179000) + (1.5 * 179500)).StringFixed(2)))
				Expect(quote.Quantity).To(Equal(decimal.NewFromFloat(3).String()))
			})
		})

		Context("get SELL quote", func() {
			It("should be ((178500*1.5) + (178000*1) + (178000*0.5))", func() {
				quote, _ := orderBook.GenerateQuote(decimal.NewFromFloat(5.0), SELL)
				Expect(quote.Price).To(Equal(decimal.NewFromFloat((178500 * 1.5) + (178000 * 1) + (178000 * 0.5)).StringFixed(2)))
				Expect(quote.Quantity).To(Equal(decimal.NewFromFloat(3).String()))
			})
		})
	})

})
