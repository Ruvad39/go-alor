package alor

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const pingInterval = 30 * time.Second
const reconnectCoolDownPeriod = 10 * time.Second

//type CandleCallback func(k Candle)

type Stream struct {
	//OnCandle              CandleCallback        // Функция обработки появления новой свечи
	candleClosedCallbacks []func(candle Candle) // Список (callbacks) на закрытие свечи
}

// OnCandleClosed добавим callback на появление новой свечи в список подписок
func (s *Stream) OnCandleClosed(cb func(candle Candle)) {
	s.candleClosedCallbacks = append(s.candleClosedCallbacks, cb)
}

// PublishCandleClosed пошлем данные по свече дальше = тем кто подписался
func (s *Stream) PublishCandleClosed(candle Candle) {
	for _, cb := range s.candleClosedCallbacks {
		cb(candle)
	}
}

type Parser func(message []byte) (interface{}, error)

// type Handler func(e interface{})
type Handler func(message []byte)

type Sender func(conn *websocket.Conn) error

type StandardStream struct {
	parser       Parser
	handler      Handler       // Функция обработчик сообщения от websocket
	sender       Sender        // Функция для отправки сообщения по websocket
	pingInterval time.Duration // С какой периодичностью делам пинг
	CloseC       chan struct{} // CloseC сигнальный канал для закрытия стрима
	ReconnectC   chan struct{} // ReconnectC сигнальный канал для необходимости реконекта
	wssURL       string        // url для соединения с websocket
	//sendMessage  []byte        // Сообщение которое нужно послать для подписки
	//guid         string        // guid подписки (для отладки)
	//sg           SyncGroup
}

func NewStandardStream() StandardStream {
	return StandardStream{
		CloseC:       make(chan struct{}),
		ReconnectC:   make(chan struct{}, 1),
		pingInterval: pingInterval,
	}
}

func (s *StandardStream) SetSender(sender Sender) {
	s.sender = sender
}
func (s *StandardStream) SetParser(parser Parser) {
	s.parser = parser
}

func (s *StandardStream) SetHandler(handler Handler) {
	s.handler = handler
}

func (s *StandardStream) SetPingInterval(interval time.Duration) {
	s.pingInterval = interval
}

func (s *StandardStream) Close() {
	// закроем сигнальный канал, что бы закончить работу
	close(s.CloseC)
}

// Reconnect в сигнальный канал рекконета пошлем сообщение
func (s *StandardStream) Reconnect() {
	log.Debug("зашли в Reconnect()")
	select {
	case s.ReconnectC <- struct{}{}:
	default:
	}
}

// Connect запускаем поток и создаем соединение с websocket
func (s *StandardStream) Connect(ctx context.Context) error {
	err := s.DialAndConnect(ctx)
	if err != nil {
		return err
	}

	// запустим программу реконнекта start в отдельной горутине
	go s.reconnector(ctx)

	return nil
}

func (s *StandardStream) reconnector(ctx context.Context) {
	for {
		select {

		case <-ctx.Done():
			return

		case <-s.CloseC:
			return

		case <-s.ReconnectC:
			log.Warn("принят сигнал reconnect", "период восстановления повторного подключения", reconnectCoolDownPeriod)
			time.Sleep(reconnectCoolDownPeriod)

			log.Warn("re-connecting...")
			if err := s.DialAndConnect(ctx); err != nil {
				log.Error("re-connect error, try to reconnect later")
				// re-emit the re-connect signal if error
				s.Reconnect()
			}
		}
	}
}

// DialAndConnect создаем соединение с websocket. Запрашиваем подписку. Вызываем чтение данных
func (s *StandardStream) DialAndConnect(ctx context.Context) error {
	conn, _, err := websocket.DefaultDialer.Dial(s.wssURL, nil)
	if err != nil {
		log.Error("DefaultDialer.Dial", "err", err.Error())
		return err
	}
	// пошлем сообщение для подписки
	hasSender := s.sender != nil
	if hasSender {
		err := s.sender(conn)
		if err != nil {
			log.Error("WriteMessage", "err", err.Error())
			return err
		}

	}
	//log.Debug("WriteMessage", slog.Any("s.sendMessage", s.sendMessage))
	//err = conn.WriteMessage(websocket.TextMessage, s.sendMessage)
	//if err != nil {
	//	log.Error("WriteMessage", "err", err.Error())
	//	return err
	//}
	//log.Info("WriteMessage успешно послано", "guid", s.guid)

	//connCtx, connCancel := context.WithCancel(ctx)
	connCtx, _ := context.WithCancel(ctx)

	//connCtx, connCancel := s.SetConn(ctx, conn)
	//s.EmitConnect()
	var wg sync.WaitGroup

	wg.Add(1)
	// запустим чтение данных с websocket
	go func() {
		defer wg.Done()
		//s.Read(connCtx, conn, connCancel)
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

// Read Чтение данных с websocket
// func (s *StandardStream) Read(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc) {
func (s *StandardStream) Read(ctx context.Context, conn *websocket.Conn) {
	//log.Debug("зашли в Read()")
	//defer func() {
	//	cancel()
	//}()

	hasHandler := s.handler != nil
	for {
		select {

		case <-ctx.Done():
			return
		case <-s.CloseC:
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("ReadMessage", "err", err.Error())
				// и дальше обрабатываем разные типы ошибок
				s.Reconnect()
				return
			}
			if !hasHandler {
				log.Error("ReadMessage: Не установлен обработчик сообщений")
				continue
			}
			s.handler(message)
		}
	}

}
