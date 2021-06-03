package realbot

import (
	"fmt"
	"math"

	rmf "github.com/CraZzier/bot/realbot/utilities"
	"github.com/kr/pretty"
)

func (bot *RealBot) TestFormation() {
	//Declaring important variables
	var it int
	fmt.Println("Test formacji")
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
		if candles[i].MacD[0][2] >= 0 &&
			candles[i-1].MacD[0][2] < 0 &&
			candles[i-1].MacD[0][0] < 0 &&
			candles[i].MacD[0][0] < 0 &&
			candles[i].Emas[0] > candles[i].Emas[1] {
			//Going over next candlesticks
			fmt.Println(bot.Pairs[coin])
			i = i - 1
			fmt.Println("1.Pierwszy zielony po czerwonym znaleziony")
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
			fmt.Println("2. Trzeci period - ok")
			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for candles[i].MacD[0][2] >= 0 && candles[i].MacD[0][0] < 0 && i > limit {
				if candles[i].Emas[0] < candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if indexPeriod2-i < 3 || indexPeriod2-i > 40 || shouldStop == 1 {
				continue
			}
			fmt.Println("3. Drugi period - ok")
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
			found, _ := rmf.CheckCrossing(candles[indexPeriod1:indexPeriod2-1], "Down", 0)
			fmt.Println("4. Pierwszy period przecięcie - ok")
			if indexPeriod1-i < 3 || indexPeriod1-i > 150 || !found || shouldStop == 1 {
				continue
			}
			fmt.Println("5. Pierwszy period  - ok")
			candleMinPeriod3 := rmf.Period1ValuesBody(candles[i:indexPeriod3], "min", "body", 0)
			macMinPeriod3 := rmf.Period1ValuesBody(candles[i:indexPeriod3], "min", "macd", 0)
			candleMinPeriod3Tail := rmf.Period1ValuesBody(candles[i:indexPeriod3], "min", "tail", 0)
			candleMinPeriod1 := rmf.Period1ValuesBody(candles[indexPeriod1:indexPeriod2-1], "min", "body", 0)
			macMinPeriod1 := rmf.Period1ValuesBody(candles[indexPeriod1:indexPeriod2-1], "min", "macd", 0)
			pretty.Println(candleMinPeriod3, macMinPeriod3, candleMinPeriod3Tail, candleMinPeriod1, macMinPeriod1)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 <= candleMinPeriod3 || macMinPeriod1 >= macMinPeriod3 {
				continue
			}
			fmt.Println("6. Candle i macd  - ok")
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
				fmt.Println("Za niski target")
				continue
			}
			fmt.Println("7. Stoploss  - ok")
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
			fmt.Println(bot.Pairs[coin])
			i = i - 1
			fmt.Println("1.Pierwszy czerwony po zielonym znaleziony")
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
			fmt.Println("2. Trzeci period - ok")
			//Checking 2 PERIOD //////////////////////////////////////////////////////////////
			indexPeriod2 := i
			shouldStop = 0
			for candles[i].MacD[1][2] <= 0 && candles[i].MacD[1][0] > 0 && i > limit {
				if candles[i].Emas[0] > candles[i].Emas[1] {
					shouldStop = 1
					break
				}
				i++
			}
			if indexPeriod2-i < 3 || indexPeriod2-i > 40 || shouldStop == 1 {
				continue
			}
			fmt.Println("3. Drugi period - ok")
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
			found, _ := rmf.CheckCrossing(candles[indexPeriod1:indexPeriod2-1], "Up", 1)
			fmt.Println("4. Pierwszy period przecięcie - ok")
			if indexPeriod1-i < 3 || indexPeriod1-i > 150 || !found || shouldStop == 1 {
				continue
			}
			fmt.Println("5. Pierwszy period  - ok")
			candleMinPeriod3 := rmf.Period1ValuesBody(candles[i:indexPeriod3], "max", "body", 1)
			macMinPeriod3 := rmf.Period1ValuesBody(candles[i:indexPeriod3], "max", "macd", 1)
			candleMinPeriod3Tail := rmf.Period1ValuesBody(candles[i:indexPeriod3], "max", "tail", 1)
			candleMinPeriod1 := rmf.Period1ValuesBody(candles[indexPeriod1:indexPeriod2-1], "max", "body", 1)
			macMinPeriod1 := rmf.Period1ValuesBody(candles[indexPeriod1:indexPeriod2-1], "max", "macd", 1)
			pretty.Println(candleMinPeriod3, macMinPeriod3, candleMinPeriod3Tail, candleMinPeriod1, macMinPeriod1)
			//Checking final conditions ///////////////////////////////////////////////////////////
			if candleMinPeriod1 >= candleMinPeriod3 || macMinPeriod1 <= macMinPeriod3 {
				continue
			}
			fmt.Println("6. Candle i macd  - ok")
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
				fmt.Println("Za niski target")
				continue
			}
			fmt.Println("7. Stoploss  - ok")
			actualPrice := bot.GetTickerPrice(bot.Pairs[coin])
			profit = actualPrice * profit
			stop = actualPrice * stop
			bot.MakeCompleteOrder(int(leverage), bot.Pairs[coin], profit, stop, "SELL")

		}
	}
}
