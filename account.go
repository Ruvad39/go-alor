package alor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

/*
/md/v2/Clients/{exchange}/{portfolio}/summary Получение информации о портфеле
GET "https://apidev.alor.ru/md/v2/Clients/MOEX/D39004/summary?format=Simple"

/md/v2/Clients/{exchange}/{portfolio}/positions Получение информации о позициях
withoutCurrency=true Исключить из ответа все денежные инструменты
https://apidev.alor.ru/md/v2/Clients/MOEX/D39004/positions?format=Simple&withoutCurrency=false
return positions

/md/v2/Clients/{exchange}/{portfolio}/positions/{symbol} Получение информации о позициях выбранного инструмента
https://apidev.alor.ru/md/v2/Clients/MOEX/D39004/positions/LKOH?format=Simple
return position

*/

// GetPortfolio Получение информации о портфеле
func (c *Client) GetPortfolio(ctx context.Context, portfolio string) (Portfolio, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "summary")
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := Portfolio{}
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

// GetPositions получение информации о позициях
func (c *Client) GetPositions(ctx context.Context, portfolio string) ([]Position, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "positions")
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := make([]Position, 0)
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

// Получение информации о позициях выбранного инструмента
func (c *Client) GetPosition(ctx context.Context, portfolio, symbol string) (Position, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "positions", symbol)
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := Position{}
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

// https://apidev.alor.ru/md/v2/Clients/P039004/positions?format=Simple
// GetLoginPositions Получение информации о позициях по логину
func (c *Client) GetLoginPositions(ctx context.Context, login string) ([]Position, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, login, "positions")
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := make([]Position, 0)
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
