package main

import (
	"fmt"

	tgadapter "github.com/Demon2050502/TG_Bot/adapter"
	"github.com/Demon2050502/TG_Bot/config"
	"github.com/Demon2050502/TG_Bot/service"
)

func main() {
	fmt.Println("Starting application...")

	conf := config.GetConfig()
	adapter := tgadapter.NewAdapter(*conf)
	service := service.NewService(adapter)

	service.StartServe()

	// log.Printf("Бот авторизован как @%s", bot.Self.UserName)
}
