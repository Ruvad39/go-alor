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
	GetSecurity(ctx context.Context, board, symbol string) (Security, error)
	// GetSecurities получить список торговых инструментов
	GetSecurities(ctx context.Context, params Params) ([]Security, error)
	// GetQuotes Получение информации о котировках для выбранных инструментов
	GetQuotes(ctx context.Context, symbols string) ([]Quote, error)
	// GetQuote Получение информации о котировках для одного выбранного инструмента.
	GetQuote(ctx context.Context, symbol string) (Quote, error)
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
