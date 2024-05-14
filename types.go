package alor

import (
	"fmt"
	"time"
)

// Security defines model for security.
type Security struct {
	Symbol                 string  `json:"symbol"`                  // Symbol Тикер (Код финансового инструмента)
	ShortName              string  `json:"shortname"`               // Shortname Краткое наименование инструмента
	Description            string  `json:"description,omitempty"`   // Description Краткое описание инструмента
	Exchange               string  `json:"exchange"`                // Exchange Биржа
	Board                  string  `json:"board"`                   //Код режима торгов (Борд):
	LotSize                float64 `json:"lotsize"`                 // Lotsize Размер лота
	MinStep                float64 `json:"minstep"`                 // Minstep Минимальный шаг цены
	PriceStep              float64 `json:"pricestep"`               // Pricestep Минимальный шаг цены, выраженный в рублях
	Cancellation           string  `json:"cancellation,omitempty"`  // Cancellation Дата и время (UTC) окончания действия
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
	TheorPrice             float64 `json:"theorPrice,omitempty"`
	TheorPriceLimit        float64 `json:"theorPriceLimit,omitempty"`
	TradingStatus          int     `json:"tradingStatus,omitempty"` // TradingStatus Торговый статус инструмента
	TradingStatusInfo      string  `json:"tradingStatusInfo"`       // TradingStatusInfo Описание торгового статуса инструмента
	Type                   string  `json:"type,omitempty"`          // Type Тип
	Volatility             float64 `json:"volatility,omitempty"`    // Volatility Волативность
	//Yield                  *string `json:"yield"`                   // может быть null
	//Yield                  *int    `json:"yield,omitempty"`

}

// Quotes
type Quote struct {
	Symbol             string  `json:"symbol"`
	Exchanges          string  `json:"exchange"`
	Description        string  `json:"description"`
	PrevClosePrice     float64 `json:"prev_close_price"` // Цена предыдущего закрытия
	LastPrice          float64 `json:"last_price"`       // PriceLast
	OpenPrice          float64 `json:"open_price"`       // PriceOpen
	HighPrice          float64 `json:"high_price"`       // PriceMaximum
	LowPrice           float64 `json:"low_price"`        // PriceMinimum
	Ask                float64 `json:"ask"`
	Bid                float64 `json:"bid"`
	AskVol             float32 `json:"ask_vol"`              // Количество лотов в ближайшем аске в биржевом стакане
	BidVol             float32 `json:"bid_vol"`              // Количество лотов в ближайшем биде в биржевом стакане
	AskVolumeTotal     int32   `json:"total_ask_vol"`        // Суммарное количество лотов во всех асках в биржевом стакане
	BidVolumeTotal     int32   `json:"total_bid_vol"`        // Суммарное количество лотов во всех бидах в биржевом стакане
	LastPriceTimestamp int64   `json:"last_price_timestamp"` // UTC-timestamp для значения поля last_price
	LotSize            float64 `json:"lotsize"`              // Размер лота
	LotValue           float64 `json:"lotvalue"`             // Суммарная стоимость лота
	FaceValue          float64 `json:"facevalue"`            // Показатель, значение которого варьируется в зависимости от выбранного рынка:
	// Для фондового рынка — номинальная стоимость единицы финансового инструмента
	// Для срочного рынка — размер одного лота
	//Д ля валютного рынка — количество валюты лота, за которое указывается цена в котировках
	OpenInterest         int64   `json:"open_interest"`   // Открытый интерес (open interest). Если не поддерживается инструментом — значение 0 или null
	AccruedInt           float64 `json:"accruedInt"`      // Начислено (НКД)
	OrderBookMSTimestamp int64   `json:"ob_ms_timestamp"` // Временная метка (UTC) сообщения о состоянии биржевого стакана в формате Unix Time Milliseconds
	Type                 string  `json:"type"`            // Полное название фьючерса
	Change               float64 `json:"change"`          // Разность цены и цены предыдущего закрытия
	ChangePercent        float64 `json:"change_percent"`  // Относительное изменение цены
}

// переведем время с UTC-timestamp в Time
func (q Quote) LastTime() time.Time {
	return time.Unix(q.LastPriceTimestamp, 0)
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

/*

positions[https://alor.dev/rawdocs2/WarpOpenAPIv2.yml#/components/schemas/positionposition

*/

// Interval период свечей
type Interval string

/*
15 — 15 секунд
60 — 60 секунд или 1 минута
3600 — 3600 секунд или 1 час
D — сутки (соответствует значению 86400)
W — неделя (соответствует значению 604800)
M — месяц (соответствует значению 2592000)
Y — год (соответствует значению 31536000)

// https://apidev.alor.ru/md/v2/history?symbol=SBER&exchange=MOEX&tf=D&from=1549000661&to=1634256000&format=Simple
*/

// Длительность таймфрейма. В качестве значения можно указать точное количество секунд или код таймфрейма
const (
	Interval_S15 Interval = "15"   // 15 секунд
	Interval_M1  Interval = "60"   // 60 секунд или 1 минута
	Interval_H1  Interval = "3600" // 3600 секунд или 1 час
	Interval_D1  Interval = "D"    // D — сутки (соответствует значению 86400)
	Interval_W1  Interval = "W"    // W — неделя (соответствует значению 604800)
	Interval_MN1 Interval = "M"    // M — месяц (соответствует значению 2592000)
	Interval_Y1  Interval = "Y"    // Y — год (соответствует значению 31536000)

)

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

// работаем над биржевым стаканом

type PriceVolume struct {
	Price  float64 `json:"price"`  // цена
	Volume int64   `json:"volume"` // объем
}

// PriceVolumeSlice Биды  Аски
type PriceVolumeSlice []PriceVolume

func (slice PriceVolumeSlice) Len() int { return len(slice) }

func (p PriceVolume) String() string {
	return fmt.Sprintf("PriceVolume{ Price: %s, Volume: %s }", p.Price, p.Volume)
}

// вернем второй элемент
func (slice PriceVolumeSlice) Second() (PriceVolume, bool) {
	if len(slice) > 1 {
		return slice[1], true
	}
	return PriceVolume{}, false
}

// вернем первый элемент
func (slice PriceVolumeSlice) First() (PriceVolume, bool) {
	if len(slice) > 0 {
		return slice[0], true
	}
	return PriceVolume{}, false
}

// вернем объем стакана
func (slice PriceVolumeSlice) SumDepth() int64 {
	var total int64
	for _, pv := range slice {
		total = total + pv.Volume
	}

	return total
}

func (slice PriceVolumeSlice) Copy() PriceVolumeSlice {
	var s = make(PriceVolumeSlice, len(slice))
	copy(s, slice)
	return s
}

// OrderBook биржевой стакан
type OrderBook struct {
	Bids        PriceVolumeSlice `json:"bids"`         // Биды
	Asks        PriceVolumeSlice `json:"asks"`         // Аски
	MsTimestamp int64            `json:"ms_timestamp"` // Время (UTC) в формате Unix Time Milliseconds
	Existing    bool             `json:"existing"`     // True - для данных из "снепшота", то есть из истории. False - для новых событий

}

func (b *OrderBook) LastTime() time.Time {
	return time.Unix(b.MsTimestamp, 0)
}

func (b *OrderBook) BestBid() (PriceVolume, bool) {
	if len(b.Bids) == 0 {
		return PriceVolume{}, false
	}

	return b.Bids[0], true
}

func (b *OrderBook) BestAsk() (PriceVolume, bool) {
	if len(b.Asks) == 0 {
		return PriceVolume{}, false
	}

	return b.Asks[0], true
}
