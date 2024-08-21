package alor

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
)

// GetSecurity получить параметры по торговому инструменту
// TODO использовать паттерн ok (Security, bool, error) ?
func (c *Client) GetSecurity(ctx context.Context, board, symbol string) (Security, bool, error) {
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
		// если ошибка "NotFound"
		if errors.Is(err, ErrNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	if err = json.Unmarshal(data, &result); err != nil {
		return result, false, err
	}
	return result, true, nil

}

// GetSecurities получить список торговых инструментов
// Объекты в ответе сортируются по объёму торгов.
// Если не указано иное значение параметра limit, в ответе возвращается только 25 объектов за раз
func (c *Client) GetSecurities(ctx context.Context, opts ...Option) ([]Security, error) {
	params := &Options{}
	// Обработаем входящие параметры
	for _, opt := range opts {
		opt(params)
	}
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
	if params.IncludeOld {
		r.setParam("includeOld", params.IncludeOld)

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

// Security defines model for security.
type Security struct {
	Symbol       string  `json:"symbol"`                 // Symbol Тикер (Код финансового инструмента)
	ShortName    string  `json:"shortname"`              // Shortname Краткое наименование инструмента
	Description  string  `json:"description,omitempty"`  // Description Краткое описание инструмента
	Exchange     string  `json:"exchange"`               // Exchange Биржа
	Board        string  `json:"board"`                  //Код режима торгов (Борд):
	LotSize      float64 `json:"lotsize"`                // Lotsize Размер лота
	MinStep      float64 `json:"minstep"`                // Minstep Минимальный шаг цены
	PriceStep    float64 `json:"pricestep"`              // Pricestep Минимальный шаг цены, выраженный в рублях
	Cancellation string  `json:"cancellation,omitempty"` // Cancellation Дата и время (UTC) окончания действия
	//Cancellation           time.Time `json:"cancellation,omitempty"`  // Cancellation Дата и время (UTC) окончания действия
	CfiCode                string  `json:"cfiCode,omitempty"`       // CfiCode Тип ценной бумаги согласно стандарту ISO 10962
	ComplexProductCategory string  `json:"complexProductCategory"`  // ComplexProductCategory Требуемая категория для осуществления торговли инструментом
	Currency               string  `json:"currency,omitempty"`      // Currency Валюта
	Facevalue              float64 `json:"facevalue,omitempty"`     // Facevalue Номинальная стоимость
	Marginbuy              float64 `json:"marginbuy,omitempty"`     // Marginbuy Цена маржинальной покупки (заемные средства)
	Marginrate             float64 `json:"marginrate,omitempty"`    // Marginrate Отношение цены маржинальной покупки к цене последней сделки
	Marginsell             float64 `json:"marginsell,omitempty"`    // Marginsell Цена маржинальной продажи (заемные средства)
	PriceMax               float64 `json:"priceMax,omitempty"`      // PriceMax Максимальная цена
	PriceMin               float64 `json:"priceMin,omitempty"`      // PriceMin Минимальная цена
	PrimaryBoard           string  `json:"primary_board,omitempty"` // PrimaryBoard Код режима торгов
	Rating                 float64 `json:"rating,omitempty"`
	OptionSide             string  `json:"optionside,omitempty"`  // Только для опционов. Сторона опциона:
	StrikePrice            float64 `json:"strikePrice,omitempty"` // Только для опционов. Цена Страйк (Цена исполнения опциона)
	TheorPrice             float64 `json:"theorPrice,omitempty"`
	TheorPriceLimit        float64 `json:"theorPriceLimit,omitempty"`
	TradingStatus          int     `json:"tradingStatus,omitempty"` // TradingStatus Торговый статус инструмента
	TradingStatusInfo      string  `json:"tradingStatusInfo"`       // TradingStatusInfo Описание торгового статуса инструмента
	Type                   string  `json:"type,omitempty"`          // Type Тип
	Volatility             float64 `json:"volatility,omitempty"`    // Volatility Волативность
	//Yield                  *string `json:"yield"`                   // может быть null
	//Yield                  *int    `json:"yield,omitempty"`

}
