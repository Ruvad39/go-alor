package alor

import (
    "time"
)


type Securities = []Security

// структура финансового инструмента
type Security struct {
    // Symbol Тикер (Код финансового инструмента)
    Symbol      string  `json:"symbol"`     
    
    // Shortname Краткое наименование инструмента
    ShortName   string `json:"shortname"`     

    // Description Краткое описание инструмента
    Description string `json:"description,omitempty"`

    // Exchange Биржа
    Exchange string `json:"exchange"` 

    //Код режима торгов (Борд):
    Board string  `json:"board"`     
    
    // Lotsize Размер лота
    LotSize float64 `json:"lotsize"`

    // Minstep Минимальный шаг цены
    MinStep float64 `json:"minstep"`

    // Pricestep Минимальный шаг цены, выраженный в рублях
    PriceStep float64 `json:"pricestep"`

    // Cancellation Дата и время (UTC) окончания действия
    Cancellation string `json:"cancellation,omitempty"`

    // CfiCode Тип ценной бумаги согласно стандарту ISO 10962
    CfiCode string `json:"cfiCode,omitempty"`

    // ComplexProductCategory Требуемая категория для осуществления торговли инструментом
    ComplexProductCategory string `json:"complexProductCategory"`

    // Currency Валюта
    Currency string `json:"currency,omitempty"`

    // Facevalue Номинальная стоимость
    Facevalue float64 `json:"facevalue,omitempty"`

    // Marginbuy Цена маржинальной покупки (заемные средства)
    Marginbuy float64 `json:"marginbuy,omitempty"`

    // Marginrate Отношение цены маржинальной покупки к цене последней сделки
    Marginrate float64 `json:"marginrate,omitempty"`

    // Marginsell Цена маржинальной продажи (заемные средства)
    Marginsell float64 `json:"marginsell,omitempty"`


    // PriceMax Максимальная цена
    PriceMax float64 `json:"priceMax,omitempty"`

    // PriceMin Минимальная цена
    PriceMin float64 `json:"priceMin,omitempty"`

    // PrimaryBoard Код режима торгов
    PrimaryBoard string  `json:"primary_board,omitempty"`
    Rating       float64 `json:"rating,omitempty"`

    TheorPrice      float64 `json:"theorPrice,omitempty"`
    TheorPriceLimit float64 `json:"theorPriceLimit,omitempty"`

    // TradingStatus Торговый статус инструмента
    TradingStatus int `json:"tradingStatus,omitempty"`

    // TradingStatusInfo Описание торгового статуса инструмента
    TradingStatusInfo string `json:"tradingStatusInfo"`

    // Type Тип
    Type string `json:"type,omitempty"`

    // Volatility Волативность
    Volatility float64 `json:"volatility,omitempty"`
    Yield      string  `json:"yield"`
}
/*
exchange
MOEX - Московская биржа
SPBX - СПБ Биржа

TradingStatus
Торговый статус инструмента:

18 - Нет торгов / торги закрыты
118 - Период открытия
103 - Период закрытия
2 - Перерыв в торгах
17 - Нормальный период торгов
102 - Аукцион закрытия
106 - Аукцион крупных пакетов
107 - Дискретный аукцион
119 - Аукцион открытия
120 - Период торгов по цене аукциона закрытия

type
Тип финансового инструмента.

Возможные значения для MOEX:

FOR — Валюта
CS — Обыкновенные акции компании
PS — Привилегированные акции компании
MF — Паевой инвестиционный фонд
RDR — Российская депозитарная расписка
EUSOV — Облигация федерального займа
MUNI — Муниципальная облигация
CORP — Корпоративная облигация
"Фьючерсный контракт X" — Фьючерсный контракт с указанием базового актива
"Марж. амер. Call X" — Опцион с указанием основных параметров контракта

tradingStatus   TradingStatusinteger($int32)
example: 17
Торговый статус инструмента:

18 - Нет торгов / торги закрыты
118 - Период открытия
103 - Период закрытия
2 - Перерыв в торгах
17 - Нормальный период торгов
102 - Аукцион закрытия
106 - Аукцион крупных пакетов
107 - Дискретный аукцион
119 - Аукцион открытия
120 - Период торгов по цене аукциона закрытия

*/

// Quotes
type Quote struct {
    Symbol              string  `json:"symbol"`
    Exchanges           string  `json:"exchange"`
    Description         string  `json:"description"`
    PrevClosePrice      float64  `json:"prev_close_price"` //Цена предыдущего закрытия
    LastPrice           float64  `json:"last_price"` // PriceLast
    OpenPrice           float64  `json:"open_price"` //  PriceOpen
    HighPrice           float64  `json:"high_price"` // PriceMaximum
    LowPrice            float64  `json:"low_price"` // PriceMinimum
    Ask                 float64  `json:"ask"`
    Bid                 float64  `json:"bid"`
    AskVol              float32  `json:"ask_vol"` //Количество лотов в ближайшем аске в биржевом стакане
    BidVol              float32  `json:"bid_vol"`  //Количество лотов в ближайшем биде в биржевом стакане
    AskVolumeTotal      int32    `json:"total_ask_vol"` //Суммарное количество лотов во всех асках в биржевом стакане
    BidVolumeTotal      int32    `json:"total_bid_vol"` //Суммарное количество лотов во всех бидах в биржевом стакане
    LastPriceTimestamp  int64    `json:"last_price_timestamp"` //UTC-timestamp для значения поля last_price
    LotSize             float64  `json:"lotsize"` //Размер лота
    LotValue            float64  `json:"lotvalue"` //Суммарная стоимость лота
    FaceValue           float64  `json:"facevalue"` //Показатель, значение которого варьируется в зависимости от выбранного рынка:
                                        // Для фондового рынка — номинальная стоимость единицы финансового инструмента
                                        // Для срочного рынка — размер одного лота
                                        //Д ля валютного рынка — количество валюты лота, за которое указывается цена в котировках
    OpenInterest        int64  `json:"open_interest"` //Открытый интерес (open interest). Если не поддерживается инструментом — значение 0 или null
    AccruedInt          float64  `json:"accruedInt"` // Начислено (НКД)
    Type                string   `json:"type"` //Полное название фьючерса
    Change              float64  `json:"change"` // Разность цены и цены предыдущего закрытия
    ChangePercent       float64  `json:"change_percent"` // Относительное изменение цены
    OrderBookMSTimestamp  int64  `json:"ob_ms_timestamp"` //Временная метка (UTC) сообщения о состоянии биржевого стакана в формате Unix Time Milliseconds
}

// переведем время с UTC-timestamp в Time
func (q Quote) LastTime() time.Time{
    return time.Unix(q.LastPriceTimestamp, 0) 
}

