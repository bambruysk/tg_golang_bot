package main

import (
	"context"
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/joho/godotenv"
)

var logger BotLog

var b *tb.Bot

func main() {
	log.Println("Bot read env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	tgApiKey := os.Getenv("TG_API_KEY")
	if tgApiKey == "" {
		log.Fatal("TG_API_KEY env not set")
	}
	log.Println("Bot read env", tgApiKey)
	// test db connection

	conn, err := ConnectDB()
	if err != nil {
		log.Println("Database connection error", err)
	} else {
		log.Println("Database connection succees")
	}
	conn.Close(context.Background())

	b, err = tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL: "https://api.telegram.org",

		//	Token:  "1762186330:AAELm54VB5FAvLDPeoFPYSnkHOuWOLaj_wk",
		Token:  tgApiKey,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal("Bot cratione err", err)
		return
	}

	log.Println("Bot started")
	logger = NewBotLog(b)
	logger.Log("Bot started!")

	AddBotLogic(b)

	b.Start()

}
