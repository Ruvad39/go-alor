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
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	// создание клиента
	client := alor.NewClient(refreshToken)
	//client.SetLogDebug(true)

	// Получение информации о заявках (ордерах)
	portfolio := "D88833" // номер счета (spot)

	// добавим колл-бэк на событие появление заявки
	client.SetOnOrder(OnOrder)

	// подпишемся на появление заявки
	err := client.SubscribeOrders(ctx, portfolio)
	if err != nil {
		slog.Error("SubscribeOrders", "err", err.Error())
		return
	}

	orders, err := client.GetOrders(ctx, portfolio)
	if err != nil {
		slog.Error("main.GetOrders", "err", err.Error())
		return
	}
	slog.Info("GetOrders", "кол-во", len(orders))
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

	client.SetPortfolioID(portfolio) // номер счета для работы по умолчанию
	// покупка по рынку
	//orderID, err := client.BuyMarket(ctx, "SBER", 1, "comment к сделке")
	// продажа по рынку
	//orderID, err := client.SellMarket(ctx, "SBER", 1, "comment к сделке")

	// лимитная продажа
	//orderID, err := client.SellLimit(ctx, "NVTK", 1, 1228, "comment к сделке")
	// лимитная покупка
	//orderID, err := client.BuyLimit(ctx, "SBER", 1, 311.1, "comment к сделке")
	////
	//if err != nil {
	//	slog.Error("main.SendOrder", "err", err.Error())
	//	//return
	//}
	//slog.Info("sendOrder", "orderID", orderID)

	// отменить ордер
	orderId := "48072118608"
	ok, err := client.CancelOrder(ctx, portfolio, orderId)
	if err != nil {
		slog.Error("main.CancelOrder", "err", err.Error(), "ok", ok)
		return
	}
	if ok {
		slog.Info("CancelOrder успешно выполнено")
	}

	//----------------------------------
	// ожидание сигнала о закрытие
	waitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	cancel()

	slog.Info("exiting...")

}

// события появления заявки
func OnOrder(order alor.Order) {
	slog.Info("OnOrder", slog.Any("order", order))
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
