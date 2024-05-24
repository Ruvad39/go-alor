package main

import (
	"context"
	"github.com/Ruvad39/go-alor"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
		//slog.Error(err.Error())
	}
	//handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	//	Level: slog.LevelDebug,
	//})
	//slog.SetDefault(slog.New(handler))
}

func main() {
	ctx := context.Background()

	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	client.SetLogDebug(true)

	// Получение информации о котировках для выбранных инструментов.
	// Принимает несколько пар биржа-тикер. Пары отделены запятыми. Биржа и тикер разделены двоеточием
	symbols := "MOEX:SIM4,MOEX:SBER"
	sec, err := client.GetQuotes(ctx, symbols)
	if err != nil {
		slog.Info("main.GetQuotes", "err", err.Error())
		return
	}

	slog.Info("GetQuotes", "кол-во", len(sec))
	for _, q := range sec {
		slog.Info("Quotes",
			"symbol", q.Symbol,
			"description", q.Description,
			"lastPrice", q.LastPrice,
			"Bid", q.Bid,
			"ask", q.Ask,
			"LotValue", q.LotValue,
			"LotSize", q.LotSize,
			"OpenInterest", q.OpenInterest,
			"LastTime", q.LastTime().String(),
		)
	}

	// Получение информации о котировках для одного выбранного инструмента.
	// Указываем тикер без указания биржи. Название биржи берется по умолчанию
	symbol := "SRM4"
	q, err := client.GetQuote(ctx, symbol)
	slog.Info("Quotes",
		"symbol", q.Symbol,
		"description", q.Description,
		"lastPrice", q.LastPrice,
		"Bid", q.Bid,
		"ask", q.Ask,
		"LotValue", q.LotValue,
		"LotSize", q.LotSize,
		"OpenInterest", q.OpenInterest,
		"LastTime", q.LastTime().String(),
	)

	// получение информации о биржевом стакане
	//symbol = "SBER"
	orderbook, err := client.GetOrderBooks(ctx, symbol)
	if err != nil {
		slog.Info("main.GetOrderBooks", "err", err.Error())
		return
	}
	//slog.Info("GetOrderBooks", "MsTimestamp", orderbook.MsTimestamp)
	slog.Info("GetOrderBooks", "orderbook", orderbook.String())
	bid, _ := orderbook.BestBid()
	ask, _ := orderbook.BestAsk()
	slog.Info("orderbook", "BestBid()", bid.Price, "BestAsk()", ask.Price)
	slog.Info("orderbook", "объем bid", orderbook.Bids.SumDepth(), "объем ask", orderbook.Asks.SumDepth())

}
