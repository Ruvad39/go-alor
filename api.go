package alor

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// какой api реализован
type IAlorClient interface {
	// GetTime текущее время сервера
	GetTime(ctx context.Context) (time.Time, error)

	// GetSecurity получить параметры по торговому инструменту
	GetSecurity(ctx context.Context, board, symbol string) (Security, bool, error)

	// GetSecurities получить список торговых инструментов
	GetSecurities(ctx context.Context, opts ...Option) ([]Security, error)

	// GetQuotes Получение информации о котировках для выбранных инструментов
	GetQuotes(ctx context.Context, symbols string) ([]Quote, error)

	// GetQuote Получение информации о котировках для одного выбранного инструмента
	GetQuote(ctx context.Context, symbol string) (Quote, error)

	// GetPositions получение информации о позициях
	GetPositions(ctx context.Context, portfolio string) ([]Position, error)

	// GetPosition Получение информации о позициях выбранного инструмента
	GetPosition(ctx context.Context, portfolio, symbol string) (Position, bool, error)

	// GetHistory Запрос истории для выбранных биржи и инструмента
	GetHistory(ctx context.Context, symbol string, interval Interval, from, to int64) (History, error)

	// GetCandles Запрос истории свечей для выбранного инструмента (вызывает GetHistory)
	GetCandles(ctx context.Context, symbol string, interval Interval, from, to int64) ([]Candle, error)

	// GetOrderBooks Получение информации о биржевом стакане
	GetOrderBooks(ctx context.Context, symbol string) (OrderBook, error)

	// GetOrders получение информации о всех заявках
	GetOrders(ctx context.Context, portfolio string) ([]Order, error)

	// GetOrder получение информации о выбранной заявке
	GetOrder(ctx context.Context, portfolio, orderId string) (Order, error)

	// SendOrder создать новый ордер
	//SendOrder(ctx context.Context, order OrderRequest) (OrderResponse, error)

	//BuyMarket покупка по рынку
	BuyMarket(ctx context.Context, symbol string, lot int32, comment string) (string, error)

	// SellMarket продажа по рынку
	SellMarket(ctx context.Context, symbol string, lot int32, comment string) (string, error)

	// BuyLimit лимитная покупка
	BuyLimit(ctx context.Context, symbol string, lot int32, price float64, comment string) (string, error)

	// SellLimit лимитная продажа
	SellLimit(ctx context.Context, symbol string, lot int32, price float64, comment string) (string, error)

	// CancelOrder отменить заявку
	CancelOrder(ctx context.Context, portfolio, orderId string) (bool, error)

	// SubscribeCandles подписка на свечи
	SubscribeCandles(ctx context.Context, symbol string, interval Interval, opts ...WSRequestOption) error

	// SubscribeQuotes подписка на котировки
	SubscribeQuotes(ctx context.Context, symbol string, opts ...WSRequestOption) error
}

// GetTime
// Запрос текущего UTC времени в формате Unix Time Seconds.
// Если этот запрос выполнен без авторизации, то будет возвращено время, которое было 15 минут назад.
func (c *Client) GetTime(ctx context.Context) (time.Time, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/md/v2/time",
	}
	//t := time.Now()

	data, err := c.callAPI(ctx, r)
	if err != nil {
		return time.Now(), err
	}
	timeUnix, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return time.Now(), err
	}
	t := time.Unix(timeUnix, 0)
	return t, err
}
