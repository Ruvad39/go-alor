package alor

/*
(1) Все сделки за сессию
https://apidev.alor.ru/md/v2/Clients/:exchange/:portfolio/trades
Запрос возвращает информацию обо всех сделках по указанному финансовому инструменту за текущую торговую сессию.

(2) Сделки по выбранному инструменту
Запрос возвращает информацию обо всех сделках по указанному финансовому инструменту.
https://apidev.alor.ru/md/v2/Clients/:exchange/:portfolio/:symbol/trade
QUERY PARAMETERS

(3) История сделок по портфелю
Запрос возвращает историческую информацию о сделках, совершённых с участием указанного профиля,
но не более 1000 записей за один запрос.
https://apidev.alor.ru/md/v2/Stats/:exchange/:portfolio/history/trades

QUERY PARAMETERS
dateFrom date Начиная с какой даты отдавать историю сделок ( пример 2021-10-13 )
ticker string Тикер\код инструмента, ISIN для облигаций
limit int32 Possible values: <= 1000Количество возвращаемых записей
from int64 Начальный номер сделки для фильтра результатов
descending boolean Флаг обратной сортировки выдачи
side string Possible values: [buy, sell] Направление сделки:
format string Possible values: [Simple, Slim, Heavy] Формат возвращаемого сервером JSON

(4) История сделок по инструменту

https://apidev.alor.ru/md/v2/Stats/:exchange/:portfolio/history/trades/:symbol
*/

type TradeRequest struct {
	symbol     string // Тикер\код инструмента, ISIN для облигаций
	dateFrom   string // Начиная с какой даты отдавать историю сделок ( пример 2021-10-13 )
	limit      int32  // <= 1000 Количество возвращаемых записей
	descending bool   // Флаг обратной сортировки выдачи
	side       string // [buy, sell] Направление сделки:
	format     string // [Simple, Slim, Heavy] Формат возвращаемого сервером JSON

}

func (c *Client) GetTrades(portfolio string, params TradeRequest) ([]Trade, error) {

	result := make([]Trade, 0)
	return result, nil
}
