# Golang SDK для работы с функционалом [Alor OpenAPI V2](https://alor.dev/docs) брокера [Алор](https://www.alorbroker.ru)




## Установка

```bash
go get github.com/Ruvad39/go-alor
```

**без авторизации задержка по времени 15 минут**
## какой api реализован 
```go
// текущее время сервера
GetTime(ctx context.Context) (time.Time, error)
// GetSecurities получить список торговых инструментов
GetSecurities(ctx context.Context, params Params) ([]Security, error)
// GetSecurity получить параметры по торговому инструменту
GetSecurity(ctx context.Context, board, symbol string) (Security, error)
// GetQuotes Получение информации о котировках для выбранных инструментов
GetQuotes(ctx context.Context, symbols string) ([]Quote, error)
// GetQuote Получение информации о котировках для одного выбранного инструмента.
GetQuote(ctx context.Context, symbol string) (Quote, error)
// GetPositions получение информации о позициях
GetPositions(ctx context.Context, portfolio string) ([]Position, error)
// GetHistory Запрос истории для выбранных биржи и инструмента
GetHistory(ctx context.Context, symbol string, interval Interval, from, to int64) (History, error)
// GetCandles Запрос истории свечей для выбранного инструмента (вызывает GetHistory)
GetCandles(ctx context.Context, symbol string, interval Interval, from, to int64) ([]Candle, error)
```
## Примеры

### Пример создание клиента. Получение текущего времени сервера

```go
ctx := context.Background()

// создание клиента     
refreshToken, _ := os.LookupEnv("ALOR_REFRESH")
// по умолчаию подключен Боевой контур
client, err := alor.NewClient(refreshToken)

if err != nil {
   slog.Error("ошибка создания alor.client: " + err.Error())
}

// получить текущее время сервера
// без авторизации задержка по времени 15 минут
servTime, err := client.GetTime(ctx)
if err != nil {
    fmt.Println(err) 
    return
 }
slog.Info("time", "servTime",servTime.String()) 

```

### Получить параметры по торговому инструменту

```go
board := "TQBR"
symbol :=  "SBER"
sec, err := client.GetSecurity(ctx, board, symbol)
slog.Info("securities",
    "Symbol", sec.Symbol,
    "Exchange", sec.Exchange,
    "board", sec.Board,
    "ShortName", sec.ShortName,
    "LotSize", sec.LotSize,
    "MinStep", sec.MinStep,
)

// запрос списка инструментов
// sector = FORTS, FOND, CURR
// Если не указано иное значение параметра limit, в ответе возвращается только 25 объектов за раз
params := alor.Params{
    Sector: "FOND",
    Board:  "TQBR",
    Query:  "",
    Limit:  400,
}
sec, err := client.GetSecurities(ctx, params)
if err != nil {
slog.Info("main.GetSecurity", "err", err.Error())
return
}
slog.Info("GetSecurity", "кол-во", len(sec))
```
### Получение информации о котировках для выбранных инструментов.
```go
// Принимает несколько пар биржа-тикер. Пары отделены запятыми. Биржа и тикер разделены двоеточием
symbols := "MOEX:SIM4,MOEX:SBER"
sec, err := client.GetQuotes(ctx, symbols)
if err != nil {
	slog.Info("main.GetQuotes", "err", err.Error())
	return
}
for _, q := range sec {
slog.Info("Quotes",
    "symbol", q.Symbol,
    "description", q.Description,
    "lastPrice", q.LastPrice,
    "Bid", q.Bid,
    "ask", q.Ask,
    "LotValue", q.LotValue,
    "LotSize", q.LotSize,
    "OpenInterest", q.OpenInterest,
    "LastTime", q.LastTime().String(),
    )
}

// Получение информации о котировках для одного выбранного инструмента.
// Указываем тикер без указания биржи. Название биржи берется по умолчанию
symbol := "SRM4"
q, err := client.GetQuote(ctx, symbol)
slog.Info("Quotes",
    "symbol", q.Symbol,
    "description", q.Description,
    "lastPrice", q.LastPrice,
    "Bid", q.Bid,
    "ask", q.Ask,
    "LotValue", q.LotValue,
    "LotSize", q.LotSize,
    "OpenInterest", q.OpenInterest,
    "LastTime", q.LastTime().String(),
    )




```

### другие примеры смотрите [тут](/example)