package alor

import (
	"time"
)

type User struct {
	Portfolio string `json:"portfolio"`
}

type Instrument struct {
	Symbol          string `json:"symbol"`
	Exchange        string `json:"exchange"`
	InstrumentGroup string `json:"instrumentGroup,omitempty"`
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
	LastPriceTimestamp   int64   `json:"last_price_timestamp"` //  Unix time seconds для значения поля last_price
	LotSize              float64 `json:"lotsize"`              // Размер лота
	LotValue             float64 `json:"lotvalue"`             // Суммарная стоимость лота
	FaceValue            float64 `json:"facevalue"`            // Показатель, значение которого варьируется в зависимости от выбранного рынка:
	OpenInterest         int64   `json:"open_interest"`        // Открытый интерес (open interest). Если не поддерживается инструментом — значение 0 или null
	OrderBookMSTimestamp int64   `json:"ob_ms_timestamp"`      // Временная метка (UTC) сообщения о состоянии биржевого стакана в формате Unix Time Milliseconds
	Type                 string  `json:"type"`                 // Полное название фьючерса
	//Change               float64 `json:"change"`               // Разность цены и цены предыдущего закрытия
	//ChangePercent        float64 `json:"change_percent"`       // Относительное изменение цены
	//AccruedInt           float64 `json:"accruedInt"`           // Начислено (НКД)
}

// FaceValue
// Для фондового рынка — номинальная стоимость единицы финансового инструмента
// Для срочного рынка — размер одного лота
// Для валютного рынка — количество валюты лота, за которое указывается цена в котировках

// переведем время с UTC-timestamp в Time
func (q Quote) LastTime() time.Time {
	return time.Unix(q.LastPriceTimestamp, 0)
}

//func (q Quote) LastTime2() time.Time {
//	return time.Unix(q.OrderBookMSTimestamp, 0)
//}

// Interval период свечей
type Interval string

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

var intervalName = map[Interval]string{
	Interval_S15: "S15",
	Interval_M1:  "M1",
	Interval_H1:  "H1",
	Interval_D1:  "D1",
	Interval_W1:  "W1",
	Interval_MN1: "MN1",
	Interval_Y1:  "Y1",
}

func (i Interval) String() string {
	return intervalName[i]
	//return string(i)
}

// Candle Параметры свечи
type Candle struct {
	Symbol   string   `json:"symbol"`   // Код финансового инструмента (Тикер)
	Interval Interval `json:"interval"` // Интервал свечи
	Time     int64    `json:"time"`     // Время (UTC) (Unix time seconds)
	Close    float64  `json:"close"`    // Цена при закрытии
	Open     float64  `json:"open"`     // Цена при открытии
	High     float64  `json:"high"`     // Максимальная цена
	Low      float64  `json:"low"`      // Минимальная цена
	Volume   int32    `json:"volume"`   // Объём
}

// GeTime вернем время начала свечи в формате time.Time
func (k *Candle) GeTime() time.Time {
	return time.Unix(k.Time, 0)
}

//func (k Candle) String() string {
//	str := fmt.Sprintf("%v,%v,%v, O:%v, H:%v, L:%v, C:%v, V:%v", k.GeTime().String(), k.Symbol, k.Interval, k.Open, k.High, k.Low, k.Close, k.Volume)
//	return str
//}

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

// направление сдели ( buy sell)
type SideType string

const (
	SideTypeBuy  SideType = "buy"
	SideTypeSell SideType = "sell"
)

// OrderType Тип заявки (limit market stop stopLimit)
type OrderType string

const (
	OrderTypeLimit     OrderType = "limit"     // Лимитная заявка
	OrderTypeMarket    OrderType = "market"    // Рыночная заявка
	OrderTypeStop      OrderType = "stop"      // Стоп-заявка
	OrderTypeStopLimit OrderType = "stopLimit" // Стоп-лимитная заявка
)

// OrderStatus статус заявки ( working filled canceled rejected)
type OrderStatus string

const (
	OrderStatusWorking  OrderStatus = "working"  // На исполнении
	OrderStatusFilled   OrderStatus = "filled"   // Полностъю исполнилась (выполнилась)
	OrderStatusCanceled OrderStatus = "canceled" // Отменена
	OrderStatusRejected OrderStatus = "rejected" // отклонена

)

// TimeInForce условие по времени действия заявки
type TimeInForce string

const (
	TimeInForceGTC    TimeInForce = "goodtillcancelled" // Активна до отмены
	TimeInForceDAY    TimeInForce = "oneday"            // До конца дня
	TimeInForceFOK    TimeInForce = "fillorkill"        // Исполнить целиком или отклонить
	TimeInForceCancel TimeInForce = "immediateorcancel" // Снять остаток
)

// ConditionType Условие срабатывания стоп/стоп-лимитной заявки
type ConditionType string

const (
	ConditionMore        ConditionType = "More"        // Цена срабатывания больше текущей цены
	ConditionLess        ConditionType = "Less"        // Цена срабатывания меньше текущей цены
	ConditionMoreOrEqual ConditionType = "MoreOrEqual" // Цена срабатывания больше или равна текущей цене
	ConditionLessOrEqual ConditionType = "LessOrEqual" // Цена срабатывания меньше или равна текущей цене
)
