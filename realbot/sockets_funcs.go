package realbot

import (
	"fmt"

	rmf "github.com/CraZzier/bot/realbot/utilities"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/kr/pretty"
)

//CollectKlinesWS is websocket that every 2000ms gets data from binance server
func (bot *RealBot) CollectKlinesWS(pair string, interval string, pairIndex int, intervalIndex int, responseChannel5m chan int, responseChannel15m chan int, responseChannel1h chan int) {
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		if event.Kline.IsFinal && bot.KlinesData[pairIndex][intervalIndex][len(bot.KlinesData)-1].OpenTime != event.Kline.StartTime {
			lastKline := rmf.ConvertWsKlineToKline(event.Kline)
			bot.CustomKline[pairIndex][intervalIndex] = append(bot.CustomKline[pairIndex][intervalIndex], rmf.ToMyKlineSingle(&lastKline))
			if intervalIndex == 0 {
				responseChannel5m <- 1
			}
			if intervalIndex == 1 {
				responseChannel15m <- 1
			}
			if intervalIndex == 2 {
				responseChannel1h <- 1
			}
		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, _, err := futures.WsKlineServe(pair, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//CollectKlinesWS is websocket that every 2000ms gets data from binance server
func (bot *RealBot) CollectKlinesWSa(pair string, interval string) {
	wsKlineHandler := func(event *futures.WsKlineEvent) {
		if event.Kline.IsFinal {
			lastKline := rmf.ConvertWsKlineToKline(event.Kline)
			pretty.Println(lastKline)
		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, _, err := futures.WsKlineServe(pair, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//UserStreamWS updates trades
func (bot *RealBot) UserStreamWS() {
	wsUserDataHandler := func(event *futures.WsUserDataEvent) {
		orderSample := futures.WsOrderTradeUpdate{}
		orderReal := event.OrderTradeUpdate
		if orderReal != orderSample {
			bot.UpdateCloseOrder(&orderReal)
		}

	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	_, _, err := futures.WsUserDataServe(bot.ListenKey, wsUserDataHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
}
