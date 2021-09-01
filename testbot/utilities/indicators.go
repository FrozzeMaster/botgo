package utilities

import (
	"log"
	"math"
	"strconv"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/kr/pretty"
)

//----Zamienianie na indicator w testbocie - konwersja strongow w candlestickach + dodanie indexstart stop i mili interval
//----Merge są takie same
//MovingAverage returns moving average struct
func MovingAverage(candlesticks []*futures.Kline, whichValue string, interval string, intervalValue int64, startTimestamp int64, stopTimestamp int64, pair string) *model.MovingAverage {
	//--------Constant SECTION----------------------------------------
	firstAvailableTimestamp := candlesticks[0].OpenTime
	var miliInterval int64
	//Checking range if it doestn exceed the candlesticks
	switch interval {
	case "1m":
		miliInterval = 60000 * intervalValue
	case "5m":
		miliInterval = 60000 * 5 * intervalValue
	case "15m":
		miliInterval = 60000 * 15 * intervalValue
	case "1h":
		miliInterval = 60000 * 60 * intervalValue
	case "1d":
		miliInterval = 60000 * 60 * 24 * intervalValue
	case "1w":
		miliInterval = 60000 * 60 * 24 * 7 * intervalValue
	default:
		miliInterval = 60000 * intervalValue
	}
	if firstAvailableTimestamp >= (startTimestamp - miliInterval) {
		log.Fatal("Not in range")
	}
	indexstart, indexstop := GetStartStopCandles(candlesticks, startTimestamp, stopTimestamp)
	//-------------End of constants section----------------------------
	//Declaring variables needed to return moving Average Object
	var fullMoving model.MovingAverage
	fullMoving.StartTimestamp = startTimestamp
	fullMoving.StopTimestamp = stopTimestamp
	fullMoving.Pair = pair
	fullMoving.Interval = interval
	//Calculating SMA for index
	for i := indexstart; i < indexstop; i++ {
		if whichValue == "open" {
			if i == indexstart {
				sum := 0.0
				//Going backwards to get average
				for p := i - int(intervalValue) + 1; p <= i; p++ {
					temp, _ := strconv.ParseFloat(candlesticks[p].Open, 32)
					sum += temp
				}
				score := sum / float64(intervalValue)
				var singleValue model.SingleMovingAverageStamp
				singleValue.Timestamp = candlesticks[i].OpenTime
				singleValue.Value = score
				fullMoving.Keys = append(fullMoving.Keys, &singleValue)
			} else {
				keysLength := len(fullMoving.Keys)
				temp, _ := strconv.ParseFloat(candlesticks[i-int(intervalValue)+1].Open, 32)
				temp1, _ := strconv.ParseFloat(candlesticks[i].Open, 32)
				sum := fullMoving.Keys[keysLength-1].Value
				sum = sum * float64(intervalValue)
				score := (sum - temp + temp1) / float64(intervalValue)
				var singleValue model.SingleMovingAverageStamp
				singleValue.Timestamp = candlesticks[i].OpenTime
				singleValue.Value = score
				fullMoving.Keys = append(fullMoving.Keys, &singleValue)
			}

		} else {
			if i == indexstart {
				sum := 0.0
				//Going backwards to get average
				for p := i - int(intervalValue) + 1; p <= i; p++ {
					temp, _ := strconv.ParseFloat(candlesticks[p].Close, 32)
					sum += temp
				}
				score := sum / float64(intervalValue)
				var singleValue model.SingleMovingAverageStamp
				singleValue.Timestamp = candlesticks[i].OpenTime
				singleValue.Value = score
				fullMoving.Keys = append(fullMoving.Keys, &singleValue)
			} else {
				keysLength := len(fullMoving.Keys)
				temp, _ := strconv.ParseFloat(candlesticks[i-int(intervalValue)+1].Close, 32)
				temp1, _ := strconv.ParseFloat(candlesticks[i].Close, 32)
				sum := fullMoving.Keys[keysLength-1].Value
				sum = sum * float64(intervalValue)
				score := (sum - temp + temp1) / float64(intervalValue)
				var singleValue model.SingleMovingAverageStamp
				singleValue.Timestamp = candlesticks[i].OpenTime
				singleValue.Value = score
				fullMoving.Keys = append(fullMoving.Keys, &singleValue)
			}
		}
	}
	return &fullMoving
}

//RSI returns moving rsi struct
func RSI(candlesticks []*futures.Kline, whichValue string, interval string, intervalValue int64, startTimestamp int64, stopTimestamp int64, pair string) *model.RSI {
	//--------Constant SECTION----------------------------------------
	firstAvailableTimestamp := candlesticks[0].OpenTime
	var miliInterval int64
	//Checking range if it doestn exceed the candlesticks
	switch interval {
	case "1m":
		miliInterval = 60000 * intervalValue
	case "5m":
		miliInterval = 60000 * 5 * intervalValue
	case "15m":
		miliInterval = 60000 * 15 * intervalValue
	case "1h":
		miliInterval = 60000 * 60 * intervalValue
	case "1d":
		miliInterval = 60000 * 60 * 24 * intervalValue
	case "1w":
		miliInterval = 60000 * 60 * 24 * 7 * intervalValue
	default:
		miliInterval = 60000 * intervalValue
	}
	if firstAvailableTimestamp >= (startTimestamp - miliInterval) {
		log.Fatal("Not in range")
	}
	indexstart, indexstop := GetStartStopCandles(candlesticks, startTimestamp, stopTimestamp)
	//-------------End of constants section----------------------------
	//Declaring variables needed to return moving Average Object
	var fullMoving model.RSI
	fullMoving.StartTimestamp = startTimestamp
	fullMoving.StopTimestamp = stopTimestamp
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
			singleValue.Close, _ = strconv.ParseFloat(candlesticks[i].Close, 32)
			fullMoving.Keys = append(fullMoving.Keys, &singleValue)
		} else {
			var singleValue model.RSIstamp
			singleValue.Timestamp = candlesticks[i].OpenTime
			singleValue.Close, _ = strconv.ParseFloat(candlesticks[i].Close, 32)
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
			singleValue.Close, _ = strconv.ParseFloat(candlesticks[i].Close, 32)
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
			singleValue.Close, _ = strconv.ParseFloat(candlesticks[i].Close, 32)
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

//EmovingAverage retursn Exponential moving average
func EmovingAverage(candlesticks []*futures.Kline, whichValue string, interval string, intervalValue int64, startTimestamp int64, stopTimestamp int64, pair string) *model.MovingAverage {
	//--------Constant SECTION----------------------------------------
	firstAvailableTimestamp := candlesticks[0].OpenTime
	var miliInterval int64
	//Checking range if it doestn exceed the candlesticks
	switch interval {
	case "1m":
		miliInterval = 60000 * intervalValue
	case "5m":
		miliInterval = 60000 * 5 * intervalValue
	case "15m":
		miliInterval = 60000 * 15 * intervalValue
	case "1h":
		miliInterval = 60000 * 60 * intervalValue
	case "1d":
		miliInterval = 60000 * 60 * 24 * intervalValue
	case "1w":
		miliInterval = 60000 * 60 * 24 * 7 * intervalValue
	default:
		miliInterval = 60000 * intervalValue
	}
	if firstAvailableTimestamp >= (startTimestamp - miliInterval) {
		pretty.Println(firstAvailableTimestamp)
		pretty.Println(startTimestamp)
		pretty.Println(miliInterval)
		pretty.Println(startTimestamp - miliInterval)
		log.Fatal("Not in range", interval)
	}
	indexstart, indexstop := GetStartStopCandles(candlesticks, startTimestamp, stopTimestamp)
	//-------------End of constants section----------------------------
	//Getting sma data and setting first EMA
	smatable := MovingAverage(candlesticks, whichValue, interval, intervalValue, startTimestamp, stopTimestamp, pair)
	var Ema model.MovingAverage
	var EmaStamps []*model.SingleMovingAverageStamp
	var FirstEmaStamp model.SingleMovingAverageStamp
	FirstEmaStamp.Timestamp = startTimestamp
	FirstEmaStamp.Value = smatable.Keys[0].Value
	margin := float64(2.00 / (1.00 + float64(intervalValue)))
	EmaStamps = append(EmaStamps, &FirstEmaStamp)
	//Getting indexes i can easily move over

	//Calculating EMA for index
	it := 0
	for i := indexstart + 1; i < indexstop; i++ {
		if whichValue == "open" {
			ema := 0.00
			//Going backwards to get average
			open, _ := strconv.ParseFloat(candlesticks[i].Open, 32)
			ema = open*margin + EmaStamps[it].Value*(1.00-margin)
			it++
			var SEMA model.SingleMovingAverageStamp
			SEMA.Timestamp = candlesticks[i].OpenTime
			SEMA.Value = ema
			EmaStamps = append(EmaStamps, &SEMA)

		} else {
			ema := 0.00
			//Going backwards to get average
			close, _ := strconv.ParseFloat(candlesticks[i].Close, 32)
			ema = close*margin + EmaStamps[it].Value*(1.00-margin)
			it++
			var SEMA model.SingleMovingAverageStamp
			SEMA.Timestamp = candlesticks[i].OpenTime
			SEMA.Value = ema
			EmaStamps = append(EmaStamps, &SEMA)
		}
	}
	Ema.Keys = EmaStamps
	Ema.Pair = pair
	Ema.StartTimestamp = startTimestamp
	Ema.StopTimestamp = stopTimestamp
	Ema.Interval = interval
	return &Ema
}

//MACD returns MACD
func MACD(candlesticks []*futures.Kline, whichValue string, signalValue1 int64, intervalValue1 int64, intervalValue2 int64, interval string, startTimestamp int64, stopTimestamp int64, pair string) *model.MACD {
	signalValue := float64(signalValue1)
	var macd model.MACD
	macd.CandleTrueInterval = (candlesticks[1].OpenTime - candlesticks[0].OpenTime) / 60000
	macd.Interval = interval
	macd.Pair = pair
	macd.StartTimestamp = startTimestamp
	macd.StopTimestamp = stopTimestamp
	margin := float64(2.00 / (1.00 + float64(signalValue)))
	//Getting sma data and setting first EMA
	ema1 := EmovingAverage(candlesticks, whichValue, interval, intervalValue1, startTimestamp, stopTimestamp, pair)
	ema2 := EmovingAverage(candlesticks, whichValue, interval, intervalValue2, startTimestamp, stopTimestamp, pair)
	var macdVal []float64
	var signalVal []float64
	var timeStamp []int64
	for i := 0; i < len(ema1.Keys); i++ {
		macdVal = append(macdVal, ema1.Keys[i].Value-ema2.Keys[i].Value)
		timeStamp = append(timeStamp, ema1.Keys[i].Timestamp)
	}
	var sumSm float64 = 0
	for i := 0; i < 0+int(signalValue); i++ {
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
			singleKey.Value = append(singleKey.Value, 0)
			singleKey.Value = append(singleKey.Value, 0)
		} else {
			singleKey.Value = append(singleKey.Value, signalVal[i-int(signalValue)+1])
			singleKey.Value = append(singleKey.Value, macdVal[i]-signalVal[i-int(signalValue)+1])
		}
		singleKey.Timestamp = timeStamp[i]
		macd.Keys = append(macd.Keys, &singleKey)
	}
	return &macd
}

//BollingerBands
func BollingerBands(candlesticks []*futures.Kline, smaValue int64, interval string, bandValue float64, startTimestamp int64, stopTimestamp int64, pair string) *model.BollingerBands {
	//--------Constant SECTION----------------------------------------
	firstAvailableTimestamp := candlesticks[0].OpenTime
	var miliInterval int64
	//Checking range if it doestn exceed the candlesticks
	switch interval {
	case "1m":
		miliInterval = 60000 * smaValue
	case "5m":
		miliInterval = 60000 * 5 * smaValue
	case "15m":
		miliInterval = 60000 * 15 * smaValue
	case "1h":
		miliInterval = 60000 * 60 * smaValue
	case "1d":
		miliInterval = 60000 * 60 * 24 * smaValue
	case "1w":
		miliInterval = 60000 * 60 * 24 * 7 * smaValue
	default:
		miliInterval = 60000 * smaValue
	}
	if firstAvailableTimestamp >= (startTimestamp - miliInterval) {
		log.Fatal("Not in range")
	}
	indexstart, indexstop := GetStartStopCandles(candlesticks, startTimestamp, stopTimestamp)
	//-------------End of constants session----------------------------
	//Getting sma data and declaring containers
	smatable := MovingAverage(candlesticks, "close", interval, smaValue, int64(indexstart), int64(indexstop), pair)
	smaiterator := 0
	var BollingerBands model.BollingerBands
	var BollingerBandsStamps []*model.SingleBollingerBandsStamp

	//Calculating Bollinger for index
	for i := indexstart + int(smaValue); i < indexstop; i++ {
		actualsma := smatable.Keys[smaiterator].Value
		stDevParts := 0.00
		//Going backwards to get average
		for o := i; o >= i-int(smaValue)+1; o-- {
			close, _ := strconv.ParseFloat(candlesticks[o].Close, 32)
			stDevParts += math.Pow((close - actualsma), 2)
		}
		stDev := stDevParts / (float64(smaValue))
		stDev = math.Sqrt(stDev)
		upperBand := actualsma + bandValue*stDev
		lowerBand := actualsma - bandValue*stDev
		Bstamp := &model.SingleBollingerBandsStamp{
			Timestamp: candlesticks[i].OpenTime,
			Value:     []float64{upperBand, actualsma, lowerBand},
		}
		Bstamp.Timestamp = candlesticks[i].OpenTime
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

//ATR
func ATR(candlesticks []*futures.Kline, ATRValue int64, interval string, startTimestamp int64, stopTimestamp int64, pair string) *model.ATR {
	//--------Constant SECTION----------------------------------------
	firstAvailableTimestamp := candlesticks[0].OpenTime
	var miliInterval int64
	//Checking range if it doestn exceed the candlesticks
	switch interval {
	case "1m":
		miliInterval = 60000 * ATRValue
	case "5m":
		miliInterval = 60000 * 5 * ATRValue
	case "15m":
		miliInterval = 60000 * 15 * ATRValue
	case "1h":
		miliInterval = 60000 * 60 * ATRValue
	case "1d":
		miliInterval = 60000 * 60 * 24 * ATRValue
	case "1w":
		miliInterval = 60000 * 60 * 24 * 7 * ATRValue
	default:
		miliInterval = 60000 * ATRValue
	}
	if firstAvailableTimestamp >= (startTimestamp - miliInterval) {
		log.Fatal("Not in range")
	}
	indexstart, indexstop := GetStartStopCandles(candlesticks, startTimestamp, stopTimestamp)
	//-------------End of constants session----------------------------
	var ATR model.ATR
	var ATRStamps []*model.SingleATRStamp

	//Calculating Bollinger for index
	for i := indexstart + int(ATRValue) - 1; i < indexstop; i++ {

		stDevParts := 0.00
		//Going backwards to get average
		for o := i; o >= i-int(ATRValue)+1; o-- {
			mx, _ := strconv.ParseFloat(candlesticks[o].High, 32)
			mn, _ := strconv.ParseFloat(candlesticks[o].Low, 32)
			stDevParts += (mx - mn)
		}
		atr := stDevParts / (float64(ATRValue))

		ATRstamp := &model.SingleATRStamp{
			Timestamp: candlesticks[i].OpenTime,
			Value:     atr,
		}
		ATRStamps = append(ATRStamps, ATRstamp)

	}
	//Model data
	ATR.Keys = ATRStamps
	ATR.Pair = pair
	ATR.StartTimestamp = int64(indexstart)
	ATR.StopTimestamp = int64(indexstop)
	ATR.Interval = interval
	ATR.ATRValue = ATRValue
	return &ATR
}
