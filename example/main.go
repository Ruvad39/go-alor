package main

import (
	"context"
	"github.com/Ruvad39/go-alor"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("No .env file found")
		//slog.Error(err.Error())
	}
}

func main() {
	ctx := context.Background()

	//refreshToken := os.Getenv("ALOR_REFRESH")

	//refreshToken, ok := os.LookupEnv("ALOR_REFRESH")
	refreshToken, _ := os.LookupEnv("ALOR_REFRESH")

	//if ok {
	//	fmt.Println(refreshToken)
	//}
	// создание клиента
	client := alor.NewClient(refreshToken)
	client.Debug = true
	//slog.Info("main.", "client.version()", client.Version())

	// получить текущее время сервера
	// без авторизации задержка по времени 15 минут
	servTime, err := client.GetTime(ctx)
	if err != nil {
		slog.Error("ошибка получения текущего времени: " + err.Error())
		return
	}
	slog.Info("time", "servTime", servTime.String())

}
