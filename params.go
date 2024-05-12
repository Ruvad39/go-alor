package alor

// Params параметры запроса по методам
type Params struct {
	Exchange string // Биржа MOEX, SPBX
	Sector   string // Рынок на бирже FORTS, FOND, CURR
	Board    string // Режим торгов (instrumentGroup)
	Symbol   string // Код инструмента
	Query    string // Query Тикер (Код финансового инструмента) ищет по вхождению
	//Format   string // Format Формат возвращаемого сервером JSON
	Limit  int32 // Ограничение на количество выдаваемых результатов поиска
	Offset int32 // Смещение начала выборки (для пагинации)
}
