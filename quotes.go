package alor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"
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

// Quotes
type Quote struct {
	Symbol               string  `json:"symbol"`
	Exchanges            string  `json:"exchange"`
	Description          string  `json:"description"`
	PrevClosePrice       float64 `json:"prev_close_price"` // Цена предыдущего закрытия
	LastPrice            float64 `json:"last_price"`       // PriceLast
	OpenPrice            float64 `json:"open_price"`       // PriceOpen
	HighPrice            float64 `json:"high_price"`       // PriceMaximum
	LowPrice             float64 `json:"low_price"`        // PriceMinimum
	Ask                  float64 `json:"ask"`
	Bid                  float64 `json:"bid"`
	AskVol               float32 `json:"ask_vol"`              // Количество лотов в ближайшем аске в биржевом стакане
	BidVol               float32 `json:"bid_vol"`              // Количество лотов в ближайшем биде в биржевом стакане
	AskVolumeTotal       int32   `json:"total_ask_vol"`        // Суммарное количество лотов во всех асках в биржевом стакане
	BidVolumeTotal       int32   `json:"total_bid_vol"`        // Суммарное количество лотов во всех бидах в биржевом стакане
	LastPriceTimestamp   int64   `json:"last_price_timestamp"` // UTC-timestamp для значения поля last_price
	LotSize              float64 `json:"lotsize"`              // Размер лота
	LotValue             float64 `json:"lotvalue"`             // Суммарная стоимость лота
	FaceValue            float64 `json:"facevalue"`            // Показатель, значение которого варьируется в зависимости от выбранного рынка:
	OpenInterest         int64   `json:"open_interest"`        // Открытый интерес (open interest). Если не поддерживается инструментом — значение 0 или null
	AccruedInt           float64 `json:"accruedInt"`           // Начислено (НКД)
	OrderBookMSTimestamp int64   `json:"ob_ms_timestamp"`      // Временная метка (UTC) сообщения о состоянии биржевого стакана в формате Unix Time Milliseconds
	Type                 string  `json:"type"`                 // Полное название фьючерса
	Change               float64 `json:"change"`               // Разность цены и цены предыдущего закрытия
	ChangePercent        float64 `json:"change_percent"`       // Относительное изменение цены
}

// FaceValue
// Для фондового рынка — номинальная стоимость единицы финансового инструмента
// Для срочного рынка — размер одного лота
// Для валютного рынка — количество валюты лота, за которое указывается цена в котировках

// переведем время с UTC-timestamp в Time
func (q Quote) LastTime() time.Time {
	return time.Unix(q.LastPriceTimestamp, 0)
}
