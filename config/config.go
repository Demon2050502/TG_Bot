package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	APIKei  string
	OwnerID int64
	Timeout int
}

func GetConfig() *BotConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}

	conf := BotConfig{}
	var ok bool

	if	conf.APIKei, ok = os.LookupEnv("APIKei"); ok == false{
		log.Fatal("Переменная окружения APIKei не установлена")
	}

	// conf.OwnerID = 0
	conf.Timeout = 60

	return &conf
}
