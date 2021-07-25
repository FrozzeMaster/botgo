package utilities

import (
	"strconv"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
)

//ToMyKline converting all ema into candlesticks
func ToMyKline(candlesticks []*futures.Kline, indexstart int, indexstop int) []*model.MyKline {

	//Variable containers
	var mycandlesticks []*model.MyKline
	for i := indexstart; i < indexstop; i++ {
		//Reading default candlesticks data
		var tempKline model.MyKline
		//Converting from string to int
		close, _ := strconv.ParseFloat(candlesticks[i].Close, 32)
		open, _ := strconv.ParseFloat(candlesticks[i].Open, 32)
		min, _ := strconv.ParseFloat(candlesticks[i].Low, 32)
		max, _ := strconv.ParseFloat(candlesticks[i].High, 32)
		opentime := candlesticks[i].OpenTime
		closetime := candlesticks[i].CloseTime
		volume, _ := strconv.ParseFloat(candlesticks[i].Volume, 32)
		//Adding to Kline simple Values
		tempKline.Close = close
		tempKline.Open = open
		tempKline.Min = min
		tempKline.Max = max
		tempKline.OpenTime = opentime
		tempKline.CloseTime = closetime
		tempKline.Volume = volume
		mycandlesticks = append(mycandlesticks, &tempKline)
	}
	return mycandlesticks
}

//CandlesToMyCandles converting all ema into candlesticks
func ToMyKlineSingle(candlesticks *futures.Kline) *model.MyKline {

	//Reading default candlesticks data
	var tempKline model.MyKline
	//Converting from string to int
	close, _ := strconv.ParseFloat(candlesticks.Close, 32)
	open, _ := strconv.ParseFloat(candlesticks.Open, 32)
	min, _ := strconv.ParseFloat(candlesticks.Low, 32)
	max, _ := strconv.ParseFloat(candlesticks.High, 32)
	opentime := candlesticks.OpenTime
	closetime := candlesticks.CloseTime
	volume, _ := strconv.ParseFloat(candlesticks.Volume, 32)
	//Adding to Kline simple Values
	tempKline.Close = close
	tempKline.Open = open
	tempKline.Min = min
	tempKline.Max = max
	tempKline.OpenTime = opentime
	tempKline.CloseTime = closetime
	tempKline.Volume = volume

	return &tempKline
}

//Period1Values returns key values from a period close min close max
func Period1ValuesBody(candlesticks []*model.MyKline, MinOrMax string, BodyOrTail string, macdNum int) (candle float64) {
	var mCandle float64
	switch MinOrMax {
	case "min":
		mCandle = 1000000
		switch BodyOrTail {
		case "body":
			for _, v := range candlesticks {
				if v.Close < mCandle {
					mCandle = v.Close
				}
			}

		case "tail":
			for _, v := range candlesticks {
				if v.Close < mCandle {
					mCandle = v.Min
				}
			}

		case "macd":
			for _, v := range candlesticks {
				if v.MacD[macdNum][0] < mCandle {
					mCandle = v.MacD[macdNum][0]
				}
			}
		}
	case "max":
		mCandle = 0
		switch BodyOrTail {
		case "body":
			for _, v := range candlesticks {
				if v.Close > mCandle {
					mCandle = v.Close
				}
			}
		case "macd":
			for _, v := range candlesticks {
				if v.MacD[macdNum][0] > mCandle {
					mCandle = v.MacD[macdNum][0]
				}
			}
		case "tail":
			for _, v := range candlesticks {
				if v.Close > mCandle {
					mCandle = v.Max
				}
			}

		}
	}
	return mCandle
}

//Checking crossing and returning index of cross
func CheckCrossing(candlesticks []*model.MyKline, toSide string, macdNum int) (bool, int) {
	found := 0
	index := -1
	for i := 1; i < len(candlesticks); i++ {
		if candlesticks[i].MacD[macdNum][1] > 0 && candlesticks[i-1].MacD[macdNum][1] <= 0 && toSide == "Up" {
			found++
			index = i
		}
		if candlesticks[i].MacD[macdNum][1] <= 0 && candlesticks[i-1].MacD[macdNum][1] > 0 && toSide == "Down" {
			found++
			index = i
		}
	}
	if found == 0 {
		return false, -1
	} else {
		return true, index
	}
}

//ConvertWsKlineToKline convert Wskline to Kline
func ConvertWsKlineToKline(event futures.WsKline) futures.Kline {
	var temp futures.Kline
	temp.OpenTime = event.StartTime
	temp.Open = event.Open
	temp.High = event.High
	temp.Low = event.Low
	temp.Close = event.Close
	temp.Volume = event.Volume
	temp.CloseTime = event.EndTime
	temp.QuoteAssetVolume = event.QuoteVolume
	temp.TradeNum = event.TradeNum
	temp.TakerBuyBaseAssetVolume = event.ActiveBuyVolume
	temp.TakerBuyQuoteAssetVolume = event.ActiveBuyQuoteVolume
	return temp
}
