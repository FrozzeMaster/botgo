package realbot

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/kr/pretty"
)

//GetTickerPrice returns last price for the Pair
func (bot *RealBot) GetTickerPrice(pair string) float64 {
	prices, err := futuresClient.NewListPricesService().Symbol(pair).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	finalPrice, _ := strconv.ParseFloat(prices[0].Price, 32)
	return finalPrice
}

//GetAccountInfo retrieves actual state of account
func (bot *RealBot) GetAccountInfo() {
	res, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	bot.Account = res
}

//GetBalanceInfo retrieves actual balance of account
func (bot *RealBot) GetBalanceInfo() {
	res, err := futuresClient.NewGetBalanceService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	bot.Balance, _ = strconv.ParseFloat(res[1].Balance, 64)
}

//GetOrderInfo retrieves actual orders
func (bot *RealBot) GetOrderInfo(orderID int64) *futures.Order {
	res, err := futuresClient.NewGetOrderService().Symbol("BTCUSDT").OrderID(orderID).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return res
}

//ChangeLeverage changes leverage
func (bot *RealBot) ChangeLeverage(leverage int, pair string) int {
	ok, err := futuresClient.NewChangeLeverageService().Leverage(leverage).Symbol(pair).Do(context.Background())
	pretty.Println(ok)
	if err != nil {
		fmt.Println(err)
	}
	return leverage
}

//CollectKlines is downloading futures candles and returns as slicee of structs
func (bot *RealBot) CollectKlines(pair string, interval string, pairIndex int, intervalIndex int, limit int) {
	klines, err := futuresClient.NewKlinesService().Symbol(pair).
		Interval(interval).Limit(limit).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	bot.KlinesData[pairIndex][intervalIndex] = klines
}

//CancelOrder is canceling order
func (bot *RealBot) CancelOrder(symbol string, orderId int64, clientId string) {
	cancelOrder := futuresClient.NewCancelOrderService()
	cancelOrder = cancelOrder.Symbol(symbol)
	cancelOrder = cancelOrder.OrderID(orderId)
	cancelOrder = cancelOrder.OrigClientOrderID(clientId)
	res, err := cancelOrder.Do(context.Background())
	if err != nil {
		fmt.Println("Wystąpił błąd w anulowaniu zlecenia")
	}
	pretty.Println(res)
}

//GetExchangeInfo retireiev futures exchange data
func (bot *RealBot) GetExchangeInfo() {
	exchange, err := futuresClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	bot.ExchangeInfo = exchange
}

//MakeMarketOrder makes an order
func (bot *RealBot) MakeMarketOrder(leverage int, pair string, orderSide string) *futures.CreateOrderResponse {
	fmt.Println("Making market order")

	//Getting last PRICE
	actualPrice := bot.GetTickerPrice(pair)
	//Getting available USDT Balance
	bot.GetBalanceInfo()
	walletBalance := bot.Balance
	//Setting up LEVERAGE
	bot.ChangeLeverage(leverage, pair)

	//Calculaing Max Amount To play with
	quantityNoFee := math.Floor(walletBalance/actualPrice*float64(leverage)*1000) / 1000
	//Again calculating Max quantity considering FEE
	quantity := quantityNoFee * (1 - bot.FeeTaker)
	quantity = bot.ReturnQuantityWithPrecision(quantity, pair)
	//Creating order
	order := futuresClient.NewCreateOrderService()
	order = order.Type(futures.OrderTypeMarket)
	order = order.Symbol(pair)
	switch orderSide {
	case "BUY":
		order = order.Side(futures.SideTypeBuy)
	case "SELL":
		order = order.Side(futures.SideTypeSell)
	}
	order = order.Quantity(fmt.Sprintf("%g", quantity))
	result, err := order.Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return result
}

//SetLimitStoploss makes an order
func (bot *RealBot) SetLimitTakeProfit(price float64, quantity float64, pair string, side string) *futures.CreateOrderResponse {
	fmt.Println("Setting takeprofit")
	price = bot.ReturnPriceWithPrecision(price, pair)
	//Creating order
	order := futuresClient.NewCreateOrderService()
	order = order.Type(futures.OrderTypeTakeProfitMarket)
	order = order.Symbol(pair)
	order = order.StopPrice(fmt.Sprintf("%g", price))
	switch side {
	case "BUY":
		order = order.Side(futures.SideTypeBuy)
	case "SELL":
		order = order.Side(futures.SideTypeSell)
	}
	order = order.Quantity(fmt.Sprintf("%g", quantity))
	result, err := order.Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return result
}

//MakeMarketOrder makes an order
func (bot *RealBot) SetLimitStoploss(price float64, quantity float64, pair string, side string) *futures.CreateOrderResponse {
	fmt.Println("Setting stoploss")
	price = bot.ReturnPriceWithPrecision(price, pair)
	//Creating order
	order := futuresClient.NewCreateOrderService()
	order = order.Type(futures.OrderTypeStopMarket)
	order = order.Symbol(pair)
	order = order.StopPrice(fmt.Sprintf("%g", price))
	switch side {
	case "BUY":
		order = order.Side(futures.SideTypeBuy)
	case "SELL":
		order = order.Side(futures.SideTypeSell)
	}
	order = order.Quantity(fmt.Sprintf("%g", quantity))
	result, err := order.Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return result
}

//MakeCompleteOrder makes order with stoploss na dstoplimit
func (bot *RealBot) MakeCompleteOrder(leverage int, pair string, tpPrice float64, slPrice float64, orderSide string) {
	actualPrice := bot.GetTickerPrice(pair)
	switch orderSide {
	case "BUY":
		if slPrice >= actualPrice || tpPrice <= actualPrice {
			fmt.Println("Wrong stoploss or take profit")
			return
		}
	case "SELL":
		if slPrice <= actualPrice || tpPrice >= actualPrice {
			fmt.Println("Wrong stoploss or take profit")
			return
		}
	}
	if bot.NumberOfActivePositions() != 0 {
		fmt.Println("Bot didnt make trade- already in position")
		return
	}
	result := bot.MakeMarketOrder(leverage, pair, orderSide)
	quantity, _ := strconv.ParseFloat(result.OrigQuantity, 64)
	if quantity == 0 {
		return
	}
	switch orderSide {
	case "BUY":
		stoploss := bot.SetLimitStoploss(slPrice, quantity, pair, "SELL")
		takeprofit := bot.SetLimitTakeProfit(tpPrice, quantity, pair, "SELL")
		var into []*futures.CreateOrderResponse
		into = append(into, stoploss)
		into = append(into, takeprofit)
		bot.OpenOrders = append(bot.OpenOrders, into)
	case "SELL":
		stoploss := bot.SetLimitStoploss(slPrice, quantity, pair, "BUY")
		takeprofit := bot.SetLimitTakeProfit(tpPrice, quantity, pair, "BUY")
		var into []*futures.CreateOrderResponse
		into = append(into, stoploss)
		into = append(into, takeprofit)
		bot.OpenOrders = append(bot.OpenOrders, into)
	}
	pretty.Println(bot.OpenOrders)
}
