package testbot

import (
	"fmt"
	"math"

	"github.com/CraZzier/bot/model"
	"github.com/CraZzier/bot/testbot/utilities"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/kr/pretty"
)

//AlgoMacdSP is algorithm MACD with SINGLE PAIR
func AlgoMacdSP(usdtEntry float64, feeMaker float64, feeTaker float64, candlesticks [][][]*futures.Kline, from string, to string, pair string) *model.AlgorithmTransactionsFull {
	var pairNum int
	switch pair {
	case "BTCUSDT":
		pairNum = 0
	case "ETHUSDT":
		pairNum = 1
	case "ADAUSDT":
		pairNum = 2
	case "NEOUSDT":
		pairNum = 3
	case "XLMUSDT":
		pairNum = 4
	case "XRPUSDT":
		pairNum = 5
	case "LINKUSDT":
		pairNum = 6
	}
	//Declaring important variables
	var it int
	transNumbers := &model.TransNumbers{
		TransAmount:  0,
		SuccessTrans: 0,
		LostTrans:    0,
		FeeSum:       0.00,
	}
	leverage := 10.00
	profit, stop := 1.01, 0.9975
	//Preparation for checking start index and stop index
	fromTimestamp, toTimestamp := DateToTimestampRange(from), DateToTimestampRange(to)
	indexstart, indexstop := utilities.GetStartStopCandles(candlesticks[pairNum][1], fromTimestamp, toTimestamp)

	//Container for Trasnactions
	var transactions []*model.AlgorithmTransactions
	var mycandlesticks []*model.MyKline
	//Getting SMA and EMA
	ema1 := utilities.EmovingAverage(candlesticks[pairNum][2], "close", "15m", 50, fromTimestamp, toTimestamp, pair)
	ema2 := utilities.EmovingAverage(candlesticks[pairNum][3], "close", "1h", 50, fromTimestamp, toTimestamp, pair)
	macd := utilities.MACD(candlesticks[pairNum][1], "close", 18, 13, 25, "5m", fromTimestamp, toTimestamp, pair)
	var tableOfEMAs []*model.MovingAverage
	tableOfEMAs = append(tableOfEMAs, ema1, ema2)
	var tableOfMACDs []*model.MACD
	tableOfMACDs = append(tableOfMACDs, macd)
	mycandlesticks = utilities.CandlesToMyCandles(candlesticks[pairNum][1], indexstart, indexstop)
	mycandlesticks = utilities.MergeEMA(mycandlesticks, tableOfEMAs, indexstart, indexstop)
	mycandlesticks = utilities.MergeMACD(mycandlesticks, tableOfMACDs, indexstart, indexstop)
	mycandlesticks = mycandlesticks[0 : len(mycandlesticks)-2]
	it = indexstop - indexstart - 2
	//Going through candlesticks
	for i := 0; i < it-3; i++ {
		if mycandlesticks[i].MacD[0][1] >= 0 && mycandlesticks[i+1].MacD[0][1] < 0 && mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
			//Going over next candlesticks
			i = i + 2

			//Checking 1 PERIOD /////////////////////////////////////////////////////////////
			indexPeriod1 := i
			shouldStop := 0
			for mycandlesticks[i].MacD[0][2] < 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] > 0 || mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod1 < 3 || i-indexPeriod1 > 150 || shouldStop == 1 {
				continue
			}
			macMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod1:i], "min", "macd", 0)

			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for mycandlesticks[i].MacD[0][2] >= 0 && mycandlesticks[i].MacD[0][0] < 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] > 0 || mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod2 < 3 || i-indexPeriod2 > 40 && shouldStop == 1 {
				continue
			}

			//Checking 3 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod3 := i
			for mycandlesticks[i].MacD[0][2] < 0 && mycandlesticks[i].MacD[0][0] < 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] > 0 || mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod3 < 3 || i-indexPeriod3 > 150 {
				continue
			}
			candleMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "min", "body", 0)
			macMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "min", "macd", 0)
			candleMinPeriod3Tail := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "min", "tail", 0)
			var beg1Period int
			for iback := indexPeriod1; mycandlesticks[iback].MacD[0][2] < 0; iback-- {
				beg1Period = iback
			}
			candleMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[beg1Period:indexPeriod2-1], "min", "body", 0)
			pretty.Println(candleMinPeriod3, macMinPeriod3, candleMinPeriod3Tail, candleMinPeriod1, macMinPeriod1)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 <= candleMinPeriod3 || macMinPeriod1 >= macMinPeriod3 {
				continue
			} else {
				i++
			}
			//Starting transaction
			stop = candleMinPeriod3Tail / mycandlesticks[i].Open
			profit = 1 + 3*(1-stop)
			if stop < 1 {
				leverage = 0.05 / (1 - stop)
			} else {
				leverage = 0.05 / (-1 + stop)
			}
			leverage = math.Round(leverage)
			if stop >= 0.995 {
				fmt.Println("Za niski target")
				continue
			}
			transData := utilities.StartingTransaction(mycandlesticks[i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Long")
			//Simulating begin in trade
			for i < it-2 {
				if mycandlesticks[i].Min <= transData.Stoploss[0].Value {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				if transData.Target >= mycandlesticks[i].Min && transData.Target <= mycandlesticks[i].Max {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Target, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				i++
			}
			if usdtEntry <= 0 {
				break
			}
		}
		if mycandlesticks[i].MacD[0][1] <= 0 && mycandlesticks[i+1].MacD[0][1] > 0 && mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
			//Going over next candlesticks
			i = i + 2

			//Checking 1 PERIOD /////////////////////////////////////////////////////////////
			indexPeriod1 := i
			shouldStop := 0
			for mycandlesticks[i].MacD[0][2] > 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] < 0 || mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod1 < 3 || i-indexPeriod1 > 150 || shouldStop == 1 {
				continue
			}
			macMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod1:i], "max", "macd", 0)

			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for mycandlesticks[i].MacD[0][2] <= 0 && mycandlesticks[i].MacD[0][0] > 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] < 0 || mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod2 < 3 || i-indexPeriod2 > 40 && shouldStop == 1 {
				continue
			}

			//Checking 3 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod3 := i
			for mycandlesticks[i].MacD[0][2] > 0 && mycandlesticks[i].MacD[0][0] > 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] < 0 || mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod3 < 3 || i-indexPeriod3 > 150 {
				continue
			}
			candleMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "max", "body", 0)
			macMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "max", "macd", 0)
			candleMinPeriod3Tail := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "max", "tail", 0)
			var beg1Period int
			for iback := indexPeriod1; mycandlesticks[iback].MacD[0][2] > 0; iback-- {
				beg1Period = iback
			}
			candleMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[beg1Period:indexPeriod2-1], "max", "body", 0)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 >= candleMinPeriod3 || macMinPeriod1 <= macMinPeriod3 {
				continue
			} else {
				i++
			}
			//Starting transaction
			stop = candleMinPeriod3Tail / mycandlesticks[i].Open
			profit = 1 + 3*(1-stop)
			if stop < 1 {
				leverage = 0.05 / (1 - stop)
			} else {
				leverage = 0.05 / (-1 + stop)
			}
			leverage = math.Round(leverage)
			if stop <= 1.005 {
				//fmt.Println("Za niski target")
				continue
			}
			transData := utilities.StartingTransaction(mycandlesticks[i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Short")
			//Simulating begin in trade
			for i < it-2 {
				if mycandlesticks[i].Max >= transData.Stoploss[0].Value {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				if transData.Target >= mycandlesticks[i].Min && transData.Target <= mycandlesticks[i].Max {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Target, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				i++
			}
			if usdtEntry <= 0 {
				break
			}
		}
	}

	finalObject := &model.AlgorithmTransactionsFull{
		Transaction:     transactions,
		Finalusdt:       usdtEntry,
		Transamount:     int64(transNumbers.TransAmount),
		SuccessfulTrans: int64(transNumbers.SuccessTrans),
		LostTrans:       int64(transNumbers.LostTrans),
	}
	return finalObject
}

//AlgoMacdSP is algorithm MACD with SINGLE PAIR - with parameters for CHART - low, high, signal(MACD)
func AlgoMacdSPC(usdtEntry float64, feeMaker float64, feeTaker float64, candlesticks [][][]*futures.Kline, from string, to string, pair string, low int, high int, signal int) *model.AlgorithmTransactionsFull {
	var pairNum int
	switch pair {
	case "BTCUSDT":
		pairNum = 0
	case "ETHUSDT":
		pairNum = 1
	case "ADAUSDT":
		pairNum = 2
	case "NEOUSDT":
		pairNum = 3
	case "XLMUSDT":
		pairNum = 4
	case "XRPUSDT":
		pairNum = 5
	case "LINKUSDT":
		pairNum = 6
	}
	//Declaring important variables
	var it int
	transNumbers := &model.TransNumbers{
		TransAmount:  0,
		SuccessTrans: 0,
		LostTrans:    0,
		FeeSum:       0.00,
	}
	leverage := 10.00
	profit, stop := 1.01, 0.9975
	//Preparation for checking start index and stop index
	fromTimestamp, toTimestamp := DateToTimestampRange(from), DateToTimestampRange(to)
	indexstart, indexstop := utilities.GetStartStopCandles(candlesticks[pairNum][1], fromTimestamp, toTimestamp)

	//Container for Trasnactions
	var transactions []*model.AlgorithmTransactions
	var mycandlesticks []*model.MyKline
	//Getting SMA and EMA
	ema1 := utilities.EmovingAverage(candlesticks[pairNum][2], "close", "15m", 50, fromTimestamp, toTimestamp, pair)
	ema2 := utilities.EmovingAverage(candlesticks[pairNum][3], "close", "1h", 50, fromTimestamp, toTimestamp, pair)
	macd := utilities.MACD(candlesticks[pairNum][1], "close", int64(signal), int64(low), int64(high), "5m", fromTimestamp, toTimestamp, pair)
	var tableOfEMAs []*model.MovingAverage
	tableOfEMAs = append(tableOfEMAs, ema1, ema2)
	var tableOfMACDs []*model.MACD
	tableOfMACDs = append(tableOfMACDs, macd)
	mycandlesticks = utilities.CandlesToMyCandles(candlesticks[pairNum][1], indexstart, indexstop)
	mycandlesticks = utilities.MergeEMA(mycandlesticks, tableOfEMAs, indexstart, indexstop)
	mycandlesticks = utilities.MergeMACD(mycandlesticks, tableOfMACDs, indexstart, indexstop)
	mycandlesticks = mycandlesticks[0 : len(mycandlesticks)-2]
	it = indexstop - indexstart - 2
	//Going through candlesticks
	for i := 0; i < it-3; i++ {
		if mycandlesticks[i].MacD[0][1] >= 0 && mycandlesticks[i+1].MacD[0][1] < 0 && mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
			//Going over next candlesticks
			i = i + 2

			//Checking 1 PERIOD /////////////////////////////////////////////////////////////
			indexPeriod1 := i
			shouldStop := 0
			for mycandlesticks[i].MacD[0][2] < 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] > 0 || mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod1 < 3 || i-indexPeriod1 > 150 || shouldStop == 1 {
				continue
			}
			macMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod1:i], "min", "macd", 0)

			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for mycandlesticks[i].MacD[0][2] >= 0 && mycandlesticks[i].MacD[0][0] < 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] > 0 || mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod2 < 3 || i-indexPeriod2 > 40 && shouldStop == 1 {
				continue
			}

			//Checking 3 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod3 := i
			for mycandlesticks[i].MacD[0][2] < 0 && mycandlesticks[i].MacD[0][0] < 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] > 0 || mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod3 < 3 || i-indexPeriod3 > 150 {
				continue
			}
			candleMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "min", "body", 0)
			macMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "min", "macd", 0)
			candleMinPeriod3Tail := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "min", "tail", 0)
			var beg1Period int
			for iback := indexPeriod1; mycandlesticks[iback].MacD[0][2] < 0; iback-- {
				beg1Period = iback
			}
			candleMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[beg1Period:indexPeriod2-1], "min", "body", 0)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 <= candleMinPeriod3 || macMinPeriod1 >= macMinPeriod3 {
				continue
			} else {
				i++
			}
			//Starting transaction
			stop = candleMinPeriod3Tail / mycandlesticks[i].Open
			profit = 1 + 3*(1-stop)
			if stop < 1 {
				leverage = 0.05 / (1 - stop)
			} else {
				leverage = 0.05 / (-1 + stop)
			}
			leverage = math.Round(leverage)
			if stop >= 0.995 {
				//fmt.Println("Za niski target")
				continue
			}
			transData := utilities.StartingTransaction(mycandlesticks[i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Long")
			//Simulating begin in trade
			for i < it-2 {
				if mycandlesticks[i].Min <= transData.Stoploss[0].Value {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				if transData.Target >= mycandlesticks[i].Min && transData.Target <= mycandlesticks[i].Max {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Target, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				i++
			}
			if usdtEntry <= 0 {
				break
			}
		}
		if mycandlesticks[i].MacD[0][1] <= 0 && mycandlesticks[i+1].MacD[0][1] > 0 && mycandlesticks[i].Emas[0] < mycandlesticks[i].Emas[1] {
			//Going over next candlesticks
			i = i + 2

			//Checking 1 PERIOD /////////////////////////////////////////////////////////////
			indexPeriod1 := i
			shouldStop := 0
			for mycandlesticks[i].MacD[0][2] > 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] < 0 || mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod1 < 3 || i-indexPeriod1 > 150 || shouldStop == 1 {
				continue
			}
			macMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod1:i], "max", "macd", 0)

			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for mycandlesticks[i].MacD[0][2] <= 0 && mycandlesticks[i].MacD[0][0] > 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] < 0 || mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod2 < 3 || i-indexPeriod2 > 40 && shouldStop == 1 {
				continue
			}

			//Checking 3 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod3 := i
			for mycandlesticks[i].MacD[0][2] > 0 && mycandlesticks[i].MacD[0][0] > 0 && i < it-3 {
				if mycandlesticks[i].MacD[0][0] < 0 || mycandlesticks[i].Emas[0] > mycandlesticks[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if i-indexPeriod3 < 3 || i-indexPeriod3 > 150 {
				continue
			}
			candleMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "max", "body", 0)
			macMinPeriod3 := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "max", "macd", 0)
			candleMinPeriod3Tail := utilities.Period1ValuesBody(mycandlesticks[indexPeriod3:i], "max", "tail", 0)
			var beg1Period int
			for iback := indexPeriod1; mycandlesticks[iback].MacD[0][2] > 0; iback-- {
				beg1Period = iback
			}
			candleMinPeriod1 := utilities.Period1ValuesBody(mycandlesticks[beg1Period:indexPeriod2-1], "max", "body", 0)
			//pretty.Println(candleMinPeriod3, macMinPeriod3, candleMinPeriod3Tail, candleMinPeriod1, macMinPeriod1)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 >= candleMinPeriod3 || macMinPeriod1 <= macMinPeriod3 {
				continue
			} else {
				i++
			}
			//Starting transaction
			stop = candleMinPeriod3Tail / mycandlesticks[i].Open
			profit = 1 + 3*(1-stop)
			if stop < 1 {
				leverage = 0.05 / (1 - stop)
			} else {
				leverage = 0.05 / (-1 + stop)
			}
			leverage = math.Round(leverage)
			if stop <= 1.005 && stop <= 1.05 {
				//fmt.Println("Za niski target")
				continue
			}
			transData := utilities.StartingTransaction(mycandlesticks[i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Short")
			//Simulating begin in trade
			for i < it-2 {
				if mycandlesticks[i].Max >= transData.Stoploss[0].Value {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				if transData.Target >= mycandlesticks[i].Min && transData.Target <= mycandlesticks[i].Max {
					transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Target, transData)
					usdtEntry = transData.FinishBalance
					transactions = append(transactions, transData)
					break
				}
				i++
			}
			if usdtEntry <= 0 {
				break
			}
		}
	}

	finalObject := &model.AlgorithmTransactionsFull{
		Transaction:     transactions,
		Finalusdt:       usdtEntry,
		Transamount:     int64(transNumbers.TransAmount),
		SuccessfulTrans: int64(transNumbers.SuccessTrans),
		LostTrans:       int64(transNumbers.LostTrans),
	}
	return finalObject
}

//AlgoMacdAP is algorithm MACD with ALL PAIRS
func AlgoMacdAP(usdtEntry float64, feeMaker float64, feeTaker float64, candlesticks [][][]*futures.Kline, from string, to string, pairs []string) *model.AlgorithmTransactionsFull {
	var pairNum []int
	for _, pair := range pairs {
		switch pair {
		case "BTCUSDT":
			pairNum = append(pairNum, 0)
		case "ETHUSDT":
			pairNum = append(pairNum, 1)
		case "ADAUSDT":
			pairNum = append(pairNum, 2)
		case "NEOUSDT":
			pairNum = append(pairNum, 3)
		case "XLMUSDT":
			pairNum = append(pairNum, 4)
		case "XRPUSDT":
			pairNum = append(pairNum, 5)
		case "LINKUSDT":
			pairNum = append(pairNum, 6)
		}
	}
	//Declaring important variables
	var it int
	transNumbers := &model.TransNumbers{
		TransAmount:  0,
		SuccessTrans: 0,
		LostTrans:    0,
		FeeSum:       0.00,
	}
	leverage := 1.00
	profit, stop := 1.01, 0.9975
	//Preparation for checking start index and stop index
	fromTimestamp, toTimestamp := DateToTimestampRange(from), DateToTimestampRange(to)
	var indexstart []int
	var indexstop []int
	for _, num := range pairNum {
		istart, istop := utilities.GetStartStopCandles(candlesticks[num][1], fromTimestamp, toTimestamp)
		indexstart = append(indexstart, istart)
		indexstop = append(indexstop, istop)
	}
	//Container for Trasnactions
	var transactions []*model.AlgorithmTransactions
	var mycandlesticks []*model.MyKline
	var allcandlesticks [][]*model.MyKline
	//Getting SMA and EMA
	for i, num := range pairNum {
		ema1 := utilities.EmovingAverage(candlesticks[num][2], "close", "15m", 50, fromTimestamp, toTimestamp, pairs[num])
		ema2 := utilities.EmovingAverage(candlesticks[num][3], "close", "1h", 50, fromTimestamp, toTimestamp, pairs[num])
		macd1 := utilities.MACD(candlesticks[num][1], "close", 7, 7, 12, "5m", fromTimestamp, toTimestamp, pairs[num])
		macd2 := utilities.MACD(candlesticks[num][1], "close", 18, 13, 25, "5m", fromTimestamp, toTimestamp, pairs[num])
		var tableOfEMAs []*model.MovingAverage
		tableOfEMAs = append(tableOfEMAs, ema1, ema2)
		var tableOfMACDs []*model.MACD
		tableOfMACDs = append(tableOfMACDs, macd1, macd2)
		mycandlesticks = utilities.CandlesToMyCandles(candlesticks[num][1], indexstart[i], indexstop[i])
		mycandlesticks = utilities.MergeEMA(mycandlesticks, tableOfEMAs, indexstart[i], indexstop[i])
		mycandlesticks = utilities.MergeMACD(mycandlesticks, tableOfMACDs, indexstart[i], indexstop[i])
		allcandlesticks = append(allcandlesticks, mycandlesticks)
		pretty.Println()

	}
	it = len(allcandlesticks[0])
	//Going through candlesticks
	for i := 0; i < it-3; i++ {
		for coin := 0; coin < len(allcandlesticks); coin++ {
			coinIndexBegin := i
			if allcandlesticks[coin][i].MacD[0][1] >= 0 && allcandlesticks[coin][i+1].MacD[0][1] < 0 && allcandlesticks[coin][i].Emas[0] > allcandlesticks[coin][i].Emas[1] {
				//Going over next candlesticks
				i = i + 2
				//Checking 1 PERIOD /////////////////////////////////////////////////////////////
				indexPeriod1 := i
				shouldStop := 0
				for allcandlesticks[coin][i].MacD[0][2] < 0 && i < it-3 {
					if allcandlesticks[coin][i].MacD[0][0] > 0 || allcandlesticks[coin][i].Emas[0] < allcandlesticks[coin][i].Emas[1] {
						shouldStop = 1
						break
					}
					i++
				}
				if i-indexPeriod1 < 3 || i-indexPeriod1 > 150 || shouldStop == 1 {
					i = coinIndexBegin
					continue
				}
				macMinPeriod1 := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod1:i], "min", "macd", 0)

				//Checking 2 PERIOD //////////////////////////////////////////////////////////////
				indexPeriod2 := i
				shouldStop = 0
				for allcandlesticks[coin][i].MacD[0][2] >= 0 && allcandlesticks[coin][i].MacD[0][0] < 0 && i < it-3 {
					if allcandlesticks[coin][i].MacD[0][0] > 0 || allcandlesticks[coin][i].Emas[0] < allcandlesticks[coin][i].Emas[1] {
						shouldStop = 1
						break
					}
					i++
				}
				if i-indexPeriod2 < 3 || i-indexPeriod2 > 40 && shouldStop == 1 {
					i = coinIndexBegin
					continue
				}

				//Checking 3 PERIOD //////////////////////////////////////////////////////////////
				indexPeriod3 := i
				for allcandlesticks[coin][i].MacD[0][2] < 0 && allcandlesticks[coin][i].MacD[0][0] < 0 && i < it-3 {
					if allcandlesticks[coin][i].MacD[0][0] > 0 || allcandlesticks[coin][i].Emas[0] < allcandlesticks[coin][i].Emas[1] {
						shouldStop = 1
						break
					}
					i++
				}
				if i-indexPeriod3 < 3 || i-indexPeriod3 > 150 {
					i = coinIndexBegin
					continue
				}
				candleMinPeriod3 := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod3:i], "min", "body", 0)
				macMinPeriod3 := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod3:i], "min", "macd", 0)
				candleMinPeriod3Tail := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod3:i], "min", "tail", 0)
				var beg1Period int
				for iback := indexPeriod1; allcandlesticks[coin][iback].MacD[0][2] < 0; iback-- {
					beg1Period = iback
				}
				candleMinPeriod1 := utilities.Period1ValuesBody(allcandlesticks[coin][beg1Period:indexPeriod2-1], "min", "body", 0)
				//Checking final conditions ///////////////////////////////////////////////////////////
				if candleMinPeriod1 <= candleMinPeriod3 || macMinPeriod1 >= macMinPeriod3 {
					i = coinIndexBegin
					continue
				} else {
					i++
				}
				//Starting transaction
				stop = candleMinPeriod3Tail / allcandlesticks[coin][i].Open
				profit = 1 + 3*(1-stop)
				if stop < 1 {
					leverage = 0.05 / (1 - stop)
				} else {
					leverage = 0.05 / (-1 + stop)
				}
				leverage = math.Round(leverage)
				if stop >= 0.995 {
					fmt.Println("Za niski target")
					i = coinIndexBegin
					continue
				}
				transData := utilities.StartingTransaction(allcandlesticks[coin][i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Long")
				//Simulating begin in trade
				for i < it-2 {
					if allcandlesticks[coin][i].Min <= transData.Stoploss[0].Value {
						transData = utilities.ClosingTransaction(allcandlesticks[coin][i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
						usdtEntry = transData.FinishBalance
						transactions = append(transactions, transData)
						break
					}
					if transData.Target >= allcandlesticks[coin][i].Min && transData.Target <= allcandlesticks[coin][i].Max {
						transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Target, transData)
						usdtEntry = transData.FinishBalance
						transactions = append(transactions, transData)
						break
					}
					i++
				}
				if usdtEntry <= 0 {
					break
				}
			}
			if allcandlesticks[coin][i].MacD[1][1] <= 0 && allcandlesticks[coin][i+1].MacD[1][1] > 0 && allcandlesticks[coin][i].Emas[0] < allcandlesticks[coin][i].Emas[1] {
				//Going over next candlesticks
				i = i + 2

				//Checking 1 PERIOD /////////////////////////////////////////////////////////////
				indexPeriod1 := i
				shouldStop := 0
				for allcandlesticks[coin][i].MacD[1][2] > 0 && i < it-3 {
					if allcandlesticks[coin][i].MacD[1][0] < 0 || allcandlesticks[coin][i].Emas[0] > allcandlesticks[coin][i].Emas[1] {
						shouldStop = 1
						break
					}
					i++
				}
				if i-indexPeriod1 < 3 || i-indexPeriod1 > 150 || shouldStop == 1 {
					i = coinIndexBegin
					continue
				}
				macMinPeriod1 := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod1:i], "max", "macd", 1)

				//Checking 2 PERIOD //////////////////////////////////////////////////////////////
				indexPeriod2 := i
				shouldStop = 0
				for allcandlesticks[coin][i].MacD[1][2] <= 0 && allcandlesticks[coin][i].MacD[1][0] > 0 && i < it-3 {
					if allcandlesticks[coin][i].MacD[1][0] < 0 || allcandlesticks[coin][i].Emas[0] > allcandlesticks[coin][i].Emas[1] {
						shouldStop = 1
						break
					}
					i++
				}
				if i-indexPeriod2 < 3 || i-indexPeriod2 > 40 && shouldStop == 1 {
					i = coinIndexBegin
					continue
				}

				//Checking 3 PERIOD //////////////////////////////////////////////////////////////
				indexPeriod3 := i
				for allcandlesticks[coin][i].MacD[1][2] > 0 && allcandlesticks[coin][i].MacD[1][0] > 0 && i < it-3 {
					if allcandlesticks[coin][i].MacD[1][0] < 0 || allcandlesticks[coin][i].Emas[0] > allcandlesticks[coin][i].Emas[1] {
						shouldStop = 1
						break
					}
					i++
				}
				if i-indexPeriod3 < 3 || i-indexPeriod3 > 150 {
					i = coinIndexBegin
					continue
				}
				candleMinPeriod3 := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod3:i], "max", "body", 1)
				macMinPeriod3 := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod3:i], "max", "macd", 1)
				candleMinPeriod3Tail := utilities.Period1ValuesBody(allcandlesticks[coin][indexPeriod3:i], "max", "tail", 1)
				var beg1Period int
				for iback := indexPeriod1; allcandlesticks[coin][iback].MacD[1][2] > 0; iback-- {
					beg1Period = iback
				}
				candleMinPeriod1 := utilities.Period1ValuesBody(allcandlesticks[coin][beg1Period:indexPeriod2-1], "max", "body", 1)
				//Checking final conditions ///////////////////////////////////////////////////////////
				if candleMinPeriod1 >= candleMinPeriod3 || macMinPeriod1 <= macMinPeriod3 {
					i = coinIndexBegin
					continue
				} else {
					i++
				}
				//Starting transaction
				stop = candleMinPeriod3Tail / allcandlesticks[coin][i].Open
				profit = 1 + 3*(1-stop)
				if stop < 1 {
					leverage = 0.05 / (1 - stop)
				} else {
					leverage = 0.05 / (-1 + stop)
				}
				leverage = math.Round(leverage)
				if stop <= 1.005 {
					//fmt.Println("Za niski target")
					continue
				}
				transData := utilities.StartingTransaction(allcandlesticks[coin][i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Short")
				//Simulating begin in trade
				for i < it-2 {
					if allcandlesticks[coin][i].Max >= transData.Stoploss[0].Value {
						transData = utilities.ClosingTransaction(allcandlesticks[coin][i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
						usdtEntry = transData.FinishBalance
						transactions = append(transactions, transData)
						break
					}
					if transData.Target >= allcandlesticks[coin][i].Min && transData.Target <= allcandlesticks[coin][i].Max {
						transData = utilities.ClosingTransaction(allcandlesticks[coin][i], transNumbers, feeMaker, transData.Target, transData)
						usdtEntry = transData.FinishBalance
						transactions = append(transactions, transData)
						break
					}
					i++
				}
				if usdtEntry <= 0 {
					break
				}
			}
		}
	}

	finalObject := &model.AlgorithmTransactionsFull{
		Transaction:     transactions,
		Finalusdt:       usdtEntry,
		Transamount:     int64(transNumbers.TransAmount),
		SuccessfulTrans: int64(transNumbers.SuccessTrans),
		LostTrans:       int64(transNumbers.LostTrans),
	}
	return finalObject
}

//AlgoMacdAP is algorithm RSI with ALL PAIRS - with parameters for CHART - target, stoploss, value(RSI)
func AlgoRsiAPC(usdtEntry float64, feeMaker float64, feeTaker float64, candlesticks [][][]*futures.Kline, pairs []string, from string, to string, target float64, stoploss float64, value int) *model.AlgorithmTransactionsFull {
	var pairNum []int
	for _, pair := range pairs {
		switch pair {
		case "BTCUSDT":
			pairNum = append(pairNum, 0)
		case "ETHUSDT":
			pairNum = append(pairNum, 1)
		case "ADAUSDT":
			pairNum = append(pairNum, 2)
		case "NEOUSDT":
			pairNum = append(pairNum, 3)
		case "XLMUSDT":
			pairNum = append(pairNum, 4)
		case "XRPUSDT":
			pairNum = append(pairNum, 5)
		case "LINKUSDT":
			pairNum = append(pairNum, 6)
		}
	}
	//Declaring important variables
	var it int
	transNumbers := &model.TransNumbers{
		TransAmount:  0,
		SuccessTrans: 0,
		LostTrans:    0,
		FeeSum:       0.00,
	}
	leverage := 1.00
	profit, stop := target, stoploss
	//Preparation for checking start index and stop index
	fromTimestamp, toTimestamp := DateToTimestampRange(from), DateToTimestampRange(to)
	var indexstart []int
	var indexstop []int
	for _, num := range pairNum {
		istart, istop := utilities.GetStartStopCandles(candlesticks[num][1], fromTimestamp, toTimestamp)
		indexstart = append(indexstart, istart)
		indexstop = append(indexstop, istop)
	}
	//Container for Trasnactions
	var transactions []*model.AlgorithmTransactions
	var mycandlesticks []*model.MyKline
	var allcandlesticks [][]*model.MyKline
	//Getting SMA and EMA
	for i, num := range pairNum {
		ema1 := utilities.EmovingAverage(candlesticks[num][2], "close", "15m", 50, fromTimestamp, toTimestamp, pairs[num])
		ema2 := utilities.EmovingAverage(candlesticks[num][3], "close", "1h", 50, fromTimestamp, toTimestamp, pairs[num])
		macd1 := utilities.MACD(candlesticks[num][1], "close", 7, 7, 12, "5m", fromTimestamp, toTimestamp, pairs[num])
		macd2 := utilities.MACD(candlesticks[num][1], "close", 18, 13, 25, "5m", fromTimestamp, toTimestamp, pairs[num])
		rsi := utilities.RSI(candlesticks[num][1], "close", "5m", 14, fromTimestamp, toTimestamp, pairs[num])
		var tableOfEMAs []*model.MovingAverage
		tableOfEMAs = append(tableOfEMAs, ema1, ema2)
		var tableOfMACDs []*model.MACD
		tableOfMACDs = append(tableOfMACDs, macd1, macd2)
		var tableOfRSIs []*model.RSI
		tableOfRSIs = append(tableOfRSIs, rsi)
		mycandlesticks = utilities.CandlesToMyCandles(candlesticks[num][1], indexstart[i], indexstop[i])
		mycandlesticks = utilities.MergeEMA(mycandlesticks, tableOfEMAs, indexstart[i], indexstop[i])
		mycandlesticks = utilities.MergeMACD(mycandlesticks, tableOfMACDs, indexstart[i], indexstop[i])
		mycandlesticks = utilities.MergeRSI(mycandlesticks, tableOfRSIs, indexstart[i], indexstop[i])
		allcandlesticks = append(allcandlesticks, mycandlesticks)

	}
	it = len(allcandlesticks[0])
	//Going through candlesticks
	for i := 0; i < it-4; i++ {
		for coin := 0; coin < len(allcandlesticks); coin++ {
			// if allcandlesticks[coin][i].RSI >= float64(value) {
			// 	//Going over next candlesticks
			// 	i = i + 1

			// 	transData := StartingTransaction(allcandlesticks[coin][i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Short")
			// 	//Simulating begin in trade
			// 	for i < it-4 {
			// 		if allcandlesticks[coin][i].Max >= transData.Stoploss[0].Value {
			// 			transData = ClosingTransaction(allcandlesticks[coin][i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
			// 			usdtEntry = transData.FinishBalance
			// 			transactions = append(transactions, transData)
			// 			break
			// 		}
			// 		if transData.Target >= allcandlesticks[coin][i].Min && transData.Target <= allcandlesticks[coin][i].Max {
			// 			transData = ClosingTransaction(allcandlesticks[coin][i], transNumbers, feeMaker, transData.Target, transData)
			// 			usdtEntry = transData.FinishBalance
			// 			transactions = append(transactions, transData)
			// 			break
			// 		}
			// 		i++
			// 	}
			// 	if usdtEntry <= 0 {
			// 		break
			// 	}
			// }
			if allcandlesticks[coin][i].RSI[0] <= float64(value) {
				//Going over next candlesticks
				i = i + 1

				transData := utilities.StartingTransaction(allcandlesticks[coin][i], transNumbers, usdtEntry, leverage, stop, profit, feeTaker, "Long")
				//Simulating begin in trade
				for i < it-4 {
					if allcandlesticks[coin][i].Min <= transData.Stoploss[0].Value {
						transData = utilities.ClosingTransaction(allcandlesticks[coin][i], transNumbers, feeMaker, transData.Stoploss[0].Value, transData)
						usdtEntry = transData.FinishBalance
						transactions = append(transactions, transData)
						break
					}
					if transData.Target >= allcandlesticks[coin][i].Min && transData.Target <= allcandlesticks[coin][i].Max {
						transData = utilities.ClosingTransaction(mycandlesticks[i], transNumbers, feeMaker, transData.Target, transData)
						usdtEntry = transData.FinishBalance
						transactions = append(transactions, transData)
						break
					}
					i++
				}
				if usdtEntry <= 0 {
					break
				}
			}
		}
	}
	finalObject := &model.AlgorithmTransactionsFull{
		Transaction:     transactions,
		Finalusdt:       usdtEntry,
		Transamount:     int64(transNumbers.TransAmount),
		SuccessfulTrans: int64(transNumbers.SuccessTrans),
		LostTrans:       int64(transNumbers.LostTrans),
	}
	return finalObject
}
