package alor

// Options параметры запроса по методам
type Options struct {
	Exchange string // Биржа MOEX, SPBX
	Sector   string // Рынок на бирже FORTS, FOND, CURR
	Board    string // Режим торгов (instrumentGroup)
	Symbol   string // Код инструмента
	Query    string // Query Тикер (Код финансового инструмента) ищет по вхождению
	Limit    int32  // Ограничение на количество выдаваемых результатов поиска
	Offset   int32  // Смещение начала выборки (для пагинации)
	//Format   string // Format Формат возвращаемого сервером JSON
}

type Option func(p *Options)

func NewOptions() *Options {
	o := &Options{}
	return o
}

// WithSector Рынок на бирже FORTS, FOND, CURR
func WithSector(param string) Option {
	return func(opts *Options) {
		opts.Sector = param
	}
}

// WithBoard Режим торгов (instrumentGroup)
func WithBoard(param string) Option {
	return func(opts *Options) {
		opts.Board = param
	}
}

// WithSymbol Код инструмента
func WithSymbol(param string) Option {
	return func(opts *Options) {
		opts.Symbol = param
	}
}

// WithExchange Биржа MOEX, SPBX
func WithExchange(param string) Option {
	return func(opts *Options) {
		opts.Exchange = param
	}
}

// WithQuery Тикер (Код финансового инструмента) ищет по вхождению
func WithQuery(param string) Option {
	return func(opts *Options) {
		opts.Query = param
	}
}

// WithLimit Ограничение на количество выдаваемых результатов поиска
func WithLimit(param int32) Option {
	return func(opts *Options) {
		opts.Limit = param
	}
}

// WithOffset Смещение начала выборки (для пагинации)
func WithOffset(param int32) Option {
	return func(opts *Options) {
		opts.Offset = param
	}
}
