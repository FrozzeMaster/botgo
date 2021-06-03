package testbot

import (
	"fmt"

	"github.com/CraZzier/bot/model"
	"github.com/adshao/go-binance/v2/futures"
)

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

//MinAndMax gets MAX and MIN from array
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

//Getting bot test data
func BotParameters(transactions *model.AlgorithmTransactionsFull) {

	longNumber := 0
	shortNumber := 0
	longWinNumber := 0
	shortWinNumber := 0
	winsInRow := 0
	lossInRow := 0
	shortWinsInRow := 0
	longWinsInRow := 0
	shortLossInRow := 0
	longLossInRow := 0
	timeInTransaction := 0
	shortTimeInTransaction := 0
	longTimeInTransaction := 0
	longtargetSmallerThan25 := 0
	shorttargetSmallerThan25 := 0

	fullFee := 0.000
	averageTransactionTime := 0.000
	longAverageTransactionTime := 0.000
	shortAverageTransactionTime := 0.000

	for _, v := range transactions.Transaction {
		timeInTransaction += int(v.ClosingTime - v.EntryTime)
		if v.LongOrShort == "SHORT" {
			longNumber++
			if v.Type == "Succesful" {
				longWinNumber++
			}
			if (v.SellingPrice-v.BuyingPrice)/v.BuyingPrice < 0.025 {
				longtargetSmallerThan25++
			}
			longTimeInTransaction += int(v.ClosingTime - v.EntryTime)
		} else {
			shortNumber++
			if v.Type == "Succesful" {
				shortWinNumber++
			}
			if (-v.SellingPrice+v.BuyingPrice)/v.BuyingPrice < 0.025 {
				shorttargetSmallerThan25++
			}
			shortTimeInTransaction += int(v.ClosingTime - v.EntryTime)
		}
		fullFee += v.EntryFee + v.ClosingFee

	}
	averageTransactionTime = float64(timeInTransaction / len(transactions.Transaction))
	longAverageTransactionTime = float64(longTimeInTransaction / longNumber)
	shortAverageTransactionTime = float64(shortTimeInTransaction / shortNumber)

	var winrowTable []int
	winrowTable = append(winrowTable, 0)
	var losrowTable []int
	losrowTable = append(losrowTable, 0)

	var winrowTableLong []int
	winrowTableLong = append(winrowTableLong, 0)
	var losrowTableShort []int
	losrowTableShort = append(losrowTableShort, 0)

	var winrowTableShort []int
	winrowTableShort = append(winrowTableShort, 0)
	var losrowTableLong []int
	losrowTableLong = append(losrowTableLong, 0)

	for _, v := range transactions.Transaction {
		if v.Type == "Succesful" {
			winrowTable = append(winrowTable, winrowTable[len(winrowTable)-1]+1)
			losrowTable = append(losrowTable, 0)
		} else {
			winrowTable = append(winrowTable, 0)
			losrowTable = append(losrowTable, losrowTable[len(losrowTable)-1]+1)
		}
	}
	for _, v := range transactions.Transaction {
		if v.LongOrShort == "Long" {
			if v.Type == "Succesful" {
				winrowTableLong = append(winrowTableLong, winrowTableLong[len(winrowTableLong)-1]+1)
				losrowTableLong = append(losrowTableLong, 0)
			} else {
				winrowTableLong = append(winrowTableLong, 0)
				losrowTableLong = append(losrowTableLong, losrowTableLong[len(losrowTableLong)-1]+1)
			}
		} else {
			if v.Type == "Succesful" {
				winrowTableShort = append(winrowTableShort, winrowTableShort[len(winrowTableShort)-1]+1)
				losrowTableShort = append(losrowTableShort, 0)
			} else {
				winrowTableShort = append(winrowTableShort, 0)
				losrowTableShort = append(losrowTableShort, losrowTableShort[len(losrowTableShort)-1]+1)
			}
		}
	}
	winsInRow = MinAndMax(winrowTable, "max")
	lossInRow = MinAndMax(losrowTable, "max")
	longWinsInRow = MinAndMax(winrowTableLong, "max")
	shortWinsInRow = MinAndMax(winrowTableShort, "max")
	longLossInRow = MinAndMax(losrowTableLong, "max")
	shortLossInRow = MinAndMax(losrowTableShort, "max")
	fmt.Printf("Liczba longów: %d \n", longNumber)
	fmt.Printf("Liczba shortow: %d \n", shortNumber)
	fmt.Printf("zwyciestwa na longach: %d \n", longWinNumber)
	fmt.Printf("zwyciestwa na shortach: %d \n", shortWinNumber)
	fmt.Printf("Max zwyciestwa z rzedu: %d \n", winsInRow)
	fmt.Printf("Max przegrane z rzedu: %d \n", lossInRow)
	fmt.Printf("Max zwyciestwa na shortach z rzedu: %d \n", shortWinsInRow)
	fmt.Printf("Max zwyciestwa na longach z rzedu: %d \n", longWinsInRow)
	fmt.Printf("Max przegrane na shortach z rzedu: %d \n", shortLossInRow)
	fmt.Printf("Max przegrane na longach z rzedu: %d \n", longLossInRow)
	fmt.Printf("Czas w transkacji w sumie: %f \n", float64(timeInTransaction)/60000)
	fmt.Printf("Czas w transkacji shorty: %f\n", float64(shortTimeInTransaction)/60000)
	fmt.Printf("Czas w transkacji longi: %f \n", float64(longTimeInTransaction)/60000)
	fmt.Printf("Transakcje na longach z targetem mniejszym niz niz 2,5 procent: %d \n", longtargetSmallerThan25)
	fmt.Printf("Transakcje na shortach z targetem mniejszym niz niz 2,5 procent: %d \n", shorttargetSmallerThan25)
	fmt.Printf("Suma fee: %f \n", fullFee)
	fmt.Printf("Sredni czas w transkacji: %f \n", averageTransactionTime/60000)
	fmt.Printf("sredni czas w transkacji na longach: %f \n", longAverageTransactionTime/60000)
	fmt.Printf("Sredni czas w transkacji na shortach: %f \n", shortAverageTransactionTime/60000)
	fmt.Println("Chart genreation started")
}
