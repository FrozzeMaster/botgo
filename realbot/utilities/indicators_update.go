package utilities

import (
	"fmt"
	"math"

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

//UpdateMACD updates MACD using the last given candle
func UpdateRSI(candlesticks []*model.MyKline, rsi *model.RSI) ([]*model.MyKline, *model.RSI) {

	//Calculating new values
	var candleVal float64
	timestamp := candlesticks[len(candlesticks)-1].OpenTime
	switch rsi.WhichValue {
	case "close":
		candleVal = candlesticks[len(candlesticks)-1].Close
	case "open":
		candleVal = candlesticks[len(candlesticks)-1].Open
	}
	var singleValue model.RSIstamp
	singleValue.Timestamp = timestamp
	singleValue.Close = candleVal
	singleValue.Change = singleValue.Close - rsi.Keys[len(rsi.Keys)-1].Close
	if singleValue.Change > 0 {
		singleValue.CurrGain = singleValue.Change
		singleValue.CurrLoss = 0
	} else {
		singleValue.CurrGain = 0
		singleValue.CurrLoss = -singleValue.Change
	}

	avggain, avgloss := rsi.Keys[len(rsi.Keys)-1].AvgGain, rsi.Keys[len(rsi.Keys)-1].AvgLoss

	singleValue.AvgGain = (avggain*float64(rsi.IntervalValue-1) + singleValue.CurrGain) / float64(rsi.IntervalValue)
	singleValue.AvgLoss = (avgloss*float64(rsi.IntervalValue-1) + singleValue.CurrLoss) / float64(rsi.IntervalValue)
	if singleValue.AvgLoss == 0 {
		singleValue.RS = 100
		singleValue.RSI = 100
	} else {
		singleValue.RS = singleValue.AvgGain / singleValue.AvgLoss
		singleValue.RSI = 100 - (100 / (1 + singleValue.RS))
	}
	rsi.Keys = append(rsi.Keys, &singleValue)

	candlesticks[len(candlesticks)-1].RSI = append(candlesticks[len(candlesticks)-1].RSI, singleValue.RSI)

	//Removing first element of an arrays to clear memory
	rsi.Keys = rsi.Keys[1:]

	return candlesticks, rsi

}

//UpdateBollingerBands
func UpdateBollingerBands(candlesticks []*model.MyKline, bb *model.BollingerBands) ([]*model.MyKline, *model.BollingerBands) {
	candleVal := candlesticks[len(candlesticks)-1].Close
	candleValForSma := candlesticks[len(candlesticks)-1-int(bb.E1)].Close
	ai := len(candlesticks) - 1
	//Updating SMA inside BOLLINGER
	lastSmaVal := bb.E1Keys[len(bb.E1Keys)-1].Value
	newSmaVal := (float64(bb.E1)*lastSmaVal - candleValForSma + candleVal) / float64(bb.E1)
	stDevParts := 0.00
	for i := ai; i >= ai-int(bb.E1)+1; i-- {
		stDevParts += math.Pow((candlesticks[i].Close - float64(newSmaVal)), 2)
	}
	stDev := stDevParts / (float64(bb.E1))
	stDev = math.Sqrt(stDev)
	upperBand := newSmaVal + bb.BandValue*stDev
	lowerBand := newSmaVal - bb.BandValue*stDev
	Bstamp := &model.SingleBollingerBandsStamp{
		Timestamp: candlesticks[ai].OpenTime,
		Value:     []float64{upperBand, newSmaVal, lowerBand},
	}
	SmaStamp := &model.SingleMovingAverageStamp{
		Timestamp: candlesticks[ai].OpenTime,
		Value:     newSmaVal,
	}
	bb.Keys = append(bb.Keys, Bstamp)
	bb.E1Keys = append(bb.E1Keys, SmaStamp)
	candlesticks[ai].BollingerBands = append(candlesticks[ai].BollingerBands, []float64{upperBand, newSmaVal, lowerBand})
	//Removing first element of an arrays to clear memory
	bb.Keys = bb.Keys[1:]
	bb.E1Keys = bb.E1Keys[1:]

	return candlesticks, bb
}

//UpdateATR
func UpdateATR(candlesticks []*model.MyKline, atr *model.ATR) ([]*model.MyKline, *model.ATR) {
	ai := len(candlesticks) - 1
	stDevParts := 0.00
	for i := ai; i >= ai-int(atr.ATRValue)+1; i-- {
		stDevParts += (candlesticks[i].Max - candlesticks[i].Min)
	}
	atrVal := stDevParts / (float64(atr.ATRValue))
	ATRstamp := &model.SingleATRStamp{
		Timestamp: candlesticks[ai].OpenTime,
		Value:     atrVal,
	}
	atr.Keys = append(atr.Keys, ATRstamp)
	candlesticks[ai].ATR = append(candlesticks[ai].ATR, atrVal)
	//Removing first element of an arrays to clear memory
	atr.Keys = atr.Keys[1:]

	return candlesticks, atr
}
