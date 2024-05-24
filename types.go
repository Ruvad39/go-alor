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

func (i Interval) String() string {
	return string(i)
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
