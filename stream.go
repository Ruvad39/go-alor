package alor

type CandleCloseFunc func(candle Candle)
type QuoteFunc func(quote Quote)
type OrderFunc func(order Order)

//type DataFeedConsumer func(Candle)

type Stream struct {
	OnCandle CandleCloseFunc // Функция обработки появления новой свечи
	OnQuote  QuoteFunc       // Функция обработки появления котировки
	OnOrder  OrderFunc       // Функция обработки появления заявках
}

// SetOnCandle регистрирует функцию для вызова OnCandleClosed
func (s *Stream) SetOnCandle(f CandleCloseFunc) {
	s.OnCandle = f
}

// SetOnQuote регистрирует функцию для вызова OnQuote
func (s *Stream) SetOnQuote(f QuoteFunc) {
	s.OnQuote = f
}

// SetOnOrder регистрирует функцию для вызова OnOrder
func (s *Stream) SetOnOrder(f OrderFunc) {
	s.OnOrder = f
}

// RegisterOnCandleClosed регистрирует функцию для вызова OnCandleClosed
//func (s *Stream) RegisterOnCandleClosed(f func(candle Candle)) {
//	s.candleClosedCallbacks = append(s.candleClosedCallbacks, f)
//}

// PublishCandleClosed пошлем данные по свече дальше = тем кто подписался
func (s *Stream) PublishCandleClosed(candle Candle) {
	hasFunction := s.OnCandle != nil
	if !hasFunction {
		log.Error("PublishCandleClosed: не зарегистирована функция OnCandle")
		return
	}
	s.OnCandle(candle)

	//for _, f := range s.candleClosedCallbacks {
	//	f(candle)
	//}
}

// RegisterOnQuotes регистрирует функцию на появление новой котировки (OnQuotes)
//func (s *Stream) RegisterOnQuotes(cb func(quote Quote)) {
//	s.quotesCallbacks = append(s.quotesCallbacks, cb)
//}

// PublishQuotes пошлем котировки = тем кто подписался
func (s *Stream) PublishQuotes(quote Quote) {
	hasFunction := s.OnQuote != nil
	if !hasFunction {
		log.Error("PublishQuotes: не зарегистрирована функция OnQuote")
		return
	}
	s.OnQuote(quote)

}

// PublishOrder пошлем заявки тем кто подписался
func (s *Stream) PublishOrder(order Order) {
	hasFunction := s.OnOrder != nil
	if !hasFunction {
		log.Error("PublishOrder: не зарегистрирована функция OnOrder")
		return
	}
	s.OnOrder(order)

}
