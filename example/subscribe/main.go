package main

import (
	"context"
	"github.com/Ruvad39/go-alor"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
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

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	slog.Info("start main ")
	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	//client.SetLogDebug(true)

	// добавим коллбек на событие появление новой свечи
	client.SetOnCandle(onCandle)
	// добавим коллбек на котировки
	client.SetOnQuote(onTick)

	// подписка на свечи
	err := client.SubscribeCandles(ctx, "SBRF-6.24", alor.Interval_M1, alor.WithFrequency(500))
	if err != nil {
		slog.Error("SubscribeCandles", "err", err.Error())
		return
	}
	_ = client.SubscribeCandles(ctx, "Si-6.24", alor.Interval_M1, alor.WithFrequency(500))

	// подписка на Котировки
	err = client.SubscribeQuotes(ctx, "Si-6.24", alor.WithFrequency(25))
	if err != nil {
		slog.Error("SubscribeQuotes", "err", err.Error())
		return
	}
	//_ = client.SubscribeQuotes(ctx, "SBRF-6.24")
	//_ = client.SubscribeQuotes(ctx, "SBER")

	//----------------------------------
	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()

	slog.Info("exiting...")
}

// сюда приходят данные по закрытым свечам
func onCandle(candle alor.Candle) {
	slog.Info("onCandle ",
		"symbol", candle.Symbol,
		"tf", candle.Interval.String(),
		"time", candle.GeTime().String(),
		"open", candle.Open,
		"high", candle.High,
		"low", candle.Low,
		"close", candle.Close,
		"volume", candle.Volume,
	)
}

func onTick(quote alor.Quote) {
	slog.Info("onTick",
		"symbol", quote.Symbol,
		"LastPriceTimestamp", quote.LastTime(),
		//"OrderBookMSTimestamp", quote.OrderBookMSTimestamp,
		"Bid", quote.Bid,
		"Ask", quote.Ask,
		"LastPrice", quote.LastPrice,
		"OpenInterest", quote.OpenInterest,
	)
}

// waitForSignal Ожидание сигнала о закрытие
func waitForSignal(ctx context.Context, signals ...os.Signal) os.Signal {
	var exit = make(chan os.Signal, 1)
	signal.Notify(exit, signals...)
	defer signal.Stop(exit)

	select {
	case sig := <-exit:
		slog.Info("WaitForSignal", "signals", sig)
		return sig
	case <-ctx.Done():
		return nil
	}

	return nil
}
