package utilities

import (
	"fmt"

	"github.com/CraZzier/bot/model"
)

//UpdateMACD updates MACD using the last given candle
func UpdateMACD(candlesticks []*model.MyKline, macd *model.MACD) ([]*model.MyKline, *model.MACD) {

	//Getting last Ema values
	lastEma1 := macd.E1Keys[len(macd.E1Keys)-1].Value
	lastEma2 := macd.E2Keys[len(macd.E2Keys)-1].Value

	//Checking if emastampare correct
	if macd.E2Keys[len(macd.E1Keys)-1].Timestamp != macd.E2Keys[len(macd.E2Keys)-1].Timestamp {
		fmt.Println("błąd w prgramie w ema")
	}
	//Calculating new values
	var candleVal float64
	timestamp := candlesticks[len(candlesticks)-1].OpenTime
	switch macd.WhichValue {
	case "close":
		candleVal = candlesticks[len(candlesticks)-1].Close
	case "open":
		candleVal = candlesticks[len(candlesticks)-1].Open
	}
	newMargin1 := float64(2.00 / (1.00 + float64(macd.E1)))
	newMargin2 := float64(2.00 / (1.00 + float64(macd.E2)))
	newMarginMacd := float64(2.00 / (1.00 + float64(macd.Signal)))
	newEma1 := candleVal*newMargin1 + lastEma1*(1.00-newMargin1)
	newEma2 := candleVal*newMargin2 + lastEma2*(1.00-newMargin2)
	newMacd := newEma1 - newEma2
	newSignal := macd.Keys[len(macd.Keys)-1].Value[1]*(1-newMarginMacd) + newMacd*newMarginMacd
	newDiff := newMacd - newSignal

	//Creating objects to be inserted into pointers
	newMacdStamp := &model.SingleMACDStamp{
		Timestamp: timestamp,
		Value:     []float64{newMacd, newSignal, newDiff},
	}
	newEma1Stamp := &model.SingleMovingAverageStamp{
		Timestamp: timestamp,
		Value:     newEma1,
	}
	newEma2Stamp := &model.SingleMovingAverageStamp{
		Timestamp: timestamp,
		Value:     newEma2,
	}
	//Updating Values inside MACD and candlesticks
	macd.Keys = append(macd.Keys, newMacdStamp)
	macd.E1Keys = append(macd.E1Keys, newEma1Stamp)
	macd.E2Keys = append(macd.E2Keys, newEma2Stamp)
	candlesticks[len(candlesticks)-1].MacD = append(candlesticks[len(candlesticks)-1].MacD, []float64{newMacd, newSignal, newDiff})

	//Removing first element of an arrays to clear memory
	macd.Keys = macd.Keys[1:]
	macd.E1Keys = macd.E1Keys[1:]
	macd.E2Keys = macd.E2Keys[1:]
	candlesticks = candlesticks[1:]

	return candlesticks, macd

}

//UpdateEMA
func UpdateEMA(candlesticks []*model.MyKline, ema *model.MovingAverage) ([]*model.MyKline, *model.MovingAverage) {
	var candleVal float64

	switch ema.WhichValue {
	case "close":
		candleVal = candlesticks[len(candlesticks)-1].Close
	case "open":
		candleVal = candlesticks[len(candlesticks)-1].Open
	}
	lastEma := ema.Keys[len(ema.Keys)-1].Value
	newMargin := float64(2.00 / (1.00 + float64(ema.IntervalValue)))
	newEma := candleVal*newMargin + lastEma*(1.00-newMargin)

	candlesticks[len(candlesticks)-1].Emas = append(candlesticks[len(candlesticks)-1].Emas, newEma)
	ema.Keys = append(ema.Keys, &model.SingleMovingAverageStamp{Timestamp: candlesticks[len(candlesticks)-1].OpenTime, Value: newEma})

	ema.Keys = ema.Keys[1:]

	return candlesticks, ema

}
