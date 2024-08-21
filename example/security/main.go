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
}

func main() {
	ctx := context.Background()

	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	//client.SetLogDebug(true)

	/*
		// получить параметры по торговому инструменту
		board := "" // "TQBR"
		symbol := "Si-9.24"
		sec, ok, err := client.GetSecurity(ctx, board, symbol)
		if err != nil {
			slog.Info("main.GetSecurity", "err", err.Error())
			return
		}
		if !ok {
			slog.Info("main.GetSecurity", symbol, "нет такого инструмента")
			//return
		}
		slog.Info("symbol", slog.Any("sec", sec))

		slog.Info("security",
			"Symbol", sec.Symbol,
			"Exchange", sec.Exchange,
			"board", sec.Board,
			"ShortName", sec.ShortName,
			"LotSize", sec.LotSize,
			"MinStep", sec.MinStep,
			"Cancellation", sec.Cancellation,
		)
	*/
	// Если не указано иное значение параметра limit, в ответе возвращается только 25 объектов за раз
	//Sec, err := client.GetSecurities(ctx,
	//	alor.WithSector("FOND"),
	//	alor.WithBoard("TQBR"),
	//	alor.WithLimit(5),
	//)
	// ROPD опционы
	Sec, err := client.GetSecurities(ctx,
		alor.WithSector("FORTS"),
		alor.WithBoard("ROPD"),
		alor.WithLimit(100),
		alor.WithOld(true),
	)

	if err != nil {
		slog.Info("main.GetSecurity", "err", err.Error())
		return
	}
	slog.Info("GetSecurity", "кол-во", len(Sec))
	for n, sec := range Sec {
		slog.Info("securities",
			"row", n,
			"Symbol", sec.Symbol,
			"Exchange", sec.Exchange,
			"board", sec.Board,
			"ShortName", sec.ShortName,
			"type", sec.Type,
			"LotSize", sec.LotSize,
			"MinStep", sec.MinStep,
			"Cancellation", sec.Cancellation,
			"OptionSide", sec.OptionSide,
			"StrikePrice", sec.StrikePrice,
		)
	}

}
