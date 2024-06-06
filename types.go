package alor

import (
	"fmt"
	"strconv"
	"time"
)

var TzMsk = initMoscow()

func initMoscow() *time.Location {
	var loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", int(3*time.Hour/time.Second))
	}
	return loc
}

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

// переведем время с UTC-timestamp в Time и сразу поменяем в Московскеое время
func (q Quote) LastTime() time.Time {
	return time.Unix(q.LastPriceTimestamp, 0).In(TzMsk)
	//return time.Unix(q.LastPriceTimestamp, 0)
}

//func (q Quote) LastTime2() time.Time {
//	return time.Unix(q.OrderBookMSTimestamp, 0)
//}

// Interval период свечей
type Interval string

func (i Interval) String() string {
	return string(i)
}

func (i Interval) ToString() string {
	return intervalToName[i]
}

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

var intervalToName = map[Interval]string{
	Interval_S15: "S15",
	Interval_M1:  "M1",
	Interval_H1:  "H1",
	Interval_D1:  "D1",
	Interval_W1:  "W1",
	Interval_MN1: "MN1",
	Interval_Y1:  "Y1",
}

var nameToInterval = map[string]Interval{
	"S15": Interval_S15,
	"M1":  Interval_M1,
	"H1":  Interval_H1,
	"D1":  Interval_D1,
	"W1":  Interval_W1,
	"Y1":  Interval_Y1,
}

// ParseToInterval преобразуем символьную стоку в Interval
func ParseToInterval(input string) (Interval, error) {
	m, ok := nameToInterval[input]
	if !ok {
		return "", fmt.Errorf("не поддерживаемый формат периода свечи %s", input)
	}
	return m, nil
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
	//t := time.Unix(k.Time, 0).In(TzMsk)
	return time.Unix(k.Time, 0).In(TzMsk)

	//return time.Unix(k.Time, 0).LoadLocation(TzMsk)
}

// <TICKER>,<PER>,<DATE>,<TIME>,<OPEN>,<HIGH>,<LOW>,<CLOSE>,<VOL>

//func (k *Candle) CsvHeader() []string {
//return []string{
//	"<TICKER>", "<PER>", "<DATE>", "<TIME>", "<OPEN>", "<HIGH>", "<LOW>", "<CLOSE>", "<VOLUME>",
//}

func (k *Candle) CsvHeader() string {
	return "<TICKER>,<PER>,<DATE>,<TIME>,<OPEN>,<HIGH>,<LOW>,<CLOSE>,<VOLUME>"
}

// LKOH Splice,1,20130108,100700,20525,20525,20485,20504,138

// возвращает строку через запятую
func (k *Candle) CsvRecord() string {
	delimiter := ","
	return fmt.Sprint(
		k.Symbol, delimiter,
		k.Interval.ToString(), delimiter,
		k.GeTime().Format("20060102"), delimiter,
		k.GeTime().Format("150405"), delimiter,
		strconv.FormatFloat(k.Open, 'f', -1, 64), delimiter,
		strconv.FormatFloat(k.High, 'f', -1, 64), delimiter,
		strconv.FormatFloat(k.Low, 'f', -1, 64), delimiter,
		strconv.FormatFloat(k.Low, 'f', -1, 64), delimiter,
		strconv.FormatFloat(k.Close, 'f', -1, 64), delimiter,
		strconv.FormatInt(int64(k.Volume), 10),
	)

}

// возвращает массив строки для записи через "encoding/csv"
func (k *Candle) StringRecord() []string {
	return []string{
		k.Symbol,
		k.Interval.String(),
		k.GeTime().Format("20060102"),
		k.GeTime().Format("150405"),
		strconv.FormatFloat(k.Open, 'f', -1, 64),
		strconv.FormatFloat(k.High, 'f', -1, 64),
		strconv.FormatFloat(k.Low, 'f', -1, 64),
		strconv.FormatFloat(k.Low, 'f', -1, 64),
		strconv.FormatFloat(k.Close, 'f', -1, 64),
		strconv.FormatInt(int64(k.Volume), 10),
	}
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

type Order struct {
	ID             string      `json:"id"`             // Уникальный идентификатор заявки
	Symbol         string      `json:"symbol"`         // Тикер (Код финансового инструмента)
	BrokerSymbol   string      `json:"brokerSymbol"`   // Пара Биржа:Тикер
	Exchange       string      `json:"exchange"`       // Биржа
	Portfolio      string      `json:"portfolio"`      // Идентификатор клиентского портфеля
	Comment        string      `json:"comment"`        // Комментарий к заявке
	Type           OrderType   `json:"type"`           // Тип заявки limit - Лимитная заявка market - Рыночная заявка
	Side           SideType    `json:"side"`           // Направление сделки. buy — Купля sell — Продажа
	Status         OrderStatus `json:"status"`         // статус заявки
	TransitionTime string      `json:"transTime"`      // Дата и время выставления (UTC)
	UpdateTime     string      `json:"updateTime"`     // Дата и время изменения статуса заявки (UTC)
	EndTime        string      `json:"endTime"`        // Дата и время завершения (UTC)
	QtyUnits       int32       `json:"qtyUnits"`       // Количество (штуки)
	QtyBatch       int32       `json:"qtyBatch"`       // Количество (лоты)
	Qty            int32       `json:"qty"`            // Количество (лоты)
	FilledQtyUnits int32       `json:"filledQtyUnits"` // Количество исполненных (штуки)
	FilledQtyBatch int32       `json:"filledQtyBatch"` // Количество исполненных (лоты)
	Filled         int32       `json:"filled"`         // Количество исполненных (лоты)
	Price          float64     `json:"price"`          // Цена
	Existing       bool        `json:"existing"`       // True - для данных из "снепшота", то есть из истории. False - для новых событий
	TimeInForce    TimeInForce `json:"timeInForce"`    // Тип заявки oneday - До конца дня goodtillcancelled - Активна до отмены
	Volume         float64     `json:"volume"`         // Объем, для рыночных заявок - null
	//Iceberg // Специальные поля для сделок со скрытой частью
}

func IsActiveOrder(o Order) bool {
	return o.Status == OrderStatusWorking
}

// структура сделки
type Trade struct {
	Id           string    `json:"id"`           // Уникальный идентификатор сделки
	OrderNo      string    `json:"orderNo"`      // Уникальный идентификатор заявки
	Comment      string    `json:"comment"`      // Пользовательский комментарий к заявке
	Symbol       string    `json:"symbol"`       // Тикер (Код финансового инструмента).
	BrokerSymbol string    `json:"brokerSymbol"` // Пара Биржа:Тикер
	Exchange     string    `json:"exchange"`     // Биржа
	Date         time.Time `json:"date"`         // Дата и время завершения (UTC)
	Board        string    `json:"board"`        // Код режима торгов (Борд):
	QtyUnits     int32     `json:"qtyUnits"`     // Количество (штуки)
	QtyBatch     int       `json:"qtyBatch"`     // Количество (лоты)
	Qty          int       `json:"qty"`          // Количество (лоты)
	Price        float64   `json:"price"`        // Цена
	AccruedInt   int       `json:"accruedInt"`   // Начислено (НКД)
	Side         string    `json:"side"`         // Направление сделки:
	Existing     bool      `json:"existing"`     // True — для данных из "снепшота", то есть из истории. False — для новых событий
	Commission   float64   `json:"commission"`   // Суммарная комиссия (null для Срочного рынка)
	//RepoSpecificFields interface{} `json:"repoSpecificFields"` // Специальные поля для сделок РЕПО
	Volume float64 `json:"volume"` // Объём, рассчитанный по средней цене
}
