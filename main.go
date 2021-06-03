package main

import (
	"github.com/CraZzier/bot/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(api.CORS)
	r.POST("backend/botCandles", api.BotCandles)
	r.POST("backend/botTest", api.BotTest)
	r.POST("backend/botChart", api.BotChart)
	r.GET("backend/init", api.InitBot)
	r.GET("backend/realBot", api.RealBot)
	r.Run("127.0.0.1:8080")
}
