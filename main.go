package main

import (
	"fmt"
	"log"
	"strconv"
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

	holdeStorage := NewHoldeStorage()

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

		addHoldeMenuKeyboard =  &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		addHoldeButton = addHoldeMenuKeyboard.Text("Добавить поместье")
		addHoldeCancelButton = addHoldeMenuKeyboard.Text("Нет, другое")



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


	b.Handle(tb.OnText, func (m* tb. Message) {
		id := UserID(m.Sender.ID)

		user, err :=  users.Get(id)
		if err != nil {
			b.Send(m.Sender, "Сkучилась какая то ошибка. давай начнем заово. Жми /start")
			return
		}
		switch user.State {
		case IDLE :{
			b.Send(m.Sender, "Не понимаю. Жми /start")
			return
		}	
		case AddHolde : {
			holde_id := strconv.Atoi( m.Text)
			if holde_id < 0 || holde_id > HoldesNumber {
				b.Send(m.Sender, "Неправильынй нмоер поместья")
				return
			}
			holde, err := holdeStorage.Get(holde_id)
			if err !=  nil {
				b.Send(m.Sender, "Неправильынй нмоер поместья")
				return
			} 
			user.CurrHolde =  holde_id
			b.Send(m.Sender, holde.ResponseText(),  addHoldeMenuKeyboard)

		}

		case EnterDice : {
			dice := strconv.Atoi( m.Text)
			if dice < 1 || dice > 10 {
				b.Send(m.Sender, "Неправильынй бросок. Введеите 1- 10")
				return
			}
			HoldeRequestItem{
				HoldeID: 0,
				Dice:    0,
			}

		} 

		  user.State



	} )

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
