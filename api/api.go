package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CraZzier/bot/model"
	RB "github.com/CraZzier/bot/realbot"
	TB "github.com/CraZzier/bot/testbot"
	"github.com/gin-gonic/gin"
)

var testbot *TB.TestBot
var realbot *RB.RealBot

//CORS accept cross origin requests
func CORS(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

//--------------------------------------------------------------------//
//------------------TESTBOT AND CANDLES-------------------------------//
//--------------------------------------------------------------------//

//InitBot downloads csv DATA and sets defualt variables
func InitBot(c *gin.Context) {
	testbot = new(TB.TestBot)
	testbot.Initialization()
	c.Writer.Header().Set("Content-type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode("Initiated")
}

//BotCandles sends to browser candlesticks
func BotCandles(c *gin.Context) {
	var rD model.CoinFormula
	err := json.NewDecoder(c.Request.Body).Decode(&rD)
	if err != nil {
		log.Fatal(err)
	}
	KlinesToSend := testbot.GetCandlesFromRange(rD.CoinName, rD.From, rD.To)
	c.Writer.Header().Set("Content-type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(KlinesToSend)
}

//BotChart generates visualization for algorithm parameteres
func BotChart(c *gin.Context) {
	var rD model.CoinFormula
	err := json.NewDecoder(c.Request.Body).Decode(&rD)
	if err != nil {
		log.Fatal(err)
	}
	testbot.GenerateChartsFromAlgorithm(rD.CoinName, rD.From, rD.To)
	c.Writer.Header().Set("Content-type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode("Generating chart")
}

//BotTest simulates transaction made by bot and returns list of them with description
func BotTest(c *gin.Context) {
	var rD model.CoinFormula
	err := json.NewDecoder(c.Request.Body).Decode(&rD)
	if err != nil {
		log.Fatal(err)
	}
	testbot.TestAlgorithm(rD.From, rD.To, rD.CoinName)
	c.Writer.Header().Set("Content-type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(testbot.TransactionFullData)
}

//-------------------------------------------------------------------//
//--------------------------REALBOT----------------------------------//
//-------------------------------------------------------------------//

//RealBot starts REAL BOT that trades on binance
func RealBot(c *gin.Context) {
	pairs := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "NEOUSDT", "XLMUSDT", "XRPUSDT"}
	intervals := []string{"5m", "15m", "1h"}
	realbot = new(RB.RealBot)
	realbot.Initialization(pairs, intervals, 1000)
	c.Writer.Header().Set("Content-type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode("Initiated")
}
