package realbot

import (
	"context"
	"fmt"

	"github.com/CraZzier/bot/model"
	"github.com/CraZzier/bot/realbot/stop"
	rmf "github.com/CraZzier/bot/realbot/utilities"
	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
)

var (
	apiKey        = "XOIxbjxLWErp4pTGqOs0svp1PWgPuUgyz7QVp735QVinF3T5i43VrpQU50jVIqQW"
	secretKey     = "tRJ75yZh8gq9Hb1gwTfESz083wlwgd2qt71BprinrODm7BijYDmqiI90YftksIAz"
	futuresClient = binance.NewFuturesClient(apiKey, secretKey) // USDT-M Futures
)

//RealBot is main struct that runs bot for single pair
type RealBot struct {
	Pairs           []string
	Intervals       []string
	KlinesData      [][][]*futures.Kline
	CustomKline     [][][]*model.MyKline
	Account         *futures.Account
	Balance         float64
	OpenOrders      [][]*futures.CreateOrderResponse
	CandleLimit     int
	ActivePositions int
	FeeTaker        float64
	FeeMaker        float64
	ListenKey       string
	ExchangeInfo    *futures.ExchangeInfo
	PrecisionTable  [][]int
	Stop            *stop.Stop
}

//Initialization is ment to be getting data of user while strating the bot
func (bot *RealBot) Initialization(pairs []string, intervals []string, candleLimit int) {
	fmt.Println("Realbot initialised")
	bot.Pairs = pairs
	bot.Intervals = intervals
	bot.CandleLimit = candleLimit
	bot.FeeMaker = 0.0040
	bot.FeeTaker = 0.0040

	bot.Stop = &stop.Stop{}
	bot.Stop.Init("macd", "bollinger")

	bot.GetAccountInfo()
	bot.GetBalanceInfo()
	bot.NumberOfActivePositions()
	bot.GetExchangeInfo()
	bot.GetPrecisionTable()
	bot.ListenKey, _ = futuresClient.NewStartUserStreamService().Do(context.Background())
	//Making space for klines
	bot.KlinesData = make([][][]*futures.Kline, len(pairs))
	bot.CustomKline = make([][][]*model.MyKline, len(pairs))
	for i := range pairs {
		bot.KlinesData[i] = make([][]*futures.Kline, len(intervals))
		bot.CustomKline[i] = make([][]*model.MyKline, len(intervals))
	}
	//First Klines Downloading
	for i, x := range intervals {
		for o, v := range pairs {
			bot.CollectKlines(v, x, o, i, bot.CandleLimit)
		}
	}
	bot.UserStreamWS()
	for i := range intervals {
		for o := range pairs {
			bot.CustomKline[o][i] = rmf.ToMyKline(bot.KlinesData[o][i], 0, 999)
		}
	}
	var ema1T, ema2T []*model.MovingAverage
	var macd1T, macd2T []*model.MACD
	var bb1T []*model.BollingerBands
	var atr1T []*model.ATR
	//Adding indicators to first candlesticks
	for num := range pairs {
		ema1 := rmf.EMA(bot.CustomKline[num][0], "close", "5m", 150, 0, 999, pairs[num])
		ema2 := rmf.EMA(bot.CustomKline[num][0], "close", "5m", 600, 0, 999, pairs[num])
		macd1 := rmf.MACD(bot.CustomKline[num][0], "close", 7, 7, 12, "5m", 0, 999, pairs[num])
		macd2 := rmf.MACD(bot.CustomKline[num][0], "close", 18, 13, 25, "5m", 0, 999, pairs[num])
		bb := rmf.BollingerBands(bot.CustomKline[num][0], 20, "5m", 5.5, 0, 999, pairs[num])
		atr := rmf.ATR(bot.CustomKline[num][0], 14, "5m", 0, 999, pairs[num])
		var tableOfEMAs []*model.MovingAverage
		var tableOfMACDs []*model.MACD
		var tableOfBollingerBands []*model.BollingerBands
		var tableOfATRs []*model.ATR
		tableOfEMAs = append(tableOfEMAs, ema1, ema2)
		tableOfMACDs = append(tableOfMACDs, macd1, macd2)
		tableOfBollingerBands = append(tableOfBollingerBands, bb)
		tableOfATRs = append(tableOfATRs, atr)
		bot.CustomKline[num][0] = rmf.MergeEMA(bot.CustomKline[num][0], tableOfEMAs, 0, 999)
		bot.CustomKline[num][0] = rmf.MergeMACD(bot.CustomKline[num][0], tableOfMACDs, 0, 999)
		bot.CustomKline[num][0] = rmf.MergeBollingerBands(bot.CustomKline[num][0], tableOfBollingerBands, 0, 999)
		bot.CustomKline[num][0] = rmf.MergeATR(bot.CustomKline[num][0], tableOfATRs, 0, 999)
		ema1T = append(ema1T, ema1)
		ema2T = append(ema2T, ema2)
		macd1T = append(macd1T, macd1)
		macd2T = append(macd2T, macd2)
		bb1T = append(bb1T, bb)
		atr1T = append(atr1T, atr)
	}
	bot.NumberOfActivePositions()
	//Creating channel for a communication
	klineChannels5m := make(chan int)
	klineChannels15m := make(chan int)
	klineChannels1h := make(chan int)
	for i, x := range intervals {
		for o, v := range pairs {
			go bot.CollectKlinesWS(v, x, o, i, klineChannels5m, klineChannels15m, klineChannels1h)
		}
	}
	sum5m, sum15m, sum1h := 0, 0, 0
	for {
		select {
		case msg1 := <-klineChannels5m:
			sum5m += msg1

		case msg2 := <-klineChannels15m:
			sum15m += msg2

		case msg3 := <-klineChannels1h:
			sum1h += msg3
		}
		if sum5m == len(pairs) {
			sum5m = 0
			//Updating
			for o := range pairs {
				rmf.UpdateMACD(bot.CustomKline[o][0], macd1T[o])
				rmf.UpdateMACD(bot.CustomKline[o][0], macd2T[o])
				rmf.UpdateEMA(bot.CustomKline[o][0], ema1T[o])
				rmf.UpdateEMA(bot.CustomKline[o][0], ema2T[o])
				rmf.UpdateBollingerBands(bot.CustomKline[o][0], bb1T[o])
				rmf.UpdateATR(bot.CustomKline[o][0], atr1T[o])
			}
			//Checking account status
			bot.GetAccountInfo()
			bot.GetBalanceInfo()
			bot.NumberOfActivePositions()
			//Testing formation
			if bot.ActivePositions == 0 {
				bot.TestMACD()
			}
		}
		if sum15m == len(pairs) {
			sum15m = 0
			//Updating userstream
			futuresClient.NewKeepaliveUserStreamService().ListenKey(bot.ListenKey).Do(context.Background())
		}
		if sum1h == len(pairs) {
			sum1h = 0
		}
	}

}
