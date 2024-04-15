# Библиотека, которая позволяет работать с функционалом [ALOR OpenAPI V2](https://www.alorbroker.ru/trading/openapi)  брокера [Алор](https://www.alorbroker.ru/) из GO



## Установка

```bash
go get github.com/Ruvad39/go-alor
```
## какой api реализован 
на текущий момент будет реализованы только методы которые не требует авторизации
**без авторизации задержка по времени 15 минут**

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

```

### Пример получения свечей

```go
// TODO
```

### другие примеры смотрите [тут](/example)