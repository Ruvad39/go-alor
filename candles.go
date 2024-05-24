package alor

import (
	"context"
	"encoding/json"
	"net/http"
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
