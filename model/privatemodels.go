package model

//SingleMovingAverageStamp keeps daa for single timestamp
type SingleMovingAverageStamp struct {
	Timestamp int64
	Value     float64
}
type SingleMACDStamp struct {
	Timestamp int64
	Value     []float64
}

//MovingAverage keeps moving average data
type MovingAverage struct {
	Keys           []*SingleMovingAverageStamp
	Pair           string
	StartTimestamp int64
	StopTimestamp  int64
	Interval       string
	IntervalValue  int64
	WhichValue     string
}
type RSI struct {
	Keys           []*RSIstamp
	Pair           string
	StartTimestamp int64
	StopTimestamp  int64
	Interval       string
	IntervalValue  int64
	WhichValue     string
}

//RSIstamp keeps data of single RSI candle
type RSIstamp struct {
	Timestamp int64
	Change    float64
	CurrGain  float64
	CurrLoss  float64
	AvgGain   float64
	AvgLoss   float64
	RS        float64
	RSI       float64
	Close     float64
}

//MACD to handle MACD
type MACD struct {
	Keys               []*SingleMACDStamp
	E1Keys             []*SingleMovingAverageStamp
	E2Keys             []*SingleMovingAverageStamp
	Pair               string
	StartTimestamp     int64
	StopTimestamp      int64
	Interval           string
	E1                 int64
	E2                 int64
	Signal             int64
	CandleTrueInterval int64
	WhichValue         string
}

//MyKline serves for handling data within one Kline
type MyKline struct {
	Open      float64
	Close     float64
	Min       float64
	Max       float64
	Volume    float64
	OpenTime  int64
	CloseTime int64
	Emas      []float64
	MacD      [][]float64
	RSI       float64
}

//TransNumbers keep transaction numbers
type TransNumbers struct {
	TransAmount  int
	SuccessTrans int
	LostTrans    int
	FeeSum       float64
}
