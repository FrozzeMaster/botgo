package utilities

import (
	"github.com/CraZzier/bot/model"
)

//MovingAverage returns moving average struct
func SMA(candlesticks []*model.MyKline, whichValue string, interval string, intervalValue int64, indexstart int, indexstop int, pair string) *model.MovingAverage {

	//Declaring variables needed to return moving Average Object
	var fullMoving model.MovingAverage
	fullMoving.StartTimestamp = int64(indexstart)
	fullMoving.StopTimestamp = int64(indexstop)
	fullMoving.Pair = pair
	fullMoving.Interval = interval
	//Checking if index is in range
	if indexstart >= indexstop {
		return nil
	}
	//Getting first SMA value
	sum := 0.0
	for i := indexstart; i < int(intervalValue); i++ {
		var temp float64
		switch whichValue {
		case "close":
			temp = candlesticks[i].Close
		case "open":
			temp = candlesticks[i].Open
		}
		sum += temp
		if i+1 == int(intervalValue) {
			var singleValue model.SingleMovingAverageStamp
			singleValue.Timestamp = candlesticks[i].OpenTime
			singleValue.Value = sum / float64(intervalValue)
			fullMoving.Keys = append(fullMoving.Keys, &singleValue)
		}
	}

	//Calculating SMA for index
	for i := int(intervalValue); i < indexstop; i++ {
		var temp, temp1 float64
		switch whichValue {
		case "close":
			temp = candlesticks[i-int(intervalValue)].Close
			temp1 = candlesticks[i].Close
		case "open":
			temp = candlesticks[i-int(intervalValue)].Open
			temp1 = candlesticks[i].Open
		}
		keysLength := len(fullMoving.Keys)
		sum := fullMoving.Keys[keysLength-1].Value
		sum = sum * float64(intervalValue)
		score := (sum - temp + temp1) / float64(intervalValue)
		//Adding to Keys
		var singleValue model.SingleMovingAverageStamp
		singleValue.Timestamp = candlesticks[i].OpenTime
		singleValue.Value = score
		fullMoving.Keys = append(fullMoving.Keys, &singleValue)
	}

	return &fullMoving
}

//EmovingAverage retursn Exponential moving average - OK80%
func EMA(candlesticks []*model.MyKline, whichValue string, interval string, intervalValue int64, indexstart int, indexstop int, pair string) *model.MovingAverage {

	//Getting sma data and setting first EMA
	smatable := SMA(candlesticks, whichValue, interval, intervalValue, indexstart, indexstop, pair)
	var Ema model.MovingAverage
	var EmaStamps []*model.SingleMovingAverageStamp

	var FirstEmaStamp model.SingleMovingAverageStamp
	FirstEmaStamp.Timestamp = smatable.Keys[0].Timestamp
	FirstEmaStamp.Value = smatable.Keys[0].Value
	EmaStamps = append(EmaStamps, &FirstEmaStamp)

	margin := float64(2.00 / (1.00 + float64(intervalValue)))

	//Calculating EMA for index

	for i := indexstart + int(intervalValue); i < indexstop; i++ {
		var candleVal float64
		switch whichValue {
		case "close":
			candleVal = candlesticks[i].Close
		case "open":
			candleVal = candlesticks[i].Open
		}
		ema := 0.00
		//Going backwards to get average
		emaLength := len(EmaStamps)
		ema = candleVal*margin + EmaStamps[emaLength-1].Value*(1.00-margin)
		var SEMA model.SingleMovingAverageStamp
		SEMA.Timestamp = candlesticks[i].OpenTime
		SEMA.Value = ema
		EmaStamps = append(EmaStamps, &SEMA)
	}
	//Model data
	Ema.Keys = EmaStamps
	Ema.Pair = pair
	Ema.StartTimestamp = int64(indexstart)
	Ema.StopTimestamp = int64(indexstop)
	Ema.Interval = interval
	Ema.IntervalValue = intervalValue
	Ema.WhichValue = whichValue
	return &Ema
}

//MACD returns MACD -Ok 60%
func MACD(candlesticks []*model.MyKline, whichValue string, signalValue1 int64, intervalValue1 int64, intervalValue2 int64, interval string, indexstart int, indexstop int, pair string) *model.MACD {
	signalValue := float64(signalValue1)
	var macd model.MACD
	macd.CandleTrueInterval = (candlesticks[1].OpenTime - candlesticks[0].OpenTime) / 60000
	macd.Interval = interval
	macd.E1 = intervalValue1
	macd.E2 = intervalValue2
	macd.Signal = signalValue1
	macd.Interval = interval
	macd.Interval = interval
	macd.Pair = pair
	macd.WhichValue = whichValue
	macd.StartTimestamp = int64(indexstart)
	macd.StopTimestamp = int64(indexstop)
	margin := float64(2.00 / (1.00 + float64(signalValue)))
	//Getting sma data and setting first EMA
	ema1 := EMA(candlesticks, whichValue, interval, intervalValue1, indexstart, indexstop, pair)
	ema2 := EMA(candlesticks, whichValue, interval, intervalValue2, indexstart, indexstop, pair)
	var macdVal []float64
	var signalVal []float64
	var timeStamp []int64
	difEma := intervalValue2 - intervalValue1
	for i := 0; i < len(ema2.Keys); i++ {
		macdVal = append(macdVal, ema1.Keys[i+int(difEma)].Value-ema2.Keys[i].Value)
		timeStamp = append(timeStamp, ema1.Keys[i+int(difEma)].Timestamp)
		emastamp1 := &model.SingleMovingAverageStamp{Timestamp: ema1.Keys[i+int(difEma)].Timestamp, Value: ema1.Keys[i+int(difEma)].Value}
		emastamp2 := &model.SingleMovingAverageStamp{Timestamp: ema1.Keys[i+int(difEma)].Timestamp, Value: ema2.Keys[i].Value}
		macd.E1Keys = append(macd.E1Keys, emastamp1)
		macd.E2Keys = append(macd.E2Keys, emastamp2)
	}
	var sumSm float64 = 0
	for i := 0; i < int(signalValue); i++ {
		sumSm += macdVal[i]
	}
	signalVal = append(signalVal, sumSm/float64(signalValue))
	for i := 0; i < len(macdVal)-int(signalValue); i++ {
		result := signalVal[i]*(1-margin) + macdVal[i+int(signalValue)]*margin
		signalVal = append(signalVal, result)
	}
	for i := 0; i < len(macdVal); i++ {
		var singleKey model.SingleMACDStamp
		singleKey.Value = append(singleKey.Value, macdVal[i])
		if i < int(signalValue) {
			// singleKey.Value = append(singleKey.Value, 0)
			// singleKey.Value = append(singleKey.Value, 0)
		} else {
			singleKey.Value = append(singleKey.Value, signalVal[i-int(signalValue)+1])
			singleKey.Value = append(singleKey.Value, macdVal[i]-signalVal[i-int(signalValue)+1])
		}
		singleKey.Timestamp = timeStamp[i]
		macd.Keys = append(macd.Keys, &singleKey)
	}
	return &macd
}

//RSI returns moving rsi struct
func RSI(candlesticks []*model.MyKline, whichValue string, interval string, intervalValue int64, indexstart int, indexstop int, pair string) *model.RSI {

	//Declaring variables needed to return moving Average Object
	var fullMoving model.RSI
	fullMoving.StartTimestamp = int64(indexstart)
	fullMoving.StopTimestamp = int64(indexstop)
	fullMoving.Pair = pair
	fullMoving.Interval = interval
	//Liczenie pierwszych wartości
	for i := indexstart; i < indexstart+int(intervalValue); i++ {
		if i == indexstart {
			var singleValue model.RSIstamp
			singleValue.Timestamp = candlesticks[i].OpenTime
			singleValue.Change = 0
			singleValue.AvgGain = 0
			singleValue.AvgLoss = 0
			singleValue.CurrGain = 0
			singleValue.CurrLoss = 0
			singleValue.RS = 0
			singleValue.RSI = 0
			singleValue.Close = candlesticks[i].Close
			fullMoving.Keys = append(fullMoving.Keys, &singleValue)
		} else {
			var singleValue model.RSIstamp
			singleValue.Timestamp = candlesticks[i].OpenTime
			singleValue.Close = candlesticks[i].Close
			singleValue.Change = fullMoving.Keys[len(fullMoving.Keys)-1].Close - singleValue.Close
			if singleValue.Change > 0 {
				singleValue.CurrGain = singleValue.Change
				singleValue.CurrLoss = 0
			} else {
				singleValue.CurrGain = 0
				singleValue.CurrLoss = -singleValue.Change
			}
			singleValue.AvgGain = 0
			singleValue.AvgLoss = 0
			singleValue.RS = 0
			singleValue.RSI = 0
			fullMoving.Keys = append(fullMoving.Keys, &singleValue)
		}
	}
	//Calculating SMA for index
	for i := indexstart + int(intervalValue); i < indexstop; i++ {
		if i == indexstart+int(intervalValue) {
			var singleValue model.RSIstamp
			singleValue.Timestamp = candlesticks[i].OpenTime
			singleValue.Close = candlesticks[i].Close
			singleValue.Change = singleValue.Close - fullMoving.Keys[len(fullMoving.Keys)-1].Close
			if singleValue.Change > 0 {
				singleValue.CurrGain = singleValue.Change
				singleValue.CurrLoss = 0
			} else {
				singleValue.CurrGain = 0
				singleValue.CurrLoss = -singleValue.Change
			}

			avggain, avgloss := 0.00, 0.00
			avggain = avggain + singleValue.CurrGain
			avgloss = avgloss + singleValue.CurrLoss
			for o := len(fullMoving.Keys) - 1; o > len(fullMoving.Keys)-(int(intervalValue)); o-- {
				avggain += fullMoving.Keys[o].CurrGain
				avgloss += fullMoving.Keys[o].CurrLoss
			}
			singleValue.AvgGain = avggain / float64(intervalValue)
			singleValue.AvgLoss = avgloss / float64(intervalValue)
			if singleValue.AvgLoss == 0 {
				singleValue.RS = 100
				singleValue.RSI = 100
			} else {
				singleValue.RS = singleValue.AvgGain / singleValue.AvgLoss
				singleValue.RSI = 100 - (100 / (1 + singleValue.RS))
			}
			fullMoving.Keys = append(fullMoving.Keys, &singleValue)
		} else {
			var singleValue model.RSIstamp
			singleValue.Timestamp = candlesticks[i].OpenTime
			singleValue.Close = candlesticks[i].Close
			singleValue.Change = singleValue.Close - fullMoving.Keys[len(fullMoving.Keys)-1].Close
			if singleValue.Change > 0 {
				singleValue.CurrGain = singleValue.Change
				singleValue.CurrLoss = 0
			} else {
				singleValue.CurrGain = 0
				singleValue.CurrLoss = -singleValue.Change
			}

			avggain, avgloss := fullMoving.Keys[len(fullMoving.Keys)-1].AvgGain, fullMoving.Keys[len(fullMoving.Keys)-1].AvgLoss

			singleValue.AvgGain = (avggain*float64(intervalValue-1) + singleValue.CurrGain) / float64(intervalValue)
			singleValue.AvgLoss = (avgloss*float64(intervalValue-1) + singleValue.CurrLoss) / float64(intervalValue)
			if singleValue.AvgLoss == 0 {
				singleValue.RS = 100
				singleValue.RSI = 100
			} else {
				singleValue.RS = singleValue.AvgGain / singleValue.AvgLoss
				singleValue.RSI = 100 - (100 / (1 + singleValue.RS))
			}
			fullMoving.Keys = append(fullMoving.Keys, &singleValue)
		}
	}

	return &fullMoving
}

//BollingerBands retursn Exponential moving average - OK80%
func BollingerBands(candlesticks []*model.MyKline, smaValue int64, interval string, bandValue float64, indexstart int, indexstop int, pair string) *model.BollingerBands {

	//Getting sma data and declaring containers
	smatable := SMA(candlesticks, "close", interval, smaValue, indexstart, indexstop, pair)
	smaiterator := 0
	var BollingerBands model.BollingerBands
	var BollingerBandsStamps []*model.SingleBollingerBandsStamp

	//Calculating Bollinger for index
	for i := indexstart + int(smaValue); i < indexstop; i++ {

		candleVal := candlesticks[i].Close
		actualsma := smatable.Keys[smaiterator].Value
		stDevParts := 0.00
		//Going backwards to get average
		for o := i - int(smaValue); o < i; o++ {
			stDevParts += (candlesticks[o].Close - actualsma)
		}
		stDev := (stDevParts + (candleVal - actualsma)) / float64((float64(smaValue) - 1))
		upperBand := actualsma + bandValue*stDev
		lowerBand := actualsma - bandValue*stDev
		Bstamp := &model.SingleBollingerBandsStamp{
			Timestamp: candlesticks[i].OpenTime,
			Value:     []float64{upperBand, actualsma, lowerBand},
		}
		BollingerBandsStamps = append(BollingerBandsStamps, Bstamp)

	}
	//Model data
	BollingerBands.Keys = BollingerBandsStamps
	BollingerBands.E1Keys = smatable.Keys
	BollingerBands.Pair = pair
	BollingerBands.StartTimestamp = int64(indexstart)
	BollingerBands.StopTimestamp = int64(indexstop)
	BollingerBands.Interval = interval
	BollingerBands.BandValue = bandValue
	BollingerBands.E1 = smaValue
	return &BollingerBands
}
