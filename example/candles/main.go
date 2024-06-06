package main

import (
	"context"
	"fmt"
	"github.com/Ruvad39/go-alor"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
		//slog.Error(err.Error())
	}
}

func main() {
	ctx := context.Background()

	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	client.SetLogDebug(true)

	// получить список свечей по инструменту
	timeFrom, _ := time.Parse("2006-01-02", "2024-01-25")
	timeTo, _ := time.Parse("2006-01-02", "2024-06-01")

	//history, err := client.GetHistory(ctx, "SBER", alor.Interval_D1, timeFrom.Unix(), timeTo.Unix())
	//if err != nil {
	//	slog.Info("main.GetHistory", "err", err.Error())
	//	return
	//}
	//slog.Info("candles", "кол-во", len(history.Candles))
	//for n, candle := range history.Candles {
	//	slog.Info("candles",
	//		"row", n,
	//		"Time", candle.GeTime(),
	//		"close", candle.Close,
	//	)
	//}
	symbol := "SBER"
	interval := alor.Interval_D1
	candles, err := client.GetCandles(ctx, symbol, interval, timeFrom.Unix(), timeTo.Unix())
	if err != nil {
		slog.Info("main.GetCandles", "err", err.Error())
		return
	}
	slog.Info("candles", "кол-во", len(candles))

	c := alor.Candle{}
	fmt.Println(c.CsvHeader())
	for _, candle := range candles {
		candle.Symbol = symbol
		candle.Interval = interval
		//slog.Info("candle", "row", n, "Time", candle.GeTime(), "close", candle.Close)
		//slog.Info("candle", "row", candle.CsvRecords())
		fmt.Println(candle.CsvRecord())
	}

}
