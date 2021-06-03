package utilities

import (
	"strconv"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
)

//GetStartStopCandles gets range of candlesticks to consider
func GetStartStopCandles(candlesticks []*futures.Kline, fromTimestamp int64, toTimestamp int64) (int, int) {
	indexstart, indexstop, startread, stopread := 0, len(candlesticks)-1, 0, 0

	//Getting indexes i can easily move over
	for i := 0; i < len(candlesticks); i++ {
		if candlesticks[i].OpenTime >= fromTimestamp && startread == 0 {
			indexstart = i
			startread = 1
		}
		if candlesticks[i].OpenTime >= toTimestamp && stopread == 0 {
			indexstop = i
			stopread = 1
			break
		}
	}
	return indexstart, indexstop
}

//StartingTransaction is made for handling enetering transactions in TestBot
func StartingTransaction(mycandlestick *model.MyKline, transNumbers *model.TransNumbers, usdtEntry float64, leverage float64, stopNumber float64, targetNumber float64, feeTaker float64, LongOrShort string) *model.AlgorithmTransactions {
	//Increasing Transactions amount
	transNumbers.TransAmount++

	//Calculating transactions value
	var tempTrans model.AlgorithmTransactions
	feeEntry := feeTaker * usdtEntry * leverage
	entryPrice := mycandlestick.Open
	stoploss := entryPrice * stopNumber
	target := entryPrice * targetNumber
	usdtEntryFee := usdtEntry - feeEntry
	sl := &model.Stoploss{
		Timestamp: mycandlestick.OpenTime,
		Value:     stoploss,
	}

	//Setting single Operation object
	tempTrans.StartBalance = usdtEntry
	tempTrans.EntryFee = feeEntry
	tempTrans.Stoploss = append(tempTrans.Stoploss, sl)
	tempTrans.Leverage = leverage
	tempTrans.BuyingPrice = entryPrice
	tempTrans.BalanceMinusFee = usdtEntryFee
	tempTrans.Target = target
	tempTrans.EntryTime = mycandlestick.OpenTime
	tempTrans.LongOrShort = LongOrShort

	//Transaction info
	// fmt.Printf("Starting transaction: usdt : %f, ", usdtEntry)
	// fmt.Printf("fee : %f, ", feeEntry)
	// fmt.Printf("entryPrice : %f, ", entryPrice)
	// fmt.Printf("stoploss : %f, ", stoploss)
	// fmt.Printf("target : %f, ", target)
	// fmt.Printf("time : %s \n", TimestampToDate(mycandlestick.OpenTime))

	//Returning
	return &tempTrans
}

//ClosingTransaction is made for handling closing tranactions in TestBot
func ClosingTransaction(mycandlestick *model.MyKline, transNumbers *model.TransNumbers, feeMaker float64, closePrice float64, tempTrans *model.AlgorithmTransactions) *model.AlgorithmTransactions {
	//Calculuating money after finish
	var usdt1 float64
	if tempTrans.LongOrShort == "Long" {
		usdt1 = tempTrans.BalanceMinusFee - (tempTrans.BalanceMinusFee * (tempTrans.BuyingPrice - closePrice) / tempTrans.BuyingPrice * tempTrans.Leverage)
	} else {
		usdt1 = tempTrans.BalanceMinusFee + (tempTrans.BalanceMinusFee * (tempTrans.BuyingPrice - closePrice) / tempTrans.BuyingPrice * tempTrans.Leverage)
	}
	feeM := usdt1 * feeMaker * tempTrans.Leverage
	usdtAfter := usdt1 - feeM
	if usdtAfter >= tempTrans.StartBalance {
		transNumbers.SuccessTrans++
		tempTrans.Type = "Succesful"
	} else {
		transNumbers.LostTrans++
		tempTrans.Type = "Lost"
	}
	tempTrans.FinishBalance = usdtAfter
	tempTrans.ClosingFee = feeM
	tempTrans.ClosingTime = mycandlestick.OpenTime
	tempTrans.SellingPrice = closePrice
	tempTrans.FeeSum = feeM + tempTrans.EntryFee
	tempTrans.Profit = usdtAfter - tempTrans.StartBalance
	tempTrans.ProfitWithoutFee = usdtAfter - tempTrans.StartBalance + feeM + tempTrans.FeeSum
	tempTrans.Percent = (tempTrans.Profit / tempTrans.StartBalance) * 100
	transNumbers.FeeSum += tempTrans.FeeSum

	//Writing out transaction data to console
	// fmt.Printf("Finishing transaction: usdt : %f, ", usdtAfter)
	// fmt.Printf("Type : succesful, ")
	// fmt.Printf("feeM :  %f, ", feeM)
	// fmt.Printf("usdt1 :  %f, ", usdt1)
	// fmt.Printf("closePrice : %f, ", tempTrans.SellingPrice)
	// fmt.Printf("percenr : %f, ", tempTrans.Percent)
	// fmt.Printf("time : %s \n", TimestampToDate(mycandlestick.OpenTime))

	//Returning Transaction finishd data
	return tempTrans
}

//CandlesToMyCandles converting all ema into candlesticks
func CandlesToMyCandles(candlesticks []*futures.Kline, indexstart int, indexstop int) []*model.MyKline {
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
		if candlesticks[i].MacD[macdNum][1] > 0 && candlesticks[i-1].MacD[macdNum][1] <= 0 && toSide == "Down" {
			found++
			index = i
		}
		if candlesticks[i].MacD[macdNum][1] <= 0 && candlesticks[i-1].MacD[macdNum][1] > 0 && toSide == "Up" {
			found++
			index = i
		}
	}
	if found == 0 {
		return false, -1
	} else if found == 1 {
		return true, index
	} else {
		return false, -1
	}
}

//MinAndMax gets the highest points on the chart and macd
func MinAndMax(array []int, mom string) int {
	switch mom {
	case "min":
		min := array[0]
		for i := 1; i < len(array); i++ {
			if min > array[i] {
				min = array[i]
			}
		}
		return min
	case "max":
		max := array[0]
		for i := 1; i < len(array); i++ {
			if max < array[i] {
				max = array[i]
			}
		}
		return max
	}
	return 0
}
