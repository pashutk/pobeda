package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	fmt.Println("Poehali!")

	key := ""
	if len(os.Args) > 1 {
		key = os.Args[1]
	} else {
		log.Fatalln("No api key")
	}

	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := update.Message

		if msg == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if cmd := msg.Command(); cmd != "" {
			log.Println(cmd)
			log.Println(msg)
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)
	}

	// to := "LCA"
	// if len(os.Args) > 1 {
	// 	to = os.Args[1]
	// }

	// flights := getFlightsForRegion(to)

	// for _, flight := range flights {
	// 	fmt.Println(flight)
	// }

	fmt.Println("Pobeda!")
}
