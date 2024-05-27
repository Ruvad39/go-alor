package alor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

/*
	/md/v2/Securities/{symbols}/quotes Получение информации о котировках для выбранных инструментов

	symbols Принимает несколько пар биржа-тикер. Пары отделены запятыми. Биржа и тикер разделены двоеточием
			MOEX:SBER,MOEX:GAZP,SPBX:AAPL
*/

// GetQuotes Получение информации о котировках для выбранных инструментов.
// Принимает несколько пар биржа-тикер. Пары отделены запятыми. Биржа и тикер разделены двоеточием
func (c *Client) GetQuotes(ctx context.Context, symbols string) ([]Quote, error) {
	queryURL, _ := url.Parse("/md/v2/Securities")
	queryURL.Path = path.Join(queryURL.Path, symbols, "quotes")

	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}

	result := make([]Quote, 0)
	data, err := c.callAPI(ctx, r)
	if err != nil {
		//c.Logger.Error("GetQuotes ", "err", err.Error())
		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		//c.Logger.Error("GetQuotes ", "err", err.Error())
		return result, err
	}
	return result, nil

}

// GetQuote Получение информации о котировках для одного выбранного инструмента.
// Указываем тикер без указания биржи. Название биржи берется по умолчанию
func (c *Client) GetQuote(ctx context.Context, symbol string) (Quote, error) {
	ticker := c.Exchange + ":" + symbol
	quotes, err := c.GetQuotes(ctx, ticker)
	if err != nil {
		return Quote{}, err
	}
	if len(quotes) > 0 {
		return quotes[0], nil
	}
	return Quote{}, nil
}
