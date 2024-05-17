package alor

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

//func (c *Client) GetCandles(ctx context.Context, login string) ([]Position, error) {

// https://apidev.alor.ru/md/v2/history?symbol=SBER&exchange=MOEX&tf=D&from=1549000661&to=1634256000&format=Simple

// GetHistory Запрос истории для выбранных биржи и инструмента
// биржу берем по умолчанию
func (c *Client) GetHistory(ctx context.Context, symbol string, interval Interval, from, to int64) (History, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/md/v2/history",
	}
	r.setParam("exchange", c.Exchange)
	r.setParam("symbol", symbol)
	r.setParam("tf", interval)
	r.setParam("from", from)
	r.setParam("to", to)

	result := History{}
	data, err := c.callAPI(ctx, r)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil

}

// GetCandles Запрос свечей для выбранного инструмента
// биржу берем по умолчанию
func (c *Client) GetCandles(ctx context.Context, symbol string, interval Interval, from, to int64) ([]Candle, error) {
	history, err := c.GetHistory(ctx, symbol, interval, from, to)
	return history.Candles, err
}

// Candle Параметры свечи
type Candle struct {
	Time   int64   `json:"time"`   // Время (UTC) (Unix time seconds)
	Close  float64 `json:"close"`  // Цена при закрытии
	Open   float64 `json:"open"`   // Цена при открытии
	High   float64 `json:"high"`   // Максимальная цена
	Low    float64 `json:"low"`    // Минимальная цена
	Volume int32   `json:"volume"` // Объём
}

// GeTime вернем время начала свечи в формате time.Time
func (k *Candle) GeTime() time.Time {
	return time.Unix(k.Time, 0)
}

type History struct {
	Candles []Candle `json:"history"` // Данные по свечам
	Next    int64    `json:"next"`    // Время (UTC) начала следующей свечи
	Prev    int64    `json:"prev"`    // Время (UTC) начала предыдущей свечи
}

// GeNextTime вернем время начала следующей свечи в формате time.Time
func (k *History) GeNextTime() time.Time {
	return time.Unix(k.Next, 0)
}

// GePrevTime вернем время начала предыдущей свечи в формате time.Time
func (k *History) GePrevTime() time.Time {
	return time.Unix(k.Prev, 0)
}
