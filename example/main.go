package main

import (
	"context"
	"github.com/Ruvad39/go-alor"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
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
	//var _logger = log.Logger{
	//	Level:      log.DebugLevel,
	//	TimeField:  "time",
	//	TimeFormat: "2006-01-02 15:04:05.999Z07:00",
	//	Caller:     0,
	//	Writer: &log.MultiEntryWriter{
	//		&log.ConsoleWriter{ColorOutput: true},
	//		&log.FileWriter{Filename: "logs/test_candle.log", MaxSize: 10 * 1024 * 1024},
	//	},
	//}

	var logger *slog.Logger = (&log.Logger{
		//Level: log.DebugLevel,
		Level:      log.InfoLevel,
		TimeField:  "time",
		TimeFormat: "2006-01-02 15:04:05.999Z07:00",
		Caller:     0,
		Writer: &log.MultiEntryWriter{
			&log.ConsoleWriter{ColorOutput: true},
			&log.FileWriter{Filename: "logs/test_candle.log", MaxSize: 10 * 1024 * 1024},
		},
	}).Slog()
	slog.SetDefault(logger)
	//slog.SetDefault((_logger).Slog())

	var logClient = slog.With(slog.String("package", "go-alor"))
	alor.SetLogger(logClient)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	slog.Info("start main ")
	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	//client.SetLogDebug(true)

	// добавим коллбек на событие появление новой свечи
	client.RegisterOnCandleClosed(func(candle alor.Candle) {
		onCandle(candle)
	})

	//подписка на свечи

	// через создание сервиса
	//err := client.NewWSCandleService().Symbol("SBER").Interval(alor.Interval_M1).Do2(ctx)
	//err := client.NewWSCandleService("SBER", alor.Interval_M1).Do(ctx)
	err := client.SubscribeCandles(ctx, "SBER", alor.Interval_M1)
	if err != nil {
		slog.Error("main.NewWSCandleService", "err", err.Error())
		return
	}
	// через метод
	_ = client.SubscribeCandles(ctx, "SBER", alor.Interval_H1)
	_ = client.SubscribeCandles(ctx, "SBER", alor.Interval_D1)

	_ = client.SubscribeCandles(ctx, "Si-6.24", alor.Interval_M1)
	_ = client.SubscribeCandles(ctx, "Si-6.24", alor.Interval_H1)
	_ = client.SubscribeCandles(ctx, "Si-6.24", alor.Interval_D1)

	_ = client.SubscribeCandles(ctx, "MIX-6.24", alor.Interval_H1)
	_ = client.SubscribeCandles(ctx, "MIX-6.24", alor.Interval_M1)
	_ = client.SubscribeCandles(ctx, "MIX-6.24", alor.Interval_D1)

	_ = client.SubscribeCandles(ctx, "LKOH", alor.Interval_H1)
	_ = client.SubscribeCandles(ctx, "LKOH", alor.Interval_M1)
	_ = client.SubscribeCandles(ctx, "LKOH", alor.Interval_D1)

	_ = client.SubscribeCandles(ctx, "ROSN", alor.Interval_H1)
	_ = client.SubscribeCandles(ctx, "ROSN", alor.Interval_M1)
	_ = client.SubscribeCandles(ctx, "ROSN", alor.Interval_D1)

	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()

	slog.Info("exiting...")
}

// Ожидание сигнала о закрытие
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
