package realbot

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/kr/pretty"
)

//UpdateCloseOrder closes second order after getting filled
func (bot *RealBot) UpdateCloseOrder(order *futures.WsOrderTradeUpdate) {
	orderId := order.ID
	orderClient := order.ClientOrderID
	orderSymbol := order.Symbol
	orderStatus := order.Status
	if orderStatus != "FILLED" {
		return
	}
	deleteRow := false
	deleteIndex := -1
	for o, v := range bot.OpenOrders {
		found := false
		index := -1
		for i, x := range v {
			if x.OrderID == orderId && x.ClientOrderID == orderClient && x.Symbol == orderSymbol {
				found = true
				index = i
			}
		}
		for i, x := range v {
			if found && index != i {
				deleteRow = true
				deleteIndex = o
				bot.CancelOrder(x.Symbol, x.OrderID, x.ClientOrderID)
			}
		}
	}
	if deleteRow {
		bot.OpenOrders[len(bot.OpenOrders)-1], bot.OpenOrders[deleteIndex] = bot.OpenOrders[deleteIndex], bot.OpenOrders[len(bot.OpenOrders)-1]
		bot.OpenOrders = bot.OpenOrders[:len(bot.OpenOrders)-1]
	}
	pretty.Println(bot.OpenOrders)
}

//NumberOfActivePositions return number of active positions
func (bot *RealBot) NumberOfActivePositions() int {
	openPositions := 0
	for _, v := range bot.Account.Positions {
		ep, _ := strconv.ParseFloat(v.EntryPrice, 32)
		if ep != 0 {
			openPositions++
		}
	}
	bot.ActivePositions = openPositions
	return openPositions
}

//ReturnPriceWithPrecision returns price with precision
func (bot *RealBot) ReturnPriceWithPrecision(price float64, pair string) float64 {
	var indexOfPair int
	for i, v := range bot.Pairs {
		if pair == v {
			indexOfPair = i
		}
	}
	pricePrecision := bot.PrecisionTable[indexOfPair][0]
	return bot.RoundToPrecisionCore(price, pricePrecision)
}

//ReturnQuantityWithPrecision returns quantity with precision
func (bot *RealBot) ReturnQuantityWithPrecision(quantity float64, pair string) float64 {
	var indexOfPair int
	for i, v := range bot.Pairs {
		if pair == v {
			indexOfPair = i
		}
	}
	pricePrecision := bot.PrecisionTable[indexOfPair][1]
	return bot.RoundToPrecisionCore(quantity, pricePrecision)
}

//RoundToPrecisionCore returns number with accurate precision
func (bot *RealBot) RoundToPrecisionCore(number float64, precision int) float64 {
	number = number * math.Pow10(precision)
	number = math.Floor(number)
	number = number / math.Pow10(precision)
	return number
}

//GetPrecisionTable
func (bot *RealBot) GetPrecisionTable() {
	for _, v := range bot.Pairs {
		for _, symbols := range bot.ExchangeInfo.Symbols {
			if v == symbols.Pair && v == symbols.Symbol {
				precisions := []int{symbols.PricePrecision, symbols.QuantityPrecision}
				bot.PrecisionTable = append(bot.PrecisionTable, precisions)
			}
		}
	}
}

//ShowlastKlines shows klines
func (bot *RealBot) ShowSavedKlines() {
	for i, v := range bot.CustomKline[0][0] {
		fmt.Print(i)
		fmt.Print(" ")
		fmt.Println(v)
		fmt.Println()
	}
}

//ShowlastSingleKlines shows klines
func (bot *RealBot) ShowSavedSingleKlines(pairNum int, intervalNum int) {
	fmt.Println(bot.CustomKline[pairNum][intervalNum][len(bot.CustomKline[pairNum][intervalNum])-1])
}

//DisplaTime displays current time

func (bot *RealBot) DisplayTime() string {
	now := time.Now()
	return now.Format("15:04:05 01-02-2006")
}

func (bot *RealBot) CalculateBollinger(newKline *model.MyKline, bollingerLength int, bollingerValue float64, atrValue int, intervalIndex int, pairIndex int) (float64, []float64) {
	var maxLengthofIndicator int
	if bollingerLength > atrValue {
		maxLengthofIndicator = bollingerLength
	} else {
		maxLengthofIndicator = atrValue
	}
	//Length of Important candlesticks
	liC := len(bot.CustomKline[pairIndex][intervalIndex])
	//Important candlesticks
	iC := bot.CustomKline[pairIndex][intervalIndex][liC-maxLengthofIndicator-1 : liC]
	iC = append(iC, newKline)

	liCC := len(iC) - 1
	averageBollinger := 0.00
	stDevParts := 0.00
	atr := 0.00
	for i := liCC; i > liCC-bollingerLength; i-- {
		averageBollinger += iC[i].Close
	}
	averageBollinger = averageBollinger / float64(bollingerLength)
	for i := liCC; i > liCC-bollingerLength; i-- {
		stDevParts += math.Pow((iC[i].Close - averageBollinger), 2)
	}
	stDev := stDevParts / (float64(bollingerLength))
	stDev = math.Sqrt(stDev)
	upperBand := averageBollinger + bollingerValue*stDev
	lowerBand := averageBollinger - bollingerValue*stDev
	for i := liCC; i > liCC-atrValue; i-- {
		atr += iC[i].Max - iC[i].Min
	}
	newAtrValue := atr / float64(atrValue)
	newBollingerValues := []float64{upperBand, averageBollinger, lowerBand}
	return newAtrValue, newBollingerValues

}
