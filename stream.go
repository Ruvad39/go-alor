package alor

type CandleCloseCallback func(candle Candle)
type QuoteCallback func(quote Quote)

//type DataFeedConsumer func(Candle)

type Stream struct {
	OnCandle CandleCloseCallback // Функция обработки появления новой свечи
	OnQuote  QuoteCallback       // Функция обработки появления котировки
	//candleClosedCallbacks []func(candle Candle) // Список (callbacks) на закрытие свечи
	//quotesCallbacks       []func(quote Quote)   // Список (callbacks) на котировки
}

// SetOnCandle регистрирует функцию для вызова OnCandleClosed
func (s *Stream) SetOnCandle(f CandleCloseCallback) {
	s.OnCandle = f
}

// SetOnQuote регистрирует функцию для вызова OnQuote
func (s *Stream) SetOnQuote(f QuoteCallback) {
	s.OnQuote = f
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

	//for _, cb := range s.quotesCallbacks {
	//	cb(quote)
	//}
}
