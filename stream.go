package alor

//type CandleCallback func(k Candle)

type Stream struct {
	//OnCandle CandleCallback  // Функция обработки появления новой свечи
	candleClosedCallbacks []func(candle Candle) // Список (callbacks) на закрытие свечи
	quotesCallbacks       []func(quote Quote)   // Список (callbacks) на котировки
}

// RegisterOnCandleClosed регистрирует функцию для вызова OnCandleClosed
func (s *Stream) RegisterOnCandleClosed(cb func(candle Candle)) {
	s.candleClosedCallbacks = append(s.candleClosedCallbacks, cb)
}

// PublishCandleClosed пошлем данные по свече дальше = тем кто подписался
func (s *Stream) PublishCandleClosed(candle Candle) {
	for _, f := range s.candleClosedCallbacks {
		f(candle)
	}
}

// RegisterOnQuotes регистрирует функцию на появление новой котировки (OnQuotes)
func (s *Stream) RegisterOnQuotes(cb func(quote Quote)) {
	s.quotesCallbacks = append(s.quotesCallbacks, cb)
}

// PublishQuotes пошлем котировки = тем кто подписался
func (s *Stream) PublishQuotes(quote Quote) {
	for _, cb := range s.quotesCallbacks {
		cb(quote)
	}
}
