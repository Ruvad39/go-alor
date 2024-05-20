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
	client.Debug = true

	// Получение информации о портфеле
	portfolio := "D88833" // номер счета (срочный)
	p, err := client.GetPortfolio(ctx, portfolio)
	if err != nil {
		slog.Info("main.GetPortfolio", "err", err.Error())
		return
	}
	slog.Info("portfolio", slog.Any("p", p))

	// получение информации о позициях
	positions, err := client.GetPositions(ctx, portfolio)
	if err != nil {
		slog.Info("main.GetPositions", "err", err.Error())
		return
	}

	slog.Info("GetPosition", "кол-во", len(positions))
	for n, pos := range positions {
		slog.Info("Positions",
			"row", n,
			slog.Any("pos", pos),
		)
	}

	//slog.Info("GetLoginPositions", "кол-во", len(positions))
	//for n, pos := range positions {
	//	slog.Info("LoginPositions",
	//		"row", n,
	//		slog.Any("pos", pos),
	//	)
	//}
	// получение информации о позициях заданного инструмента
	// выдает "HTTP 404: Not Found" если нет позиций
	symbol := "MOEX" //"CRM4"
	position, ok, err := client.GetPosition(ctx, portfolio, symbol)
	if err != nil {
		slog.Info("main.GetPosition", "err", err.Error())
		return
	}
	if !ok {
		slog.Info("main.GetPosition", symbol, "нет позиции")
		//return
	}
	slog.Info("Position", slog.Any("position", position))
	symbol = "RUB" //денежные средства
	position, _, err = client.GetPosition(ctx, portfolio, symbol)
	if err != nil {
		slog.Info("main.GetPosition", "err", err.Error())
		return
	}
	slog.Info("Position", "Свободные средства", position.Qty)

}
