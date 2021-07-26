package realbot

import (
	"math"

	rmf "github.com/CraZzier/bot/realbot/utilities"
	"github.com/kr/pretty"
)

func (bot *RealBot) TestMACD() {
	//Declaring important variables
	var it int
	notifier := NewNotifier()
	notifier.AddMsg(bot.DisplayTime(), " Test formacji MACD")

	leverage := 1.00
	profit, stop := 1.01, 0.9975
	//Going through candlesticks
	for coin := 0; coin < len(bot.Pairs); coin++ {
		//Defining limits not to enter empty values of indicators
		it = len(bot.CustomKline[coin][0])
		candles := bot.CustomKline[coin][0]
		i := it - 1
		limit := 600
		//Going manually though the candles
		coinIndexBegin := i
		//Long
		//fmt.Println(candles[i].ATR[0])
		if candles[i].MacD[0][2] >= 0 &&
			candles[i-1].MacD[0][2] < 0 &&
			candles[i-1].MacD[0][0] < 0 &&
			candles[i].MacD[0][0] < 0 &&
			candles[i].Emas[0] > candles[i].Emas[1] {
			//Going over next candlesticks
			i = i - 1
			notifier.AddMsg(bot.DisplayTime(), "Potencjalny long: ", bot.Pairs[coin])
			//Checking 3 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod3 := i
			shouldStop := 0
			for candles[i].MacD[0][2] < 0 && candles[i].MacD[0][0] < 0 && i > limit {
				if candles[i].Emas[0] < candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i--
			}
			if indexPeriod3-i < 3 || indexPeriod3-i > 150 || shouldStop == 1 {
				continue
			}
			notifier.AddMsg("Trzeci okres: ", indexPeriod3-i, " czerwonych candlesticków")
			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for candles[i].MacD[0][2] >= 0 && candles[i].MacD[0][0] < 0 && i > limit {
				if candles[i].Emas[0] < candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i--
			}
			if indexPeriod2-i < 3 || indexPeriod2-i > 40 || shouldStop == 1 {
				continue
			}
			notifier.AddMsg("Drugi okres: ", indexPeriod2-i, " zielonych candlesticków")
			//Checking 1 PERIOD /////////////////////////////////////////////////////////////
			indexPeriod1 := i
			shouldStop = 0
			for candles[i].MacD[0][2] < 0 && i > limit {
				if candles[i].Emas[0] < candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i--
			}
			if indexPeriod1-i < 3 || indexPeriod1-i > 150 || shouldStop == 1 {
				continue
			}
			notifier.AddMsg("Pierwszy okres: ", indexPeriod1-i, " czerwonych candlesticków")
			found, position := rmf.CheckCrossing(candles[i:indexPeriod1], "Down", 0)
			if !found {
				notifier.AddMsg("Pierwszy okres: nie znaleziono przecięcia")
				continue
			} else {
				notifier.AddMsg("Pierwszy okres: znaleziono przecięcie: ", position)
			}
			candleMinPeriod3 := rmf.Period1ValuesBody(candles[indexPeriod2:indexPeriod3], "min", "body", 0)
			macMinPeriod3 := rmf.Period1ValuesBody(candles[indexPeriod2:indexPeriod3], "min", "macd", 0)
			candleMinPeriod3Tail := rmf.Period1ValuesBody(candles[indexPeriod2:indexPeriod3], "min", "tail", 0)
			candleMinPeriod1 := rmf.Period1ValuesBody(candles[i:indexPeriod1], "min", "body", 0)
			macMinPeriod1 := rmf.Period1ValuesBody(candles[i:indexPeriod1], "min", "macd", 0)
			pretty.Println(candleMinPeriod3, macMinPeriod3, candleMinPeriod3Tail, candleMinPeriod1, macMinPeriod1)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 <= candleMinPeriod3 || macMinPeriod1 >= macMinPeriod3 {
				notifier.AddMsg("Rozbieżności MACD i wykresu nie zgadzają się")
				continue
			}
			notifier.AddMsg("Rozbieżności MACD i wykresu zgadzają się")
			//Starting transaction
			stop = candleMinPeriod3Tail / candles[coinIndexBegin].Close
			profit = 1 + 3*(1-stop)
			if stop < 1 {
				leverage = 0.05 / (1 - stop)
			} else {
				leverage = 0.05 / (-1 + stop)
			}
			leverage = math.Round(leverage)
			if stop >= 0.995 {
				notifier.AddMsg("Za niski target")
				continue
			}
			notifier.AddMsg("Stoploss jest w porzadku")
			actualPrice := bot.GetTickerPrice(bot.Pairs[coin])
			profit = actualPrice * profit
			stop = actualPrice * stop
			bot.MakeCompleteOrder(int(leverage), bot.Pairs[coin], profit, stop, "BUY")

		}
		//Short
		if candles[i].MacD[1][2] <= 0 &&
			candles[i-1].MacD[1][2] > 0 &&
			candles[i-1].MacD[1][0] > 0 &&
			candles[i].MacD[1][0] > 0 &&
			candles[i].Emas[0] < candles[i].Emas[1] {
			//Going over next candlesticks
			i = i - 1
			notifier.AddMsg(bot.DisplayTime(), "Potencjalny short: ", bot.Pairs[coin])
			//Checking 3 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod3 := i
			shouldStop := 0
			for candles[i].MacD[1][2] > 0 && candles[i].MacD[1][0] > 0 && i > limit {
				if candles[i].Emas[0] > candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i--
			}
			if indexPeriod3-i < 3 || indexPeriod3-i > 150 || shouldStop == 1 {
				continue
			}
			notifier.AddMsg("Trzeci okres: ", indexPeriod3-i, " zielonych candlesticków")
			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for candles[i].MacD[1][2] <= 0 && candles[i].MacD[1][0] > 0 && i > limit {
				if candles[i].Emas[0] > candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i--
			}
			if indexPeriod2-i < 3 || indexPeriod2-i > 40 || shouldStop == 1 {
				continue
			}
			notifier.AddMsg("Drugi okres: ", indexPeriod2-i, " czerwonych candlesticków")
			//Checking 1 PERIOD /////////////////////////////////////////////////////////////
			indexPeriod1 := i
			shouldStop = 0
			for candles[i].MacD[1][2] > 0 && i > limit {
				if candles[i].Emas[0] > candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i--
			}
			if indexPeriod1-i < 3 || indexPeriod1-i > 150 || shouldStop == 1 {
				continue
			}
			notifier.AddMsg("Pierwszy okres: ", indexPeriod1-i, " zielonych candlesticków")
			found, position := rmf.CheckCrossing(candles[i:indexPeriod1], "Up", 0)
			if !found {
				notifier.AddMsg("Pierwszy okres: nie znaleziono przecięcia")
				continue
			} else {
				notifier.AddMsg("Pierwszy okres: znaleziono przecięcie: ", position)
			}
			candleMinPeriod3 := rmf.Period1ValuesBody(candles[indexPeriod2:indexPeriod3], "max", "body", 0)
			macMinPeriod3 := rmf.Period1ValuesBody(candles[indexPeriod2:indexPeriod3], "max", "macd", 0)
			candleMinPeriod3Tail := rmf.Period1ValuesBody(candles[indexPeriod2:indexPeriod3], "max", "tail", 0)
			candleMinPeriod1 := rmf.Period1ValuesBody(candles[i:indexPeriod1], "max", "body", 0)
			macMinPeriod1 := rmf.Period1ValuesBody(candles[i:indexPeriod1], "max", "macd", 0)
			pretty.Println(candleMinPeriod3, macMinPeriod3, candleMinPeriod3Tail, candleMinPeriod1, macMinPeriod1)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 >= candleMinPeriod3 || macMinPeriod1 <= macMinPeriod3 {

				notifier.AddMsg("Rozbieżności MACD i wykresu nie zgadzają się")
				continue
			}
			notifier.AddMsg("Rozbieżności MACD i wykresu zgadzają się")
			//Starting transaction
			stop = candleMinPeriod3Tail / candles[coinIndexBegin].Close
			profit = 1 + 3*(1-stop)
			if stop < 1 {
				leverage = 0.05 / (1 - stop)
			} else {
				leverage = 0.05 / (-1 + stop)
			}
			leverage = math.Round(leverage)
			if stop >= 0.995 {
				notifier.AddMsg("Za niski target")
				continue
			}
			notifier.AddMsg("Stoploss w porządku")
			actualPrice := bot.GetTickerPrice(bot.Pairs[coin])
			profit = actualPrice * profit
			stop = actualPrice * stop
			bot.MakeCompleteOrder(int(leverage), bot.Pairs[coin], profit, stop, "SELL")

		}
	}
	notifier.PrintAllOut()
	if notifier.MsgNumber > 1 {
		notifier.SendEmail()
	}
}
