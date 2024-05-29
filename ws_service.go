package alor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

const (
	pingInterval            = 30 * time.Second
	reconnectCoolDownPeriod = 10 * time.Second
)

const (
	OnCandleSubscribe    = "BarsGetAndSubscribe"        // Подписка на историю цен (свечи)
	onQuotesSubscribe    = "QuotesSubscribe"            // Подписка на информацию о котировках
	onOrderBookSubscribe = "OrderBookGetAndSubscribe"   //  Подписка на биржевой стакан
	onAllTradesSubscribe = "AllTradesGetAndSubscribe"   // — Подписка на все сделки
	onPositionSubscribe  = "PositionsGetAndSubscribeV2" //  — Подписка на информацию о текущих позициях по торговым инструментам и деньгам
)

// IwsRequest Интерфейс которым должна обладать структура запроса для подписки
type IwsRequest interface {
	Marshal() ([]byte, error)
	GetOpCode() string
	GetGuid() string
	GetCode() string
	GetInterval() Interval
	SetToken(token string)
	SetExchange(exchange string)
}

type WsMessage struct {
	Message     string `json:"message"`
	HttpCode    int    `json:"httpCode"`
	RequestGuid string `json:"requestGuid"`
}

func (m *WsMessage) String() string {
	return fmt.Sprintf("<APIError> requestGuid=%s, httpCode=%d, message=%s", m.RequestGuid, m.HttpCode, m.Message)

}

type WSResponse struct {
	Data      *json.RawMessage `json:"data"` // Данные по ответу
	Guid      string           `json:"guid"` // Уникальный идентификатор запроса
	WsMessage                  // Системное сообщение
}

type WSRequestOption func(r *WSRequestBase)

// поля которые должны быть во всех запросах на подписку по websocket
type WSRequestBase struct {
	OpCode          string   `json:"opcode"`                    // Код операции
	Guid            string   `json:"guid"`                      // Уникальный идентификатор запроса. Все ответные сообщения будут иметь такое же значение поля guid
	Token           string   `json:"token"`                     // Access Токен для авторизации запроса
	Exchange        string   `json:"exchange"`                  // Биржа: MOEX — Московская Биржа SPBX — СПБ Биржа
	Frequency       int32    `json:"freq"`                      // Максимальная частота отдачи данных сервером в миллисекундах
	Format          string   `json:"format"`                    // Формат представления возвращаемых данных: Simple, Slim, Heavy
	Code            string   `json:"code,omitempty"`            // Код финансового инструмента (Тикер)
	Interval        Interval `json:"tf"`                        // Длительность таймфрейма в секундах или код (D — дни, W — недели, M — месяцы, Y — годы)
	From            int64    `json:"from,omitempty"`            // Дата и время (UTC) для первой запрашиваемой свечи
	SkipHistory     bool     `json:"skipHistory,omitempty"`     // Флаг отсеивания исторических данных: true — отображать только новые данные false — отображать в том числе данные из истории
	Depth           int32    `json:"depth,omitempty"`           // Глубина стакана. Стандартное и максимальное значение — 20 (20х20).
	Portfolio       string   `json:"portfolio,omitempty"`       // Идентификатор клиентского портфеля
	InstrumentGroup string   `json:"instrumentGroup,omitempty"` // Код режима торгов (Борд):
	OrderStatuses   string   `json:"orderStatuses,omitempty"`   // Опциональный фильтр по статусам заявок. Влияет только на фильтрацию первичных исторических данных при подписке. Возможные значения:
}

func (r *WSRequestBase) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (r *WSRequestBase) GetOpCode() string {
	return r.OpCode
}
func (r *WSRequestBase) GetGuid() string {
	return r.Guid
}
func (r *WSRequestBase) GetCode() string {
	return r.Code
}
func (r *WSRequestBase) GetInterval() Interval {
	return r.Interval
}
func (r *WSRequestBase) SetToken(token string) {
	r.Token = token
}
func (r *WSRequestBase) SetExchange(exchange string) {
	r.Exchange = exchange
}

func WithFrequency(param int32) WSRequestOption {
	return func(r *WSRequestBase) {
		r.Frequency = param
	}
}

// WSService сервис для подписок
type WsService struct {
	c          *Client       // Ссылка на основного клиента
	WsRequest  IwsRequest    // Структура запроса для подписки (wsRequest)
	CloseC     chan struct{} // CloseC сигнальный канал для закрытия коннекта
	ReconnectC chan struct{} // ReconnectC сигнальный канал для необходимости реконекта
	prevCandle Candle        // Предыдущая свеча (для работы с onCandleSubscribe)
}

func (c *Client) NewWsService(wsRequest IwsRequest) *WsService {
	s := &WsService{
		c:          c,
		WsRequest:  wsRequest,
		CloseC:     make(chan struct{}),
		ReconnectC: make(chan struct{}, 1),
	}
	s.WsRequest.SetExchange(c.Exchange)

	return s

}

// Close закроем сигнальный канал, что бы закончить работу
func (s *WsService) Close() {
	close(s.CloseC)
}

// Reconnect в сигнальный канал рекконета пошлем сообщение
func (s *WsService) Reconnect() {
	log.Debug("зашли в Reconnect()", "guid", s.WsRequest.GetGuid())
	select {
	case s.ReconnectC <- struct{}{}:
	default:
	}
}

// Connect запускаем поток и создаем соединение с websocket
func (s *WsService) Connect(ctx context.Context) error {
	err := s.DialAndConnect(ctx)
	if err != nil {
		return err
	}

	// запустим программу реконнекта start в отдельной горутине
	go s.reconnector(ctx)

	return nil
}

func (s *WsService) reconnector(ctx context.Context) {
	for {
		select {

		case <-ctx.Done():
			return

		case <-s.CloseC:
			return

		case <-s.ReconnectC:
			log.Warn("принят сигнал reconnect",
				"период восстановления повторного подключения", reconnectCoolDownPeriod,
				"guid", s.WsRequest.GetGuid(),
			)
			time.Sleep(reconnectCoolDownPeriod)

			log.Warn("re-connecting...")
			if err := s.DialAndConnect(ctx); err != nil {
				log.Error("re-connect error, try to reconnect later", "guid", s.WsRequest.GetGuid())
				// re-emit the re-connect signal if error
				s.Reconnect()
			}
		}
	}
}

// DialAndConnect создаем соединение с websocket. Запрашиваем подписку. Вызываем чтение данных
func (s *WsService) DialAndConnect(ctx context.Context) error {
	wssURL := getWsEndpoint()
	conn, resp, err := websocket.DefaultDialer.Dial(wssURL, nil)
	if err != nil {
		log.Error("websocket.Dial", "guid", s.WsRequest.GetGuid(), "err", err.Error())
		return err
	}
	// TODO создать и вернуть ошибку
	//if resp.StatusCode != 200 {
	if resp.StatusCode >= http.StatusBadRequest {
		log.Error("websocket.Dial", "guid", s.WsRequest.GetGuid(), "resp.StatusCode", resp.StatusCode)
	}
	// пошлем сообщение для подписки
	err = s.sendMessage(conn)
	if err != nil {
		return err
	}

	//connCtx, connCancel := context.WithCancel(ctx)
	connCtx, _ := context.WithCancel(ctx)

	//connCtx, connCancel := s.SetConn(ctx, conn)
	//s.EmitConnect()
	var wg sync.WaitGroup

	wg.Add(1)
	// запустим чтение данных с websocket
	go func() {
		defer wg.Done()
		s.Read(connCtx, conn)
	}()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	// запустим ping
	//	s.ping(connCtx, conn, connCancel)
	//}()

	wg.Wait()
	return nil
}
func (s *WsService) sendMessage(conn *websocket.Conn) error {
	log.Debug("зашли в sendMessage", "guid", s.WsRequest.GetGuid())
	// получим токен доступа
	token, err := s.c.GetJWT()
	if err != nil {
		log.Error("sendMessage: ошибка поучение токена доступа", "guid", s.WsRequest.GetGuid(), "err", err.Error())
		return err
	}
	s.WsRequest.SetToken(token)
	buf, err := s.WsRequest.Marshal()
	if err != nil {
		return err
	}
	//log.Debug("sendMessage", slog.Any("buf", buf))
	err = conn.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Error("sendMessage: conn.WriteMessage", "guid", s.WsRequest.GetGuid(), "err", err.Error())
		return err
	}
	log.Debug("sendMessage успешно послано", "guid", s.WsRequest.GetGuid())

	return nil
}

// Read Чтение данных с websocket
// func (s *StandardStream) Read(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc) {
func (s *WsService) Read(ctx context.Context, conn *websocket.Conn) {
	//log.Debug("зашли в Read()")
	//defer func() {
	//	cancel()
	//}()

	for {
		select {

		case <-ctx.Done():
			return
		case <-s.CloseC:
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("ReadMessage", "guid", s.WsRequest.GetGuid(), "err", err.Error())
				// и дальше обрабатываем разные типы ошибок
				_ = conn.Close()
				s.Reconnect()
				return
			}
			// пошлем в обработчик
			s.handler(message)
		}
	}

}

// TODO создать обработчик по всем событиям
// или создать буферизированный канал с данными json.RawMessage и посылать данные туда
// предварительно запустив "НУЖНЫЙ" обработчик данных, где читать этот канал
func (s *WsService) handler(message []byte) {
	log.Debug("handler", "guid", s.WsRequest.GetGuid(), "message", string(message))
	msg := new(WSResponse)
	err := json.Unmarshal(message, msg)
	if err != nil {
		log.Error("handlerEvent", "guid", s.WsRequest.GetGuid(), "error json.Unmarshal", err.Error())
		return
	}
	//log.Debug("handler", "msg", msg)
	// системное сообщение
	if msg.HttpCode != 0 {

		if msg.HttpCode >= 400 {
			log.Error("handlerEvent", "guid", s.WsRequest.GetGuid(), "err", msg.String())
			// закроем обработчик
			s.Close()

		}
		return
	}
	// иначе информационное сообщение
	// TODO посылаем данные в канал ?
	switch s.WsRequest.GetOpCode() {
	case OnCandleSubscribe:
		s.onCandle(msg.Data)
	case onQuotesSubscribe:
		s.onQuote(msg.Data)
	default:
		log.Error("WsService.handler", "guid", s.WsRequest.GetGuid(), "OpCode неизвеcтен", s.WsRequest.GetOpCode())
	}
}

// onCandle обработка handler получения свечей
func (s *WsService) onCandle(data *json.RawMessage) {
	candle := Candle{}
	err := json.Unmarshal(*data, &candle)
	if err != nil {
		log.Error("WsService.onCandle", "guid", s.WsRequest.GetGuid(), "json.Unmarshaljson err", err.Error())
		return
	}

	candle.Symbol = s.WsRequest.GetCode()
	candle.Interval = s.WsRequest.GetInterval()

	// первый запуск
	if s.prevCandle.Time == 0 {
		s.prevCandle = candle
	}
	// # Пришла обновленная версия текущего бара
	if candle.Time == s.prevCandle.Time {
		s.prevCandle = candle
	}
	// # Пришла новая свеча
	if candle.Time > s.prevCandle.Time {
		// новая свеча
		log.Debug("WsService OnCandle", "guid", s.WsRequest.GetGuid(), "time", s.prevCandle.GeTime(), "candle", s.prevCandle)
		s.c.PublishCandleClosed(s.prevCandle) // пошлем в рассылку
		s.prevCandle = candle

	}

}

func (s *WsService) onQuote(data *json.RawMessage) {
	quote := Quote{}
	err := json.Unmarshal(*data, &quote)
	if err != nil {
		log.Error("WsService.onQuote", "guid", s.WsRequest.GetGuid(), "json.Unmarshaljson err", err.Error())
		return
	}
	//log.Debug("onQuote", slog.Any("Quote", quote))
	s.c.PublishQuotes(quote) // пошлем в рассылку

}

// Do запуск сервиса
func (s *WsService) Do(ctx context.Context) error {
	go func() {

		err := s.Connect(ctx)
		if err != nil {
			log.Error("s.Connect(ctx)", "guid", s.WsRequest.GetGuid(), "err", err)
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-s.CloseC:
			return

		}

	}()
	return nil
}
