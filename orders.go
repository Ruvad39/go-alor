package alor

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"path"
)

// GetOrders получение информации о всех заявках
func (c *Client) GetOrders(ctx context.Context, portfolio string) ([]Order, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "orders")
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := make([]Order, 0)
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

// GetOrder получение информации о выбранной заявке
func (c *Client) GetOrder(ctx context.Context, portfolio, orderId string) (Order, error) {
	queryURL, _ := url.Parse("/md/v2/Clients")
	queryURL.Path = path.Join(queryURL.Path, c.Exchange, portfolio, "orders", orderId)
	r := &request{
		method:   http.MethodGet,
		endpoint: queryURL.String(),
	}
	result := Order{}
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

// снятие заявки
// ​/commandapi​/warptrans​/TRADE​/v2​/client​/orders​/{orderId}
// https://apidev.alor.ru/commandapi/warptrans/TRADE/v2/client/orders/93713183?portfolio=D39004&exchange=MOEX&stop=false&jsonResponse=true&format=Simple

//{
//"code": "OrderToCancelNotFound (404)",
//"message": "Order to cancel not found"
//}

// структура ответа
type OrderResponse struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	OrderNumber string `json:"orderNumber"`
}

// если ордер с таким номером не найден (или уже отменен) = error HTTP 400: Bad Request
// если 200 = то успешно
// как обрабатывать разные ситуации? нет ошибки = Ок

// CancelOrder отменить заявку
// TODO решить что возвращать
func (c *Client) CancelOrder(ctx context.Context, portfolio, orderId string) (bool, error) {
	queryURL, _ := url.Parse("/commandapi/warptrans/TRADE/v2/client/orders")
	queryURL.Path = path.Join(queryURL.Path, orderId)
	r := &request{
		method:   http.MethodDelete,
		endpoint: queryURL.String(),
	}
	r.setParam("exchange", c.Exchange)
	r.setParam("portfolio", portfolio)
	r.setParam("jsonResponse", "true")

	//result := OrderResponse{}
	data, err := c.callAPI(ctx, r)
	if err != nil {
		//return result, err
		return false, err
	}
	log.Debug("CancelOrder", slog.Any("response body", string(data)))
	//err = json.Unmarshal(data, &result)
	//if err != nil {
	//	return false, err
	//}
	return true, nil
}

// NewCreateOrderService создать новый ордер
func (c *Client) NewCreateOrderService() *CreateOrderService {
	user := User{
		Portfolio: c.Portfolio, // проставим по умолчанию
	}
	ticker := Instrument{
		Symbol:   "",
		Exchange: c.Exchange,
		//InstrumentGroup: board,
	}
	return &CreateOrderService{
		c: c,
		order: OrderRequest{
			//OrderType:   orderType,      // лимитная или рыночная
			//Side:        side,           // направление: покупка или продажа
			TimeInForce: TimeInForceDAY, // сразу проставим "До конца дня"
			Instrument:  ticker,
			User:        user,
			Quantity:    0,
			Price:       0,
		},
	}
}

// послать команду на создание нового ордера
// возвращает ID созданной заявки
func (s *CreateOrderService) Do(ctx context.Context) (string, error) {
	r := &request{
		method:   http.MethodPost,
		endpoint: "/commandapi/warptrans/TRADE/v2/client/orders/actions/market",
	}
	if s.order.OrderType == OrderTypeLimit {
		r.endpoint = "/commandapi/warptrans/TRADE/v2/client/orders/actions/limit"
	}

	// в request.body надо записать order
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(&s.order)
	r.body = buf

	// "X-ALOR-REQID"
	// Требуется уникальная случайная строка в качестве идентификатора запроса.
	// Если уже приходил запрос с таким идентификатором, то заявка не будет исполнена повторно,
	// а в качестве ответа будет возвращена копия ответа на первый запрос с таким значением идентификатора
	// Текущее время в наносекундах, прошедших с 01.01.1970 в UTC
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

// BuyMarket покупка по рынку
func (c *Client) BuyMarket(ctx context.Context, symbol string, lot int32, comment string) (string, error) {
	return c.NewCreateOrderService().Side(SideTypeBuy).OrderType(OrderTypeMarket).
		Symbol(symbol).Qty(lot).Comment(comment).Do(ctx)
}

// BuyLimit лимитная покупка
func (c *Client) BuyLimit(ctx context.Context, symbol string, lot int32, price float64, comment string) (string, error) {
	return c.NewCreateOrderService().Side(SideTypeBuy).OrderType(OrderTypeLimit).
		Symbol(symbol).Qty(lot).Price(price).Comment(comment).Do(ctx)
}

// SellMarket продажа по рынку
func (c *Client) SellMarket(ctx context.Context, symbol string, lot int32, comment string) (string, error) {
	return c.NewCreateOrderService().Side(SideTypeSell).OrderType(OrderTypeMarket).
		Symbol(symbol).Qty(lot).Comment(comment).Do(ctx)
}

// SellLimit лимитная продажа
func (c *Client) SellLimit(ctx context.Context, symbol string, lot int32, price float64, comment string) (string, error) {
	return c.NewCreateOrderService().Side(SideTypeSell).OrderType(OrderTypeLimit).
		Symbol(symbol).Qty(lot).Price(price).Comment(comment).Do(ctx)
}

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

type OrderRequest struct {
	OrderType       OrderType   `json:"-"`
	Side            SideType    `json:"side"`                      // Направление сделки: buy — Купля sell — Продажа
	Quantity        int32       `json:"quantity"`                  // Количество (лоты)
	Price           float64     `json:"price,omitempty"`           // Цена (только для лимитной)
	Comment         string      `json:"comment"`                   // Пользовательский комментарий к заявке
	Instrument      Instrument  `json:"instrument"`                // тикер
	User            User        `json:"user"`                      // данные portfolio
	TimeInForce     TimeInForce `json:"timeInForce"`               // Условие по времени действия заявки
	IcebergFixed    int32       `json:"icebergFixed,omitempty"`    // Видимая постоянная часть айсберг-заявки в лотах
	IcebergVariance float64     `json:"icebergVariance,omitempty"` // Амплитуда отклонения (в % от icebergFixed) случайной надбавки к видимой части айсберг-заявки. Только срочный рынок
}

// CreateOrderService создать новую заявку (ордер)
type CreateOrderService struct {
	c     *Client
	order OrderRequest
}

// Side установим направление ордера
func (s *CreateOrderService) Side(side SideType) *CreateOrderService {
	s.order.Side = side
	return s
}

// Comment установим комментарий
func (s *CreateOrderService) Comment(comment string) *CreateOrderService {
	s.order.Comment = comment
	return s
}

// Symbol установим символ
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.order.Instrument.Symbol = symbol
	return s
}

// Board установим Код режима торгов
func (s *CreateOrderService) Board(board string) *CreateOrderService {
	s.order.Instrument.InstrumentGroup = board
	return s
}

// Qty установим кол-во лот (Quantity)
func (s *CreateOrderService) Qty(quantity int32) *CreateOrderService {
	s.order.Quantity = quantity
	return s
}

// Price установить цену. Для лимитной заявки
func (s *CreateOrderService) Price(price float64) *CreateOrderService {
	s.order.Price = price
	return s
}

// TimeInForce установим Условие по времени действия заявки
func (s *CreateOrderService) TimeInForce(timeInForce TimeInForce) *CreateOrderService {
	s.order.TimeInForce = timeInForce
	return s
}

// Portfolio установим номер торгового счета
func (s *CreateOrderService) Portfolio(portfolio string) *CreateOrderService {
	s.order.User.Portfolio = portfolio
	return s
}

// OrderType установим тип заявки
func (s *CreateOrderService) OrderType(orderType OrderType) *CreateOrderService {
	s.order.OrderType = orderType
	return s
}
