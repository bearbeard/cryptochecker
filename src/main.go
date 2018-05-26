package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"encoding/json"
)

type Config struct {
	TelegramBotToken string
}

func main() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Configuration file %s read successful.", file.Name())
	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	//bot, err := tgbotapi.NewBotAPI("617766288:AAFZYpiVwLi37oLJvRpMo3Hl0ziygiWcanc")
	checkError(err)
	log.Printf("Authorized on account %s.", bot.Self.UserName)
	processUpdate(bot)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}


