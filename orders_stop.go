package alor

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

/*
	/commandapi/warptrans/TRADE/v2/client/orders/actions/stopLimit
{
  "side": "buy",
  "condition": "More",
  "triggerPrice": 191.33,
  "stopEndUnixTime": 1590094740,
  "price": 191.33,
  "quantity": 1,
  "instrument": {
    "symbol": "SBER",
    "exchange": "MOEX",
    "instrumentGroup": "TQBR"
  },
  "user": {
    "portfolio": "D39004"
  },
  "timeInForce": "oneday",
  "icebergFixed": 100,
  "icebergVariance": 2,
  "activate": true
}


	/commandapi/warptrans/TRADE/v2/client/orders/actions/stop

{
  "side": "buy",
  "condition": "More",
  "triggerPrice": 191.33,
  "stopEndUnixTime": 1590094740,
  "quantity": 1,
  "instrument": {
    "symbol": "SBER",
    "exchange": "MOEX",
    "instrumentGroup": "TQBR"
  },
  "user": {
    "portfolio": "D39004"
  },
  "activate": true
}

*/

// OrderStopRequest запрос на создание stop/stopLimit заявки
type OrderStopRequest struct {
	OrderType       OrderType     `json:"-"`
	Condition       ConditionType `json:"condition"`                 // Условие срабатывания стоп/стоп-лимитной заявки:
	Side            SideType      `json:"side"`                      // Направление сделки: buy — Купля sell — Продажа
	Quantity        int32         `json:"quantity"`                  // Количество (лоты)
	TriggerPrice    float64       `json:"triggerPrice"`              // Стоп-цена
	Price           float64       `json:"price,omitempty"`           // Цена выставления стоп-лимитной заявки
	StopEndUnixTime int64         `json:"stopEndUnixTime"`           // Срок действия (UTC) в формате Unix Time seconds
	Instrument      Instrument    `json:"instrument"`                // тикер
	User            User          `json:"user"`                      // данные portfolio
	TimeInForce     TimeInForce   `json:"timeInForce"`               // Условие по времени действия заявки
	IcebergFixed    int32         `json:"icebergFixed,omitempty"`    // Видимая постоянная часть айсберг-заявки в лотах
	IcebergVariance float64       `json:"icebergVariance,omitempty"` // Амплитуда отклонения (в % от icebergFixed) случайной надбавки к видимой части айсберг-заявки. Только срочный рынок
	Activate        bool          `json:"activate"`                  // Флаг указывает, создать активную заявку, или не активную. Не активная заявка отображается в системе, но не участвует в процессе выставления на биржу, пока не станет активной. Данный флаг необходим при создании группы заявок с типом TriggerBracketOrders
}

// CreateOrderStopService создать новую stop/stopLimit заявку
type CreateOrderStopService struct {
	c     *Client
	order OrderStopRequest
}

// OrderType установим тип заявки
func (s *CreateOrderStopService) OrderType(orderType OrderType) *CreateOrderStopService {
	s.order.OrderType = orderType
	return s
}

// Condition Условие срабатывания стоп/стоп-лимитной заявки
func (s *CreateOrderStopService) Condition(condition ConditionType) *CreateOrderStopService {
	s.order.Condition = condition
	return s
}

// Side установим направление ордера
func (s *CreateOrderStopService) Side(side SideType) *CreateOrderStopService {
	s.order.Side = side
	return s
}

// Symbol установим символ
func (s *CreateOrderStopService) Symbol(symbol string) *CreateOrderStopService {
	s.order.Instrument.Symbol = symbol
	return s
}

// Board установим Код режима торгов
func (s *CreateOrderStopService) Board(board string) *CreateOrderStopService {
	s.order.Instrument.InstrumentGroup = board
	return s
}

// Qty установим кол-во лот (Quantity)
func (s *CreateOrderStopService) Qty(quantity int32) *CreateOrderStopService {
	s.order.Quantity = quantity
	return s
}

// Price установить цену. Для лимитной заявки
func (s *CreateOrderStopService) Price(price float64) *CreateOrderStopService {
	s.order.Price = price
	return s
}

// TriggerPrice Стоп-цена
func (s *CreateOrderStopService) TriggerPrice(price float64) *CreateOrderStopService {
	s.order.TriggerPrice = price
	return s
}

// TimeInForce установим Условие по времени действия заявки
func (s *CreateOrderStopService) TimeInForce(timeInForce TimeInForce) *CreateOrderStopService {
	s.order.TimeInForce = timeInForce
	return s
}

// Portfolio установим номер торгового счета
func (s *CreateOrderStopService) Portfolio(portfolio string) *CreateOrderStopService {
	s.order.User.Portfolio = portfolio
	return s
}

func (s *CreateOrderStopService) Activate(activate bool) *CreateOrderStopService {
	s.order.Activate = activate
	return s
}

// NewCreateOrderStopService создать новую stop/stopLimit заявку
func (c *Client) NewCreateOrderStopService() *CreateOrderStopService {
	user := User{
		Portfolio: c.Portfolio, // проставим по умолчанию
	}
	ticker := Instrument{
		Symbol:   "",
		Exchange: c.Exchange,
		//InstrumentGroup: board,
	}
	return &CreateOrderStopService{
		c: c,
		order: OrderStopRequest{
			OrderType:   OrderTypeStop,  // stop or stopLimit
			TimeInForce: TimeInForceDAY, // сразу проставим "До конца дня"
			Instrument:  ticker,
			User:        user,
			Quantity:    0,
			Price:       0,
			Activate:    true,
		},
	}
}

// Do послать команду на создание новой stop/stopLimit
// возвращает ID созданной заявки
func (s *CreateOrderStopService) Do(ctx context.Context) (string, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/commandapi/warptrans/TRADE/v2/client/orders/actions/stopLimit",
	}
	if s.order.OrderType == OrderTypeStop {
		r.endpoint = "/commandapi/warptrans/TRADE/v2/client/orders/actions/stop"
	}

	// в request.body надо записать order
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(&s.order)
	r.body = buf

	// Требуется уникальная случайная строка в качестве идентификатора запроса.
	r.setHeader("X-ALOR-REQID", s.c.getRequestID())
	// установим заголовок json
	r.setHeader("Content-Type", "application/json")

	result := OrderResponse{}
	data, err := s.c.callAPI(ctx, r)
	if err != nil {
		return "", err
	}

	log.Debug("SendOrder", slog.Any("response body", string(data)))
	err = json.Unmarshal(data, &result)
	if err != nil {
		return "", err
	}

	return result.OrderNumber, nil
}
