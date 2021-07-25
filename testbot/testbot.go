package testbot

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

//TestBot is main struct that runs bot for single pair
type TestBot struct {
	KlinesDataCsv       [][][]*futures.Kline
	WalletFuturesUSDT   float64
	FeeTaker            float64
	FeeMaker            float64
	TransactionFullData *model.AlgorithmTransactionsFull
}

//Initialization is ment to be getting data of user while strating the bot
func (bot *TestBot) Initialization() {
	bot.CollectKlinesCsv()
	bot.WalletFuturesUSDT = 1000
	bot.FeeMaker = 0.00036
	bot.FeeTaker = 0.00036
}

//CollectKlinesCsv lets to get historical pair data
func (bot *TestBot) CollectKlinesCsv() {
	pair := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "NEOUSDT", "XLMUSDT", "XRPUSDT", "LINKUSDT"}
	intervalsToDownload := []string{"1m", "5m", "15m", "1h", "1d", "1w"}
	bot.KlinesDataCsv = make([][][]*futures.Kline, len(pair))
	for x := 0; x < len(pair); x++ {
		for k := 0; k < len(intervalsToDownload); k++ {
			csvFile, err := os.Open("testbot/candledata/" + pair[x] + intervalsToDownload[k] + ".csv")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Successfully Opened CSV file:" + pair[x] + intervalsToDownload[k] + ".csv")
			defer csvFile.Close()

			csvLines, err := csv.NewReader(csvFile).ReadAll()
			if err != nil {
				fmt.Println(err)
			}
			var tableOfKlines []*futures.Kline
			for i, line := range csvLines {
				if i == 0 {
					continue
				}
				opentime := DateToTimestamp(line[0])
				closetime := opentime + 59999
				open := line[1]
				high := line[2]
				low := line[3]
				close := line[4]
				volume := line[5]
				quoteAssetVolume := line[6]
				numoftrad, _ := strconv.ParseInt(line[7], 10, 64)
				takerbuybase := line[8]
				takerbuyquote := line[9]

				tempKline := futures.Kline{
					OpenTime:                 opentime,
					Open:                     open,
					High:                     high,
					Low:                      low,
					Close:                    close,
					Volume:                   volume,
					CloseTime:                closetime,
					QuoteAssetVolume:         quoteAssetVolume,
					TradeNum:                 numoftrad,
					TakerBuyBaseAssetVolume:  takerbuybase,
					TakerBuyQuoteAssetVolume: takerbuyquote,
				}
				tableOfKlines = append(tableOfKlines, &tempKline)

			}
			bot.KlinesDataCsv[x] = append(bot.KlinesDataCsv[x], tableOfKlines)
		}
	}
}

//TestAlgorithm servers for testign algorithm conception
func (bot *TestBot) TestAlgorithm(from string, to string, pair string) {
	pair1 := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "NEOUSDT", "XLMUSDT", "XRPUSDT"}
	fullTrans := AlgoMacdAP(float64(bot.WalletFuturesUSDT), float64(bot.FeeMaker), float64(bot.FeeTaker), bot.KlinesDataCsv, from, to, pair1)
	bot.TransactionFullData = fullTrans
	usdt := fullTrans.Finalusdt
	transactionAmount := fullTrans.Transamount
	successTrans := fullTrans.SuccessfulTrans
	lostTrans := fullTrans.LostTrans
	procent := (usdt/bot.WalletFuturesUSDT)*100 - 100
	fmt.Printf("USDT Left: %f\n", usdt)
	fmt.Printf("Gained percent: %f\n", procent)
	fmt.Printf("Number of entries into transaction: %d\n", transactionAmount)
	fmt.Printf("Number of succesful transactions: %d\n", successTrans)
	fmt.Printf("Number of lost transactions: %d\n", lostTrans)
	//Generating finalusdt chart)
	x := make([]string, 0)
	data := make([]opts.LineData, 0)
	for i := 0; i < len(fullTrans.Transaction); i++ {
		x = append(x, TimestampToDate1(fullTrans.Transaction[i].ClosingTime))
		data = append(data, opts.LineData{Value: fullTrans.Transaction[i].FinishBalance})
	}
	surface3d := charts.NewLine()
	surface3d.SetGlobalOptions(
		charts.WithSingleAxisOpts(opts.SingleAxis{
			Type:   "time",
			Bottom: "10%",
		}),
		charts.WithTitleOpts(opts.Title{Title: "basic surface3D example"}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        220000,
			Min:        0,
		}),
	)

	surface3d.SetXAxis(x).AddSeries("Category A", data)

	//Adding Data to series
	name := "wykres-usdmonth.html"
	page := components.NewPage()
	page.AddCharts(surface3d)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	fmt.Println("its going to chart")
	page.Render(io.MultiWriter(f))
	defer f.Close()
}

//GetCandlesFromRange gives us candlesticks in the given range
func (bot *TestBot) GetCandlesFromRange(pair string, start string, stop string) []*futures.Kline {
	var ind int
	switch pair {
	case "BTCUSDT":
		ind = 0
	case "ETHUSDT":
		ind = 1
	case "ADAUSDT":
		ind = 2
	case "NEOUSDT":
		ind = 3
	case "XLMUSDT":
		ind = 4
	case "XRPUSDT":
		ind = 5
	case "LINKUSDT":
		ind = 6
	default:
		ind = 0
	}
	startTimestamp := DateToTimestampRange(start)
	stopTimestamp := DateToTimestampRange(stop)
	var interstingKlines []*futures.Kline
	for _, v := range bot.KlinesDataCsv[ind][0] {
		if v.OpenTime >= startTimestamp && v.OpenTime <= stopTimestamp {
			interstingKlines = append(interstingKlines, v)
		}

	}

	return interstingKlines
}

//Chart of worthness
func (bot *TestBot) GenerateChartsFromAlgorithm(pair string, from string, to string) {
	fmt.Println("Chart genreation started")
	var tp []float64
	var pro []float64
	var value []int
	pair1 := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "NEOUSDT", "XLMUSDT", "XRPUSDT"}
	for i := 1005; i <= 1050; i = i + 5 {
		pro = append(pro, float64(float64(i)/1000.00))
	}
	for i := 15; i <= 100; i = i + 10 {
		tp = append(tp, float64(float64(i)/100.00))
	}
	for i := 10; i <= 30; i = i + 3 {
		value = append(value, i)
	}
	for s := 0; s < len(value); s++ {
		data := make([][3]interface{}, 0)
		min := 10000000000.00
		max := 0.00
		for i := 0; i < len(pro); i++ {
			for o := 0; o < len(tp); o++ {
				var fin float64 = 0
				//----------------------------------------//
				///INSERTING GIVEN ALGORITHM TO BE CHECKED//
				//----------------------------------------//
				fullTrans := AlgoRsiAPC(float64(bot.WalletFuturesUSDT), float64(bot.FeeMaker), float64(bot.FeeTaker), bot.KlinesDataCsv, pair1, from, to, tp[o], pro[i], value[s])
				fin += fullTrans.Finalusdt
				fmt.Println(value[s], pro[i], tp[o], fullTrans.Transamount, fullTrans.SuccessfulTrans)
				data = append(data, [3]interface{}{pro[i], tp[o], fin})
				if fin < min {
					min = fin
				}
				if fin > max {
					max = fin
				}
			}
		}
		surface3d := charts.NewSurface3D()
		surface3d.SetGlobalOptions(
			charts.WithZAxis3DOpts(opts.ZAxis3D{Min: min, Max: max}),
			charts.WithXAxis3DOpts(opts.XAxis3D{Min: pro[0], Max: pro[len(pro)-1]}),
			charts.WithYAxis3DOpts(opts.YAxis3D{Min: tp[len(tp)-1], Max: tp[0]}),
			charts.WithTitleOpts(opts.Title{Title: "basic surface3D example"}),
			charts.WithVisualMapOpts(opts.VisualMap{
				Calculable: true,
				Max:        0,
				Min:        10000,
			}),
		)

		//Inserting Data into series
		ret := make([]opts.Chart3DData, 0)
		for _, d := range data {
			ret = append(ret, opts.Chart3DData{
				Value: []interface{}{d[0], d[1], d[2]},
			})
		}

		//Adding Data to series
		name := "wykres-signal" + strconv.Itoa(value[s]) + ".html"
		surface3d.AddSeries(name, ret)
		page := components.NewPage()
		page.AddCharts(surface3d)
		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		fmt.Println("its going to chart")
		page.Render(io.MultiWriter(f))
		defer f.Close()
	}
}
