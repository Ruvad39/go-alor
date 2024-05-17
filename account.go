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

// Portfolio информация о портфеле
type Portfolio struct {
	BuyingPowerAtMorning           float64 `json:"buyingPowerAtMorning"`           //Покупательская способность на утро
	BuyingPower                    float64 `json:"buyingPower"`                    // Покупательская способность
	Profit                         float64 `json:"profit"`                         // Прибыль за сегодня
	ProfitRate                     float64 `json:"profitRate"`                     // Норма прибыли, %
	PortfolioEvaluation            float64 `json:"portfolioEvaluation"`            // Ликвидный портфель
	PortfolioLiquidationValue      float64 `json:"portfolioLiquidationValue"`      // Оценка портфеля
	InitialMargin                  float64 `json:"initialMargin"`                  // Маржа
	RiskBeforeForcePositionClosing float64 `json:"riskBeforeForcePositionClosing"` // Риск до закрытия
	Commission                     float64 `json:"commission"`                     // Суммарная комиссия (null для Срочного рынка)
}

type Position struct {
	Portfolio         string  `json:"portfolio"`         // Идентификатор клиентского портфеля
	Symbol            string  `json:"symbol"`            // Тикер (Код финансового инструмента)
	BrokerSymbol      string  `json:"brokerSymbol"`      // Пара Биржа:Тикер
	Exchange          string  `json:"exchange"`          // Биржа
	ShortName         string  `json:"shortName"`         // Короткое наименование
	Volume            float64 `json:"volume"`            // Объём, рассчитанный по средней цен
	CurrentVolume     float64 `json:"currentVolume"`     // Объём, рассчитанный по текущей цене
	AvgPrice          float64 `json:"avgPrice"`          // Средняя цена
	QtyUnits          float64 `json:"qtyUnits"`          // Количество (штуки)
	OpenUnits         float64 `json:"openUnits"`         // Количество открытых позиций на момент открытия (начала торгов)
	LotSize           float64 `json:"lotSize"`           // Размер лота
	QtyT0             float64 `json:"qtyT0"`             // Агрегированное количество T0 (штуки)
	QtyT1             float64 `json:"qtyT1"`             // Агрегированное количество T1 (штуки)
	QtyT2             float64 `json:"qtyT2"`             // Агрегированное количество T2 (штуки)
	QtyTFuture        float64 `json:"qtyTFuture"`        // Количество (штуки)
	QtyT0Batch        float64 `json:"qtyT0Batch"`        // Агрегированное количество T0 (лоты)
	QtyT1Batch        float64 `json:"qtyT1Batch"`        // Агрегированное количество T1 (лоты)
	QtyT2Batch        float64 `json:"qtyT2Batch"`        // Агрегированное количество T2 (лоты)
	QtyTFutureBatch   float64 `json:"qtyTFutureBatch"`   // Агрегированное количество TFuture (лоты)
	QtyBatch          float64 `json:"qtyBatch"`          // Агрегированное количество TFuture
	OpenQtyBatch      float64 `json:"openQtyBatch"`      // Агрегированное количество на момент открытия (начала торгов) (лоты)
	Qty               float64 `json:"qty"`               // Агрегированное количество (лоты)
	Open              float64 `json:"open"`              // Агрегированное количество на момент открытия (начала торгов) (штуки)
	DailyUnrealisedPl float64 `json:"dailyUnrealisedPl"` // Суммарная прибыль или суммарный убыток за день в процентах
	UnrealisedPl      float64 `json:"unrealisedPl"`      // Суммарная прибыль или суммарный убыток за день в валюте расчётов
	IsCurrency        bool    `json:"isCurrency"`        // True для валютных остатков (денег), false - для торговых инструментов
}

// Lot вернем кол-во лот
func (p *Position) Lot() int64 {
	return int64(p.Qty)
}
