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

	// получить параметры по торговому инструменту
	board := "" // "TQBR"
	symbol := "SIM4"
	sec, err := client.GetSecurity(ctx, board, symbol)

	if err != nil {
		slog.Info("main.GetSecurity", "err", err.Error())
		return
	}
	slog.Info("symbol", slog.Any("sec", sec))

	slog.Info("security",
		"Symbol", sec.Symbol,
		"Exchange", sec.Exchange,
		"board", sec.Board,
		"ShortName", sec.ShortName,
		"LotSize", sec.LotSize,
		"MinStep", sec.MinStep,
	)
	//return
	// запрос списка инструментов
	// sector = FORTS, FOND, CURR
	// Если не указано иное значение параметра limit, в ответе возвращается только 25 объектов за раз
	params := alor.Params{
		Sector: "FOND",
		Board:  "TQBR",
		Query:  "",
		Limit:  400,
	}
	Sec, err := client.GetSecurities(ctx, params)
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
			"LotSize", sec.LotSize,
			"MinStep", sec.MinStep,
		)
	}

}
