package model

//CoinFormula is a model for chartData from Frontend
type CoinFormula struct {
	CoinName string
	From     string
	To       string
}

//AlgorithmTransactionsFull have all one test data
type AlgorithmTransactionsFull struct {
	Transaction     []*AlgorithmTransactions
	Finalusdt       float64
	Transamount     int64
	SuccessfulTrans int64
	LostTrans       int64
}

//AlgorithmTransactions keeps packet of trasnactions main data
type AlgorithmTransactions struct {
	StartBalance    float64
	EntryFee        float64
	Stoploss        []*Stoploss
	Leverage        float64
	BuyingPrice     float64
	BalanceMinusFee float64
	EntryTime       int64
	LongOrShort     string
	Target          float64

	FinishBalance float64
	ClosingFee    float64
	Type          string
	ClosingTime   int64
	SellingPrice  float64

	FeeSum           float64
	ProfitWithoutFee float64
	Profit           float64
	Percent          float64
	Pair             string
}

//Stoploss is a struct that keeps stolosess
type Stoploss struct {
	Timestamp int64
	Value     float64
}
