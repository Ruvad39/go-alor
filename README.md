# Библиотека, которая позволяет работать с функционалом [ALOR OpenAPI V2](https://www.alorbroker.ru/trading/openapi)  брокера [Алор](https://www.alorbroker.ru/) из GO



## Установка

```bash
go get github.com/Ruvad39/go-alor
```

на текущий момент будут реализованы только методы которые не требует авторизации  **без авторизации задержка по времени 15 минут**
## какой api реализован 
```go

// текущее время сервера
GetTime(ctx context.Context) (time.Time, error)


```
## Примеры

### Пример создание клиента. Получение текущего времени сервера

```go
ctx := context.Background()

// создание клиента     
// по умолчаию подключен Боевой контур
client, err := alor.NewClient()

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

// получение текущих рыночных данных по иструменту
market, err := client.GetQuotes(ctx, "SBER")
//market, err := client.GetQuotes(ctx, "RTS-6.24")
if err != nil {
    slog.Error("ошибка GetQuotes: " + err.Error())
}
slog.Info("Quotes", 
    "symbol",       market[0].Symbol,
    "description",  market[0].Description,
    "lastPrice",    market[0].LastPrice,
    "Bid",          market[0].Bid,
    "ask",          market[0].Ask,
    "LotValue",     market[0].LotValue,
    "LotSize",      market[0].LotSize,
    "OpenInterest", market[0].OpenInterest,
    "LastTime",     market[0].LastTime().String(),
)

```

### Пример получения свечей

```go
// TODO
```

### другие примеры смотрите [тут](/example)