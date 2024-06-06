package alor

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
)

/*
/md/v2/Clients/{exchange}/{portfolio}/summary Получение информации о портфеле
/md/v2/Clients/{exchange}/{portfolio}/positions Получение информации о позициях
/md/v2/Clients/{exchange}/{portfolio}/positions/{symbol} Получение информации о позициях выбранного инструмента

TODO информации о сделках
/md/v2/Clients/{exchange}/{portfolio}/trades Получение информации о сделках  (только за текущую торговую сессию)
/md/v2/Clients/{exchange}/{portfolio}/{symbol}/trades

/md/v2/Stats/{exchange}/{portfolio}/history/trades Запрос списка сделок за предыдущие дни (не более 1000 сделок за один запрос)
https://apidev.alor.ru/md/v2/Stats/MOEX/D39004/history/trades?dateFrom=2024-04-10&ticker=SBER&from=93713183&limit=50&side=buy&format=Simple

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

// GetPortfolioRisk Получение информации по портфельным рискам
func (c *Client) GetPortfolioRisk(ctx context.Context, portfolio string) (PortfolioRisk, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "risk")
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := PortfolioRisk{}
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

// GetPortfolioFortsRisk Получение информации по портфельным рискам
func (c *Client) GetPortfolioFortsRisk(ctx context.Context, portfolio string) (PortfolioFortsRisk, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "fortsrisk")
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := PortfolioFortsRisk{}
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
// TODO использовать паттерн ok ([]Position, bool, error) ?
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
// TODO использовать паттерн ok (Position, bool, error) ?
// для этого нужно правильно обраьатывать ошибку
func (c *Client) GetPosition(ctx context.Context, portfolio, symbol string) (Position, bool, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "positions", symbol)
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := Position{}
	data, err := c.callAPI(ctx, r)
	if err != nil {
		// если ошибка "NotFound"
		if errors.Is(err, ErrNotFound) {
			return result, false, nil
		}
		return result, false, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, false, err
	}
	return result, true, nil
}

// https://apidev.alor.ru/md/v2/Clients/P039004/positions?format=Simple
// GetLoginPositions Получение информации о позициях по логину
// TODO использовать паттерн ok ([]Position, bool, error) ?
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

type PortfolioRisk struct {
	Portfolio                 string  `json:"portfolio"`                 // Идентификатор клиентского портфеля
	Exchange                  string  `json:"exchange"`                  // Биржа:
	PortfolioEvaluation       float64 `json:"portfolioEvaluation"`       // Общая стоимость портфеля
	PortfolioLiquidationValue float64 `json:"portfolioLiquidationValue"` // Стоимость ликвидного портфеля
	InitialMargin             float64 `json:"initialMargin"`             // Начальная маржа
	MinimalMargin             float64 `json:"minimalMargin"`             // Минимальная маржа
	CorrectedMargin           float64 `json:"correctedMargin"`           // Скорректированная маржа
	RiskCoverageRatioOne      float64 `json:"riskCoverageRatioOne"`      // НПР1
	RiskCoverageRatioTwo      float64 `json:"riskCoverageRatioTwo"`      // НПР2
	RiskCategoryId            int32   `json:"riskCategoryId"`            // Категория риска.
	ClientType                string  `json:"clientType"`                // Тип клиента:
	HasForbiddenPositions     bool    `json:"hasForbiddenPositions"`     // Имеются ли запретные позиции
	HasNegativeQuantity       bool    `json:"hasNegativeQuantity"`       // Имеются ли отрицательные количества
}

type PortfolioFortsRisk struct {
	Portfolio          string  `json:"portfolio"`          // Идентификатор клиентского портфеля
	MoneyFree          float64 `json:"moneyFree"`          // Свободные средства. Сумма рублей и залогов, дисконтированных в рубли, доступная для открытия позиций. (MoneyFree = MoneyAmount + VmInterCl – MoneyBlocked – VmReserve – Fee)
	MoneyBlocked       float64 `json:"moneyBlocked"`       // Средства, заблокированные под ГО
	Fee                float64 `json:"fee"`                // Списанный сбор
	MoneyOld           float64 `json:"moneyOld"`           // Общее количество рублей и дисконтированных в рубли залогов на начало сессии
	MoneyAmount        float64 `json:"moneyAmount"`        // Общее количество рублей и дисконтированных в рубли залогов
	MoneyPledgeAmount  float64 `json:"moneyPledgeAmount"`  // Сумма залогов, дисконтированных в рубли
	VmInterCl          float64 `json:"vmInterCl"`          // Вариационная маржа, списанная или полученная в пром. клиринг
	VmCurrentPositions float64 `json:"vmCurrentPositions"` // Сагрегированная вармаржа по текущим позициям
	VarMargin          float64 `json:"varMargin"`          // Вариационная маржа, рассчитанная по формуле VmCurrentPositions + VmInterCl
	IsLimitsSet        bool    `json:"isLimitsSet"`        // Наличие установленных денежного и залогового лимитов
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
