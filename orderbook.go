package alor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// /md/v2/orderbooks/{exchange}/{symbol} Получение информации о биржевом стакане
// https://apidev.alor.ru/md/v2/orderbooks/MOEX/LKOH?depth=20&format=Simple

// GetOrderBooks Получение информации о биржевом стакане
func (c *Client) GetOrderBooks(ctx context.Context, symbol string) (OrderBook, error) {
	queryURL, _ := url.Parse("/md/v2/orderbooks")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, symbol)

	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	// выставим максимальное
	r.setParam("depth", 20)

	result := OrderBook{}
	data, err := c.callAPI(ctx, r)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	return result, nil

}
