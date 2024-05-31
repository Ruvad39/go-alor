package alor

import (
	"context"
	"fmt"
)

// SubscribeCandles подписка на свечи
func (c *Client) SubscribeCandles(ctx context.Context, symbol string, interval Interval, opts ...WSRequestOption) error {
	_, ok, err := c.GetSecurity(ctx, "", symbol)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("инструмент %s не найден", symbol)
	}
	r := &WSRequestBase{
		OpCode:      OnCandleSubscribe,
		Code:        symbol,
		Interval:    interval,
		SkipHistory: true,
		Exchange:    c.Exchange,
		Frequency:   1000, //  По умолчанию обновление в 1сек (1000 мс)
	}
	// обрабратаем входящие параметры
	for _, opt := range opts {
		opt(r)
	}
	r.Guid = "candle|" + r.Code + "|" + r.Interval.String()

	s := c.NewWsService(r)
	return s.Do(ctx)
}

// SubscribeQuotes подписка на котировки
func (c *Client) SubscribeQuotes(ctx context.Context, symbol string, opts ...WSRequestOption) error {
	_, ok, err := c.GetSecurity(ctx, "", symbol)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("инструмент %s не найден", symbol)
	}

	r := &WSRequestBase{
		OpCode:    onQuotesSubscribe,
		Code:      symbol,
		Exchange:  c.Exchange,
		Frequency: 175, //  По умолчанию 175 Минимальное значение параметра зависит от выбранного формата возвращаемого JSON-объекта: Slim — 10 миллисекунд
	}
	// обрабратаем входящие параметры
	for _, opt := range opts {
		opt(r)
	}
	// TODO создать метод создания guid
	r.Guid = "quote|" + r.Code

	s := c.NewWsService(r)
	return s.Do(ctx)
}

// SubscribeOrders подписка на получение информации обо всех биржевых заявках с участием указанного портфеля
func (c *Client) SubscribeOrders(ctx context.Context, portfolio string, opts ...WSRequestOption) error {

	r := &WSRequestBase{
		OpCode:    onOrdersSubscribe,
		Portfolio: portfolio,
		Exchange:  c.Exchange,
		//OrderStatuses:  ["filled"],
	}
	// обработаем входящие параметры
	for _, opt := range opts {
		opt(r)
	}
	// TODO создать метод создания guid
	r.Guid = "orders|" + r.Portfolio

	s := c.NewWsService(r)
	return s.Do(ctx)
}
