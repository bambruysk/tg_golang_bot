package main

import (
	"fmt"
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

	// b.Handle("/hello", func(m *tb.Message) {
	// 	log.Println("Message", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName)
	// 	b.Send(m.Sender, "Hello World!")
	// })

	// b.Handle(tb.OnText, func(m *tb.Message) {

	// 	log.Println("User", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName)
	// 	log.Println("Message", m.Text)
	// 	b.Send(m.Sender, "Hello World!")
	// })

	// b.Start()

	var (

		// Окно приветсвтия - регистрация  -

		// кнопка вызова главного меню
		btnMainMenu = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("В главное меню")
		// Главное меню - Настройки | Считать
		btnSettings   = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("Настройки")
		btnCalculator = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("Поместья")

		// Настройки - Oпция 1 | Oпция 2 | Oпция 3
		// Считать  -  Имя игрока - Добавить поместья - Снять кэш - Улучшить поместье.
		//                             |      |
		//								> ----^

		// Universal markup builders.
		menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		selector = &tb.ReplyMarkup{}

		// Reply buttons.
		btnHelp = menu.Text("ℹ Help")
		//	btnSettings = menu.Text("⚙ Settings")

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind,
		// as it has to be for callback routing to work.
		//
		btnPrev = selector.Data("⬅", "prev")
		//		btnNext = selector.Data("➡", "next")
	)

	var users UserStorager
	users = NewUsers()

	// menu.Reply(
	// 	menu.Row(btnHelp),
	// 	menu.Row(btnSettings),
	// )
	// selector.Inline(
	// 	selector.Row(btnPrev, btnNext),
	// )

	UpdateUserState := func(id UserID, state DialogState) {
		user, _ := users.Get(id)
		user.SetState(MainMenu)
		users.Update(id, user)
	}

	mainMenu := DialogNode{
		Content: DialogContent{
			Message: "Выберите, что бы вы хотели сделать",
			Media:   nil,
		},
		Keyboard: &tb.ReplyMarkup{},
	}
	mainMenu.Keyboard.Reply(
		mainMenu.Keyboard.Row(btnCalculator),
		mainMenu.Keyboard.Row(btnSettings),
	)

	// Command: /start <PAYLOAD>
	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		id := UserID(m.Sender.ID)
		user, created := users.GetOrCreate(id)
		user.SetState(MainMenu)
		users.Update(id, user)

		//  Для нового пользовтеля отсылаемприветсвенное соообщение
		if created {
			//
			b.Send(m.Sender, fmt.Sprintf("Рад с тобой познаокмиться %s. Я бот и я буду тебе помогать в экономике", m.Sender.FirstName))
		} else {
			b.Send(m.Sender, fmt.Sprintf("с возвращением %s", m.Sender.FirstName))
		}
		b.Send(m.Sender, "Начнем", mainMenu)
	})

	// On reply button pressed (message)
	b.Handle(&btnHelp, func(m *tb.Message) {
		log.Println("User", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName)
		log.Println("Message", m.Text)
		b.Send(m.Sender, "Hello World!", selector)
	})

	b.Handle(&btnCalculator, func(m *tb.Message) {
		UpdateUserState(UserID(m.Sender.ID), HoldeCalc)
		b.Send(m.Sender, "Сообщи, пожалуйста, мне имя игрока")
	})

	// Send hand

	// On inline button pressed (callback)
	b.Handle(&btnPrev, func(c *tb.Callback) {
		// ...
		// Always respond!
		b.Respond(c, &tb.CallbackResponse{
			Text: c.Message.Text,
		})
	})

	type DialogNode struct {
		// Текст ии дугое сообщение оторжаео при переходе на данный узел
		Content DialogContent
		// Клавиатура
		Keyboard *tb.ReplyMarkup
		//
	}

	// On reply button pressed (message)
	b.Handle(&btnMainMenu, func(m *tb.Message) {

		log.Println("main menu", m.Text)

		b.Send(m.Sender, "Выберите, что бы вы хотели сделать", mainMenu)

	})

	b.Start()

}
