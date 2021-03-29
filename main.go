
package main

import (
	"log"
	"time"


	tb "gopkg.in/tucnak/telebot.v2"
)




func main() {
	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		URL: "https://api.telegram.org",

		Token:  "1762186330:AAELm54VB5FAvLDPeoFPYSnkHOuWOLaj_wk",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Bot started")

	b.Handle("/hello", func(m *tb.Message) {
		log.Println("Message", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName )
		b.Send(m.Sender, "Hello World!")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {

		log.Println("User", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName )
		log.Println("Message", m.Text )
		b.Send(m.Sender, "Hello World!")
	})

	b.Start()
}