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
	}
}

func main() {
	ctx := context.Background()

	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	//client.SetLogDebug(true)

	// Получение информации о заявках (ордерах)
	portfolio := "D88833" // номер счета (spot)
	orders, err := client.GetOrders(ctx, portfolio)
	if err != nil {
		slog.Error("main.GetOrders", "err", err.Error())
		return
	}
	slog.Info("orders", "кол-во", len(orders))
	for n, order := range orders {
		slog.Info("order",
			"row", n,
			"order", slog.Any("o", order),
		)
	}

	// создать новую заявку

	//orderID, err := client.NewCreateOrderService().
	//	Symbol("SBER").
	//	Side(alor.SideTypeBuy).
	//	OrderType(alor.OrderTypeLimit).
	//	Qty(1).
	//	Price(320).
	//	Portfolio(portfolio).
	//	Comment("комментарий к сделке").
	//	Do(ctx)

	client.Portfolio = portfolio // номер счета должен быть указан в клиенте
	// покупка по рынку
	//orderID, err := client.BuyMarket(ctx, "NVTK", 1, "comment к сделке")
	// продажа по рынку
	//orderID, err := client.SellMarket(ctx, "MOEX", 1, "comment к сделке")

	// лимитная продажа
	//orderID, err := client.SellLimit(ctx, "NVTK", 1, 1228, "comment к сделке")
	// лимитная покупка
	//orderID, err := client.BuyLimit(ctx, "SBER", 1, 322.1, "comment к сделке")
	//
	//if err != nil {
	//	slog.Error("main.SendOrder", "err", err.Error())
	//	return
	//}
	//slog.Info("sendOrder", "orderID", orderID)

	// отменить ордер
	orderId := "47050802385"
	ok, err := client.CancelOrder(ctx, portfolio, orderId)
	if err != nil {
		slog.Error("main.CancelOrder", "err", err.Error(), "ok", ok)
		return
	}
	if ok {
		slog.Info("CancelOrder успешно выполнено")
	}

}
