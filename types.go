package alor

import "time"

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
	Yield                  string  `json:"yield"`
}

// Quotes
type Quote struct {
	Symbol             string  `json:"symbol"`
	Exchanges          string  `json:"exchange"`
	Description        string  `json:"description"`
	PrevClosePrice     float64 `json:"prev_close_price"` //Цена предыдущего закрытия
	LastPrice          float64 `json:"last_price"`       // PriceLast
	OpenPrice          float64 `json:"open_price"`       //  PriceOpen
	HighPrice          float64 `json:"high_price"`       // PriceMaximum
	LowPrice           float64 `json:"low_price"`        // PriceMinimum
	Ask                float64 `json:"ask"`
	Bid                float64 `json:"bid"`
	AskVol             float32 `json:"ask_vol"`              //Количество лотов в ближайшем аске в биржевом стакане
	BidVol             float32 `json:"bid_vol"`              //Количество лотов в ближайшем биде в биржевом стакане
	AskVolumeTotal     int32   `json:"total_ask_vol"`        //Суммарное количество лотов во всех асках в биржевом стакане
	BidVolumeTotal     int32   `json:"total_bid_vol"`        //Суммарное количество лотов во всех бидах в биржевом стакане
	LastPriceTimestamp int64   `json:"last_price_timestamp"` //UTC-timestamp для значения поля last_price
	LotSize            float64 `json:"lotsize"`              //Размер лота
	LotValue           float64 `json:"lotvalue"`             //Суммарная стоимость лота
	FaceValue          float64 `json:"facevalue"`            //Показатель, значение которого варьируется в зависимости от выбранного рынка:
	// Для фондового рынка — номинальная стоимость единицы финансового инструмента
	// Для срочного рынка — размер одного лота
	//Д ля валютного рынка — количество валюты лота, за которое указывается цена в котировках
	OpenInterest         int64   `json:"open_interest"`   //Открытый интерес (open interest). Если не поддерживается инструментом — значение 0 или null
	AccruedInt           float64 `json:"accruedInt"`      // Начислено (НКД)
	OrderBookMSTimestamp int64   `json:"ob_ms_timestamp"` //Временная метка (UTC) сообщения о состоянии биржевого стакана в формате Unix Time Milliseconds
	Type                 string  `json:"type"`            //Полное название фьючерса
	Change               float64 `json:"change"`          // Разность цены и цены предыдущего закрытия
	ChangePercent        float64 `json:"change_percent"`  // Относительное изменение цены
}

// переведем время с UTC-timestamp в Time
func (q Quote) LastTime() time.Time {
	return time.Unix(q.LastPriceTimestamp, 0)
}
