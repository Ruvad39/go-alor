package alor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

// GetSecurity получить параметры по торговому инструменту
func (c *Client) GetSecurity(ctx context.Context, board, symbol string) (Security, error) {
	queryURL, _ := url.Parse("/md/v2/Securities")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange)
	queryURL.Path = path.Join(queryURL.Path, symbol)

	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}

	if board != "" {
		r.setParam("instrumentGroup", board)
	}

	result := Security{}
	data, err := c.callAPI(ctx, r)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil

}

// GetSecurities получить список торговых инструментов
// Объекты в ответе сортируются по объёму торгов.
// Если не указано иное значение параметра limit, в ответе возвращается только 25 объектов за раз
func (c *Client) GetSecurities(ctx context.Context, params Params) ([]Security, error) {
	r := &request{
		method:   http.MethodGet,
		endpoint: "/md/v2/Securities",
	}
	// если входящая биржа пустая, берем по умолчанию
	if params.Exchange == "" {
		params.Exchange = c.Exchange

	}
	r.setParam("exchange", params.Exchange)

	if params.Query != "" {
		r.setParam("query", params.Query)
	}

	if params.Board != "" {
		r.setParam("instrumentGroup", params.Board)
	}
	if params.Sector != "" {
		r.setParam("sector", params.Sector)
	}
	if params.Limit != 0 {
		r.setParam("limit", params.Limit)
	}
	if params.Offset != 0 {
		r.setParam("offset", params.Offset)
	}

	result := make([]Security, 0)
	data, err := c.callAPI(ctx, r)
	if err != nil {
		return result, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, err
	}
	return result, nil

}
