package candlesticks

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
)

//CandlestickName returns the name of given calndlestick
func CandlestickName(openentry string, closeentry string, maxentry string, minentry string) (string, int) {
	open, err := strconv.ParseFloat(openentry, 32)
	if err != nil {
		fmt.Println(err)
	}
	close, err := strconv.ParseFloat(closeentry, 32)
	if err != nil {
		fmt.Println(err)
	}
	min, err := strconv.ParseFloat(minentry, 32)
	if err != nil {
		fmt.Println("Couldnt convert to flow")
	}
	max, err := strconv.ParseFloat(maxentry, 32)
	if err != nil {
		fmt.Println("Couldnt convert to flow")
	}
	//Greencandles
	if close > open {
		colour := 1
		size := max - min
		body := (close - open) / size
		uptail := (max - close) / size
		downtail := (open - min) / size
		if uptail+downtail < body {
			return "Green3Candle", colour
		}
		//DojiLike
		if body < 0.20 {
			return "DojiLike", colour
		}

	} else {
		//Redcandles
		colour := -1
		size := max - min
		body := (open - close) / size
		uptail := (max - open) / size
		downtail := (close - min) / size
		if uptail < 0.4*downtail && body < 0.4 {
			return "BearishCross", colour
		}
	}
	return "Undefined", 0
}

//TimestampToDateFile converts for format in csv to date
func TimestampToDateFile(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond)).UTC()
	strDate := t.Format("2006-01-02 15:04:05")
	return strDate
}

func Fix1mCandles(candlesticks []*futures.Kline, interval string, pair string) {
	var data [][]string
	for i := 0; i < len(candlesticks)-2; i++ {
		if i == 0 {
			var oneline []string
			oneline = append(oneline, "open_time")
			oneline = append(oneline, "open")
			oneline = append(oneline, "high")
			oneline = append(oneline, "low")
			oneline = append(oneline, "close")
			oneline = append(oneline, "volume")
			oneline = append(oneline, "quote_asset_volume")
			oneline = append(oneline, "number_of_trades")
			oneline = append(oneline, "taker_buy_base_asset_volume")
			oneline = append(oneline, "taker_buy_quote_asset_volume")
			data = append(data, oneline)
		}
		var oneline []string
		if candlesticks[i].OpenTime != candlesticks[i+1].OpenTime-60000 {
			takeOpen := candlesticks[i].OpenTime
			for takeOpen <= candlesticks[i+1].OpenTime-60000 {
				fmt.Println("found: " + pair + " " + TimestampToDateFile(takeOpen) + " " + fmt.Sprintf("%d", takeOpen))
				var oneline1 []string
				oneline1 = append(oneline1, TimestampToDateFile(takeOpen))
				oneline1 = append(oneline1, candlesticks[i].Open)
				oneline1 = append(oneline1, candlesticks[i].High)
				oneline1 = append(oneline1, candlesticks[i].Low)
				oneline1 = append(oneline1, candlesticks[i].Close)
				oneline1 = append(oneline1, candlesticks[i].Volume)
				oneline1 = append(oneline1, candlesticks[i].QuoteAssetVolume)
				oneline1 = append(oneline1, fmt.Sprintf("%d", candlesticks[i].TradeNum))
				oneline1 = append(oneline1, candlesticks[i].TakerBuyBaseAssetVolume)
				oneline1 = append(oneline1, candlesticks[i].TakerBuyQuoteAssetVolume)
				data = append(data, oneline1)
				takeOpen += 60000
			}
		} else {
			oneline = append(oneline, TimestampToDateFile(candlesticks[i].OpenTime))
			oneline = append(oneline, candlesticks[i].Open)
			oneline = append(oneline, candlesticks[i].High)
			oneline = append(oneline, candlesticks[i].Low)
			oneline = append(oneline, candlesticks[i].Close)
			oneline = append(oneline, candlesticks[i].Volume)
			oneline = append(oneline, candlesticks[i].QuoteAssetVolume)
			oneline = append(oneline, fmt.Sprintf("%d", candlesticks[i].TradeNum))
			oneline = append(oneline, candlesticks[i].TakerBuyBaseAssetVolume)
			oneline = append(oneline, candlesticks[i].TakerBuyQuoteAssetVolume)
			data = append(data, oneline)
		}
	}
	basepath := path.Dir("testbot/candledata/")
	pair = pair + interval + ".csv"
	filename := path.Base(pair)
	file, err := os.Create(strings.Join([]string{basepath, filename}, "/"))
	if err != nil {
		fmt.Println("error while creating file")
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//ConvertInterval converts cnadlesticks into diffrent timeframe
func ConvertInterval(candlesticks []*futures.Kline, interval string, pair string) {
	var miliseconds int64
	switch interval {
	case "5m":
		miliseconds = 60000 * 5
	case "15m":
		miliseconds = 60000 * 15
	case "1h":
		miliseconds = 60000 * 60
	case "1d":
		miliseconds = 60000 * 60 * 24
	case "1w":
		miliseconds = 60000 * 60 * 24 * 7
	default:
		miliseconds = 60000 * 60
	}

	var newcandles []*futures.Kline
	for i := 0; i < len(candlesticks)-1; i++ {
		max, _ := strconv.ParseFloat(candlesticks[i].High, 32)
		min, _ := strconv.ParseFloat(candlesticks[i].Low, 32)
		open, _ := strconv.ParseFloat(candlesticks[i].Open, 32)
		if candlesticks[i].OpenTime%miliseconds == 0 {
			dayopen := open
			minday := min
			maxday := max
			dayopenTime := candlesticks[i].OpenTime
			volume, _ := strconv.ParseFloat(candlesticks[i].Volume, 32)
			for candlesticks[i].CloseTime%miliseconds != miliseconds-1 && i < len(candlesticks)-1 {
				maxtemp, _ := strconv.ParseFloat(candlesticks[i].High, 32)
				mintemp, _ := strconv.ParseFloat(candlesticks[i].Low, 32)
				voltemp, _ := strconv.ParseFloat(candlesticks[i].Volume, 32)
				volume += voltemp
				if maxtemp > maxday {
					maxday = maxtemp
				}
				if mintemp < minday {
					minday = mintemp
				}
				i++
			}
			maxtemp, _ := strconv.ParseFloat(candlesticks[i].High, 32)
			mintemp, _ := strconv.ParseFloat(candlesticks[i].Low, 32)
			voltemp, _ := strconv.ParseFloat(candlesticks[i].Volume, 32)
			volume += voltemp
			if maxtemp > maxday {
				maxday = maxtemp
			}
			if mintemp < minday {
				minday = mintemp
			}
			daycloseTime := candlesticks[i].CloseTime
			dayclose := candlesticks[i].Close
			singleCandle := &futures.Kline{
				OpenTime:                 dayopenTime,
				Open:                     fmt.Sprintf("%f", dayopen),
				High:                     fmt.Sprintf("%f", maxday),
				Low:                      fmt.Sprintf("%f", minday),
				Close:                    dayclose,
				Volume:                   fmt.Sprintf("%f", volume),
				CloseTime:                daycloseTime,
				QuoteAssetVolume:         "nd",
				TradeNum:                 0,
				TakerBuyBaseAssetVolume:  "nd",
				TakerBuyQuoteAssetVolume: "nd",
			}
			newcandles = append(newcandles, singleCandle)
		}
	}
	var data [][]string
	for x := 0; x < len(newcandles); x++ {
		if x == 0 {
			var oneline []string
			oneline = append(oneline, "open_time")
			oneline = append(oneline, "open")
			oneline = append(oneline, "high")
			oneline = append(oneline, "low")
			oneline = append(oneline, "close")
			oneline = append(oneline, "volume")
			oneline = append(oneline, "quote_asset_volume")
			oneline = append(oneline, "number_of_trades")
			oneline = append(oneline, "taker_buy_base_asset_volume")
			oneline = append(oneline, "taker_buy_quote_asset_volume")
			data = append(data, oneline)
		}
		var oneline []string
		oneline = append(oneline, TimestampToDateFile(newcandles[x].OpenTime))
		oneline = append(oneline, newcandles[x].Open)
		oneline = append(oneline, newcandles[x].High)
		oneline = append(oneline, newcandles[x].Low)
		oneline = append(oneline, newcandles[x].Close)
		oneline = append(oneline, newcandles[x].Volume)
		oneline = append(oneline, newcandles[x].QuoteAssetVolume)
		oneline = append(oneline, fmt.Sprintf("%d", newcandles[x].TradeNum))
		oneline = append(oneline, newcandles[x].TakerBuyBaseAssetVolume)
		oneline = append(oneline, newcandles[x].TakerBuyQuoteAssetVolume)
		data = append(data, oneline)
	}
	basepath := path.Dir("testbot/candledata/")
	pair = pair + interval + ".csv"
	filename := path.Base(pair)
	file, err := os.Create(strings.Join([]string{basepath, filename}, "/"))
	if err != nil {
		fmt.Println("error while creating file")
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			fmt.Println(err)
		}
	}
}

//ConvertInterval converts cnadlesticks into diffrent timeframe
func ConvertToBinary(candlesticks []*model.MyKline, deviation float64) {
	txt := ""
	reminder := 0.0
	for i, candlestick := range candlesticks {
		if candlestick.Close >= candlestick.Open {
			riseToMax := (candlestick.Max / candlestick.Open) - 1
			howMany1 := int(riseToMax / deviation)
			reminder += riseToMax/deviation - float64(howMany1)
			txt += strings.Repeat("1", howMany1+int(reminder))
			if reminder >= 1 {
				reminder -= 1
			}
			if candlestick.Close < candlestick.Max {
				fallToClose := (candlestick.Max / candlestick.Close) - 1
				howMany0 := int(fallToClose / deviation)
				reminder -= fallToClose/deviation - float64(howMany0)
				txt += strings.Repeat("0", howMany0+int(-reminder))
				if reminder <= -1 {
					reminder += 1
				}
			}
		} else {
			fallToMin := (candlestick.Open / candlestick.Min) - 1
			howMany0 := int(fallToMin / deviation)
			reminder -= fallToMin/deviation - float64(howMany0)
			txt += strings.Repeat("0", howMany0+int(-reminder))
			if reminder <= -1 {
				reminder += 1
			}
			if candlestick.Close > candlestick.Min {
				riseToClose := (candlestick.Close / candlestick.Min) - 1
				howMany1 := int(riseToClose / deviation)
				reminder += riseToClose/deviation - float64(howMany1)
				txt += strings.Repeat("1", howMany1+int(reminder))
				if reminder >= 1 {
					reminder -= 1
				}
			}
		}
		if i%1000 == 0 {
			// txt += TimestampToDate(candlestick.OpenTime) + "\n"
		}
	}
	f, err := os.Create("data.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(txt)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}

//TimestampToDate Timestamp(miliseconds) -> 02/01/2006 15:04:05
func TimestampToDate(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Millisecond))
	strDate := t.Format("02/01/2006 15:04:05")
	return strDate
}
