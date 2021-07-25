package utilities

import (
	"github.com/CraZzier/bot/model"
)

//MergeMACD adds macd to mycandles
func MergeMACD(candlesticks []*model.MyKline, macd []*model.MACD, indexstart int, indexstop int) []*model.MyKline {
	//Variable containers
	for i := 0; i < len(candlesticks); i++ {
		candlesticks[i].MacD = make([][]float64, len(macd))
	}
	for o := range macd {
		//Iterator for minutes moving Averages
		it1m := 0
		it5m := 0
		it15m := 0
		it1h := 0
		it1d := 0
		for i := 0; i < len(candlesticks); i++ {
			//Reading default candlesticks data
			opentime := candlesticks[i].OpenTime
			switch macd[o].Interval {
			case "1m":
				candlesticks[i].MacD[o] = macd[o].Keys[it1m].Value
				it1m++
			case "5m":
				if opentime >= macd[o].Keys[it5m].Timestamp {
					if it5m+1 < len(macd[o].Keys) {
						if opentime >= macd[o].Keys[it5m+1].Timestamp {
							it5m++
						}
						candlesticks[i].MacD[o] = macd[o].Keys[it5m].Value
					} else {
						candlesticks[i].MacD[o] = macd[o].Keys[it5m].Value
					}
				}
			case "15m":
				if opentime >= macd[o].Keys[it15m].Timestamp {
					if it15m+1 < len(macd[o].Keys) {
						if opentime >= macd[o].Keys[it15m+1].Timestamp {
							it15m++
						}
						candlesticks[i].MacD[o] = macd[o].Keys[it15m].Value
					} else {
						candlesticks[i].MacD[o] = macd[o].Keys[it15m].Value
					}
				}
			case "1h":
				if opentime >= macd[o].Keys[it1h].Timestamp {
					if it1h+1 < len(macd[o].Keys) {
						if opentime >= macd[o].Keys[it1h+1].Timestamp {
							it1h++
						}
						candlesticks[i].MacD[o] = macd[o].Keys[it1h].Value
					} else {
						candlesticks[i].MacD[o] = macd[o].Keys[it1h].Value
					}
				}
			case "1d":
				if opentime >= macd[o].Keys[it1d].Timestamp {
					if it1d+1 < len(macd[o].Keys) {
						if opentime >= macd[o].Keys[it1d+1].Timestamp {
							it1d++
						}
						candlesticks[i].MacD[o] = macd[o].Keys[it1d].Value
					} else {
						candlesticks[i].MacD[o] = macd[o].Keys[it1d].Value
					}
				}

			}
		}
	}

	return candlesticks
}

//MergeEMA mergers ema values with model.MyKline
func MergeEMA(candlesticks []*model.MyKline, emas []*model.MovingAverage, indexstart int, indexstop int) []*model.MyKline {
	for o := 0; o < len(emas); o++ {
		//Iterator for minutes moving Averages
		it1m := 0
		it5m := 0
		it15m := 0
		it1h := 0
		it1d := 0
		for i := 0; i < len(candlesticks); i++ {
			opentime := candlesticks[i].OpenTime
			switch emas[o].Interval {
			case "1m":
				candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it1m].Value)
			case "5m":
				if opentime >= emas[o].Keys[it5m].Timestamp {
					if it5m < len(emas[o].Keys) {
						if opentime >= emas[o].Keys[it5m+1].Timestamp {
							it5m++
						}
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it5m].Value)
					} else {
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it5m].Value)
					}
				}
			case "15m":
				if opentime >= emas[o].Keys[it15m].Timestamp {
					if it15m+1 < len(emas[o].Keys) {
						if opentime >= emas[o].Keys[it15m+1].Timestamp {
							it15m++
						}
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it15m].Value)
					} else {
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it15m].Value)
					}
				}
			case "1h":
				if opentime >= emas[o].Keys[it1h].Timestamp {
					if it1h+1 < len(emas[o].Keys) {
						if opentime >= emas[o].Keys[it1h+1].Timestamp {
							it1h++
						}
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it1h].Value)
					} else {
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it1h].Value)
					}
				}
			case "1d":
				if opentime >= emas[o].Keys[it1d].Timestamp {
					if it1d+1 < len(emas[o].Keys) {
						if opentime >= emas[o].Keys[it1d+1].Timestamp {
							it1d++
						}
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it1d].Value)
					} else {
						candlesticks[i].Emas = append(candlesticks[i].Emas, emas[o].Keys[it1d].Value)
					}
				}
			}
		}
	}
	return candlesticks
}

//MergeRSI mergres RSI with MyKline candles
func MergeRSI(candlesticks []*model.MyKline, rsi []*model.RSI, indexstart int, indexstop int) []*model.MyKline {
	for x := 0; x < len(rsi); x++ {
		it1m := 0
		it5m := 0
		it15m := 0
		it1h := 0
		it1d := 0

		for i := 0; i < len(candlesticks); i++ {
			//Reading default candlesticks data
			opentime := candlesticks[i].OpenTime
			switch rsi[x].Interval {
			case "1m":
				candlesticks[i].RSI[x] = rsi[x].Keys[it1m].RSI
				it1m++
			case "5m":
				if opentime >= rsi[x].Keys[it5m].Timestamp {
					if it5m+1 < len(rsi[x].Keys) {
						if opentime >= rsi[x].Keys[it5m+1].Timestamp {
							it5m++
						}
						candlesticks[i].RSI[x] = rsi[x].Keys[it5m].RSI
					} else {
						candlesticks[i].RSI[x] = rsi[x].Keys[it5m].RSI
					}
				}
			case "15m":
				if opentime >= rsi[x].Keys[it15m].Timestamp {
					if it15m+1 < len(rsi[x].Keys) {
						if opentime >= rsi[x].Keys[it15m+1].Timestamp {
							it15m++
						}
						candlesticks[i].RSI[x] = rsi[x].Keys[it15m].RSI
					} else {
						candlesticks[i].RSI[x] = rsi[x].Keys[it15m].RSI
					}
				}
			case "1h":
				if opentime >= rsi[x].Keys[it1h].Timestamp {
					if it1h+1 < len(rsi[x].Keys) {
						if opentime >= rsi[x].Keys[it1h+1].Timestamp {
							it1h++
						}
						candlesticks[i].RSI[x] = rsi[x].Keys[it1h].RSI
					} else {
						candlesticks[i].RSI[x] = rsi[x].Keys[it1h].RSI
					}
				}
			case "1d":
				if opentime >= rsi[x].Keys[it1d].Timestamp {
					if it1d+1 < len(rsi[x].Keys) {
						if opentime >= rsi[x].Keys[it1d+1].Timestamp {
							it1d++
						}
						candlesticks[i].RSI[x] = rsi[x].Keys[it1d].RSI
					} else {
						candlesticks[i].RSI[x] = rsi[x].Keys[it1d].RSI
					}
				}

			}
		}
	}
	return candlesticks
}

//MergeRSI mergres RSI with MyKline candles
func MergeBollingerBands(candlesticks []*model.MyKline, bb []*model.BollingerBands, indexstart int, indexstop int) []*model.MyKline {
	for x := 0; x < len(bb); x++ {
		it1m := 0
		it5m := 0
		it15m := 0
		it1h := 0
		it1d := 0

		for i := 0; i < len(candlesticks); i++ {
			//Reading default candlesticks data
			opentime := candlesticks[i].OpenTime
			switch bb[x].Interval {
			case "1m":
				candlesticks[i].BollingerBands[x] = bb[x].Keys[it1m].Value
				it1m++
			case "5m":
				if opentime >= bb[x].Keys[it5m].Timestamp {
					if it5m+1 < len(bb[x].Keys) {
						if opentime >= bb[x].Keys[it5m+1].Timestamp {
							it5m++
						}
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it5m].Value
					} else {
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it5m].Value
					}
				}
			case "15m":
				if opentime >= bb[x].Keys[it15m].Timestamp {
					if it15m+1 < len(bb[x].Keys) {
						if opentime >= bb[x].Keys[it15m+1].Timestamp {
							it15m++
						}
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it15m].Value
					} else {
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it15m].Value
					}
				}
			case "1h":
				if opentime >= bb[x].Keys[it1h].Timestamp {
					if it1h+1 < len(bb[x].Keys) {
						if opentime >= bb[x].Keys[it1h+1].Timestamp {
							it1h++
						}
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it1h].Value
					} else {
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it1h].Value
					}
				}
			case "1d":
				if opentime >= bb[x].Keys[it1d].Timestamp {
					if it1d+1 < len(bb[x].Keys) {
						if opentime >= bb[x].Keys[it1d+1].Timestamp {
							it1d++
						}
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it1d].Value
					} else {
						candlesticks[i].BollingerBands[x] = bb[x].Keys[it1d].Value
					}
				}

			}
		}
	}
	return candlesticks
}
