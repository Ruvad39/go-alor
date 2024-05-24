package alor

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

// запрос на подписку свечей
type WSCandleRequest struct {
	Opcode      string   `json:"opcode"`      // Код операции
	Code        string   `json:"code"`        // Код финансового инструмента (Тикер)
	Interval    Interval `json:"tf"`          // Длительность таймфрейма в секундах или код (D — дни, W — недели, M — месяцы, Y — годы)
	From        int64    `json:"from"`        // Дата и время (UTC) для первой запрашиваемой свечи
	Delayed     bool     `json:"delayed"`     // Данные c задержкой в 15 минут. Для авторизованых клиентов задержка не применяется.
	SkipHistory bool     `json:"skipHistory"` // Флаг отсеивания исторических данных: true — отображать только новые данные false — отображать в том числе данные из истории
	Exchange    string   `json:"exchange"`    // Биржа: MOEX — Московская Биржа SPBX — СПБ Биржа
	Format      string   `json:"format"`      // Формат представления возвращаемых данных: Simple, Slim, Heavy
	Frequency   int32    `json:"freq"`        // Максимальная частота отдачи данных сервером в миллисекундах
	Guid        string   `json:"guid"`        // Уникальный идентификатор запроса. Все ответные сообщения будут иметь такое же значение поля guid
	Token       string   `json:"token"`       // Access Токен для авторизации запроса

}

type WSCandleResponse struct {
	Candle    Candle `json:"data"`
	Guid      string `json:"guid"`
	WsMessage        // Системное сообщение
}

// WSCandleService подписка на свечи
type WSCandleService struct {
	StandardStream
	c          *Client
	wsRequest  WSCandleRequest
	prevCandle Candle // Предыдущая свеча
}

func (c *Client) NewWSCandleService(symbol string, interval Interval) *WSCandleService {
	opCode := "BarsGetAndSubscribe"
	//_ = c.GetJWT() // получим токен доступа

	//
	stream := &WSCandleService{
		StandardStream: NewStandardStream(),
		c:              c,
		wsRequest: WSCandleRequest{
			Opcode:      opCode,
			Code:        symbol,
			Interval:    interval,
			SkipHistory: true,
			Exchange:    c.Exchange,
			Token:       c.accessToken,
			Frequency:   1000, //  По умолчанию обновление в 1сек (1000 мс)

		},
	}
	stream.SetHandler(stream.handlerMessage)
	stream.SetSender(stream.sendMessage)

	return stream
}

// Symbol установим символ
func (s *WSCandleService) Symbol(symbol string) *WSCandleService {
	s.wsRequest.Code = symbol
	return s
}

func (s *WSCandleService) Exchange(exchange string) *WSCandleService {
	s.wsRequest.Exchange = exchange
	return s
}

func (s *WSCandleService) Interval(interval Interval) *WSCandleService {
	s.wsRequest.Interval = interval
	return s
}

// From Дата и время (UTC) для первой запрашиваемой свечи
func (s *WSCandleService) From(from int64) *WSCandleService {
	s.wsRequest.From = from
	return s
}

// SkipHistory Флаг отсеивания исторических данных: true — отображать только новые данные false — отображать в том числе данные из истории
func (s *WSCandleService) SkipHistory(skip bool) *WSCandleService {
	s.wsRequest.SkipHistory = skip
	return s
}

// Delayed данные c задержкой в 15 минут. Для авторизованых клиентов задержка не применяется.
func (s *WSCandleService) Delayed(delayed bool) *WSCandleService {
	s.wsRequest.Delayed = delayed
	return s
}

// Frequency Максимальная частота отдачи данных сервером в миллисекундах (default: 175)
func (s *WSCandleService) Frequency(frequency int32) *WSCandleService {
	s.wsRequest.Frequency = frequency
	return s
}

func (s *WSCandleService) makeGuid() string {
	return "candle|" + s.wsRequest.Code + "|" + s.wsRequest.Interval.String()
}

func (s *WSCandleService) Do(ctx context.Context) error {
	s.wssURL = getWsEndpoint()
	//s.guid = s.makeGuid()
	s.wsRequest.Guid = s.makeGuid()

	// этот кусок кода должен быть в SendMessage
	//buf, err := json.Marshal(s.wsRequest)
	//if err != nil {
	//	return err
	//}
	//s.sendMessage = buf
	// вызов в отдельной горутине
	go func() {

		log.Info("s.Connect(ctx)")
		err := s.Connect(ctx)
		if err != nil {
			log.Error("s.Connect(ctx)", "err", err)
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

// создать подписку на свечиx через вызов wsServe (старый вариант)
func (s *WSCandleService) WsServe(ctx context.Context) error {
	endPoint := getWsEndpoint()
	_ = s.c.GetJWT()
	guid := s.makeGuid()
	s.wsRequest.Guid = guid

	buf, err := json.Marshal(s.wsRequest)
	if err != nil {
		return err
	}

	var handleError = func(err error) {
		log.Error("handleError", "err", err.Error())
	}
	// TODO как правильно оброботать ошибку авторизации ?  {"requestGuid":"candle|Si-6.24|60","httpCode":400,"message":"Request message should contain 'token' field!"}
	go func() {
		// вызов wsServe
		for {
			//done, _, _ := wsServe(endPoint, buf, handler, handleError)
			log.Info("вызвали wsServe")
			done, _, _ := wsServe(endPoint, buf, s.handlerMessage, handleError)

			//<-done		}
			//done, stopC, _ := wsServe(endPoint, buf, handler, handleError)
			select {
			case <-ctx.Done():
				return
			case <-s.CloseC:
				return
			case <-done:
				//return // если сделать return = закончим обработку
				//case <-stopC:

			}
		}
	}()
	return nil
}

// sendMessage пошлем сообщение для подписки
// этот метод придется дублировать во все сервисы с работой websocket
func (s *WSCandleService) sendMessage(conn *websocket.Conn) error {
	log.Info("зашли в sendMessage", "guid", s.wsRequest.Guid)

	// получим токен доступа
	err := s.c.GetJWT()

	if err != nil {
		log.Error("sendMessage: ошибка поучение токена доступа", "err", err.Error())
		return err
	}
	s.wsRequest.Token = s.c.accessToken
	buf, err := json.Marshal(s.wsRequest)
	if err != nil {
		return err
	}
	//log.Debug("sendMessage", slog.Any("buf", buf))
	err = conn.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Error("sendMessage: conn.WriteMessage", "err", err.Error())
		return err
	}
	log.Info("sendMessage успешно послано", "guid", s.wsRequest.Guid)

	return nil
}

// обработчик
// func (s *WSCandleService) handlerEvent(e interface{}) {
// приходит RawMessage
func (s *WSCandleService) handlerMessage(message []byte) {
	log.Debug("handlerEvent", "message", string(message))
	//d, ok := e.(WSCandleResponse)
	//if !ok {
	//	return
	//}
	d := new(WSCandleResponse)
	err := json.Unmarshal(message, d)
	if err != nil {
		log.Error("handlerEvent", "error json.Unmarshal", err.Error())
		return
	}

	//log.Debug("handler", "data.HttpCode", d.HttpCode)
	if d.HttpCode >= http.StatusBadRequest {
		log.Error("handler", "data.HttpCode", d.HttpCode)
		s.Close()
		return
	}

	d.Candle.Symbol = s.wsRequest.Code
	d.Candle.Interval = s.wsRequest.Interval
	// первый запуск
	if s.prevCandle.Time == 0 {
		s.prevCandle = d.Candle
	}
	if d.Candle.Time > s.prevCandle.Time {
		// новая свеча
		log.Debug("WSCandleService OnCloseCandle", "time", s.prevCandle.GeTime(), "candle", s.prevCandle)
		s.c.PublishCandleClosed(s.prevCandle) // пошлем в рассылку
		s.prevCandle = d.Candle

	}
}

// SubscribeCandles подписка на свечи
func (c *Client) SubscribeCandles(ctx context.Context, symbol string, interval Interval) error {
	return c.NewWSCandleService(symbol, interval).Do(ctx)
}
