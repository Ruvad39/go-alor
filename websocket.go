package alor

import (
	"github.com/gorilla/websocket"
	"time"
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WebsocketTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WebsocketKeepalive = true
)

/*
 "message": "Handled successfully",
  "httpCode": 200,
  "requestGuid": "c328fcf1-e495-408a-a0ed-e20f95d6b813"
*/

/*
opcode string

OrderBookGetAndSubscribe — Подписка на биржевой стакан
BarsGetAndSubscribe — Подписка на историю цен (свечи)
QuotesSubscribe — Подписка на информацию о котировках
InstrumentsGetAndSubscribeV2 — Подписка на изменение информации о финансовых инструментах на выбранной бирже
AllTradesGetAndSubscribe — Подписка на все сделки
PositionsGetAndSubscribeV2 — Подписка на информацию о текущих позициях по торговым инструментам и деньгам
SummariesGetAndSubscribeV2 — Подписка на сводную информацию по портфелю
RisksGetAndSubscribe — Подписка на сводную информацию по портфельным рискам
SpectraRisksGetAndSubscribe — Подписка на информацию по рискам срочного рынка (FORTS)
TradesGetAndSubscribeV2 — Подписка на информацию о сделках
OrdersGetAndSubscribeV2 — Подписка на информацию о текущих заявках на рынке для выбранных биржи и финансового инструмента
StopOrdersGetAndSubscribeV2 — Подписка на информацию о текущих заявках на рынке для выбранных биржи и финансового инструмента
Unsubscribe — Отмена существующей подписки

*/

type WsMessage struct {
	Message     string `json:"message"`
	HttpCode    int    `json:"httpCode"`
	RequestGuid string `json:"requestGuid"`
}

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// var wsServe = func(endPoint string, sendMessage []byte, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
// коннект
func wsServe(endPoint string, sendMessage []byte, handler WsHandler, errHandler ErrHandler) (doneC, stopC chan struct{}, err error) {
	log.Info("Start wsServe")
	conn, _, err := websocket.DefaultDialer.Dial(endPoint, nil)
	if err != nil {
		log.Error("Dial", "err", err.Error())
		return nil, nil, err
	}
	conn.SetReadLimit(655350)
	doneC = make(chan struct{})
	stopC = make(chan struct{})
	// пошлем сообщение
	err = conn.WriteMessage(websocket.TextMessage, sendMessage)
	if err != nil {
		log.Error("WriteMessage", "err", err.Error())
		return nil, nil, err
	}
	//log.Info("wsServe подписка прошла успешна", "guid", string(sendMessage))
	// в цикле читаем сообщения
	go func() {
		// This function will exit either on error from
		// websocket.Conn.ReadMessage or when the stopC channel is
		// closed by the client.
		defer close(doneC)
		if WebsocketKeepalive {
			keepAlive(conn, WebsocketTimeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			select {
			case <-stopC:
				silent = true
			case <-doneC:
			}
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
					//return
				}
				return
			}
			handler(message)
		}
	}()

	return
}

func keepAlive(conn *websocket.Conn, timeout time.Duration) {
	//log.Info("Start keepAlive")
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				conn.Close()
				return
			}
		}
	}()
}
