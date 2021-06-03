# Ninety-Ninety

Ninety-Ninety is a crypto bot software made for generating
huge amounts of passive income. Best specialists are taking care of delivering the best possible strategies for algotrading

This is repository directed only for the backend of the application. Frontend can be found using name "botgofront"

## Installation

Install the latest version of [go](https://golang.org/dl/). Then change the binance API KEY nested in realbot package.

To run a program use: 

```bash
go get https://github.com/CraZzier/botgo

go run main.go
```

## Usage

Current module uses GOMODULES

Here are Api commands that will be in use for the development process
```go
    r.POST("backend/botCandles", api.BotCandles)
    r.POST("backend/botTest", api.BotTest)
    r.POST("backend/botChart", api.BotChart)
    r.GET("backend/init", api.InitBot)
    r.GET("backend/realBot", api.RealBot)
    r.Run("127.0.0.1:8080")
```
## Note
API KEY in module is not valid. Please change it to your own to use the programm.
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)