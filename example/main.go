package main

import (
	"context"
    "os"
	"log/slog"

    "github.com/Ruvad39/go-alor"
)


func main(){

	ctx := context.Background()
	// для отладки
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{	Level: slog.LevelDebug,	})		
	logger_ := slog.New(handler)

	
	// alor.WithServer("https://apidev.alor.ru") // подключить Тестовый контур
	//client, err := alor.NewClient( alor.WithLogger(logger_), alor.WithServer("https://apidev.alor.ru"))
	// создание клиента		
	// по умолчаию подключен Боевой контур
	client, err := alor.NewClient( alor.WithLogger(logger_))

	if err != nil {
		slog.Error("ошибка создания alor.client: " + err.Error())
	}

	// получить текущее время сервера
	// без авторизации задержка по времени 15 минут
	servTime, err := client.GetTime(ctx)
	if err != nil {
		slog.Error("ошибка получения текущего времени: " + err.Error())
		return
	 }
	slog.Info("time", "servTime",servTime.String()) 


}