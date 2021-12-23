package main

import (
	"log"

	"encoding/json"
	"fmt"
	"os"

	utopiago "github.com/Sagleft/utopialib-go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TelegramBotToken string
	UtpToken         string
	UtpPort          int
	IdChannel        string
}

func main() {

	//read congig
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(configuration.TelegramBotToken)
	fmt.Println(configuration.UtpToken)
	fmt.Println(configuration.UtpPort)

	// bot-token

	bot, err := tgbotapi.NewBotAPI(configuration.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// ini channel
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updatesChann := bot.GetUpdatesChan(ucfg)

	//utp
	client := utopiago.UtopiaClient{
		Protocol: "http",
		Token:    configuration.UtpToken,
		Host:     "127.0.0.1",
		Port:     configuration.UtpPort,
	}

	fmt.Println(client.GetBalance())

	var send = false
	//send bool

	// update
	for {
		select {

		case update := <-updatesChann:
			// User bot
			UserName := update.Message.From.UserName

			// ID chat.

			ChatID := update.Message.Chat.ID

			// Text massage user
			Text := update.Message.Text

			log.Printf("[%s] %d %s", UserName, ChatID, Text)

			if send {

				client.SendChannelMessage(configuration.IdChannel, Text)
				send = false
				log.Println("massage go")
			}

			//commands

			switch Text {

			case "/SendMassage":

				//send to channel

				log.Println("/send true")

				send = true

			case "/How":

				log.Println("/How")

				msg := tgbotapi.NewMessage(ChatID, "create config.json and use")

				bot.Send(msg)

			default:

				fmt.Println("commands")

				reply := "Commands:\n /SendMassage \n /How"
				msg := tgbotapi.NewMessage(ChatID, reply)

				bot.Send(msg)

			}

		}

	}
}
