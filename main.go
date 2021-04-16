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

		//	Token:  "1762186330:AAELm54VB5FAvLDPeoFPYSnkHOuWOLaj_wk",
		Token:  "1088448942:AAGbDckx7aVCoa005afOE2bVwVejgiPMS4c",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Bot started")
	logger := NewBotLog(b)
	logger.Log("Bot started!")

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
	playerStorage := NewPlayerStorage(&holdeStorage)

	gspread := NewGspreadHoldes()

	gameSettings := gspread.ReadSettings()

	var (

		// Окно приветсвтия - регистрация  -

		// кнопка вызова главного меню
		//btnMainMenu = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("В главное меню")
		// Главное меню - Настройки | Считать
		menuMain      = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		btnSettings   = menuMain.Data("Настройки", "settings")
		btnCalculator = menuMain.Data("Поместья", "holdes")

		// Настройки - Oпция 1 | Oпция 2 | Oпция 3
		// Считать  -  Имя игрока - Добавить поместья - Снять кэш - Улучшить поместье.
		//                             |      |
		//								> ----^

		// Universal markup builders.
		// selector = &tb.ReplyMarkup{}

		addHoldeMenuKeyboard = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		addHoldeButton       = addHoldeMenuKeyboard.Data("Добавить поместье", "add_holde")
		addHoldeCancelButton = addHoldeMenuKeyboard.Data("Нет, другое", "cancel_add_holde")

		addNewHoldeMenuKeyboard = &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		calcHoldeButton    = addNewHoldeMenuKeyboard.Data("Обсчитать поместья", "calc_all_holde")
		addHoldeMoreButton = addNewHoldeMenuKeyboard.Data("Добавить еще поместий", "add_more_holde")

		menuSettingsKbd      = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		changeUsernameButton = menuSettingsKbd.Data("Изменить имя", "change_user_name")
		changeUserLocButton  = menuSettingsKbd.Data("Изменить локацию", "change_user_loc")
		backToMainMenu       = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Data("В главное меню", "back_main_menu")

		// Reply buttons.
		//btnHelp = menu.Text("ℹ Help")
		//	btnSettings = menu.Text("⚙ Settings")

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind,
		// as it has to be for callback routing to work.
		//
		// btnPrev = selector.Data("⬅", "prev")
		//		btnNext = selector.Data("➡", "next")
	)

	var users UserStorager
	users = NewUsers()

	log.Println(holdeStorage)
	log.Println(playerStorage)
	log.Println(users)
	// menu.Reply(
	// 	menu.Row(btnHelp),
	// 	menu.Row(btnSettings),
	// )
	// selector.Inline(
	// 	selector.Row(btnPrev, btnNext),
	// )

	// UpdateUserState := func(id UserID, state DialogState) {
	// 	user, _ := users.Get(id)
	// 	user.SetState(MainMenu)
	// 	users.Update(id, user)
	// }

	// mainMenu := DialogNode{
	// 	Content: DialogContent{
	// 		Message: "Выберите, что бы вы хотели сделать",
	// 		Media:   nil,
	// 	},
	// 	Keyboard: &tb.ReplyMarkup{},
	// }
	// mainMenu.Keyboard.Reply(
	// 	mainMenu.Keyboard.Row(btnCalculator),
	// 	mainMenu.Keyboard.Row(btnSettings),
	// )

	menuMain.Inline(
		menuMain.Row(btnSettings),
		menuMain.Row(btnCalculator),
	)

	addHoldeMenuKeyboard.Inline(
		addHoldeMenuKeyboard.Row(addHoldeButton),
		addHoldeMenuKeyboard.Row(addHoldeCancelButton),
	)

	addNewHoldeMenuKeyboard.Inline(
		addNewHoldeMenuKeyboard.Row(addHoldeMoreButton),
		addNewHoldeMenuKeyboard.Row(calcHoldeButton),
	)

	menuSettingsKbd.Inline(
		menuSettingsKbd.Row(changeUsernameButton),
		menuSettingsKbd.Row(changeUserLocButton),
		menuSettingsKbd.Row(backToMainMenu),
	)

	// Command: /start <PAYLOAD>
	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		} else {
			logger.AddSubscriber(tb.ChatID(m.Chat.ID))
			logger.Log("added new subscriber")
		}
		id := UserID(m.Sender.ID)
		user, created := users.GetOrCreate(id)
		user.SetState(MainMenu)
		users.Update(id, user)
		//UpdateUserState(id,MainMenu)

		//  Для нового пользовтеля отсылаем приветсвенное соообщение
		if created {
			//
			b.Send(m.Sender, fmt.Sprintf("Рад с тобой познакомиться %s. Я бот и я буду тебе помогать в экономике", m.Sender.FirstName))
		} else {
			b.Send(m.Sender, fmt.Sprintf("с возвращением %s", m.Sender.FirstName))
		}
		b.Send(m.Sender, "Начнем", menuMain)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		id := UserID(m.Sender.ID)

		user, err := users.Get(id)
		if err != nil {
			b.Send(m.Sender, "Случилась какая то ошибка. давай начнем заово. Жми /start")
			return
		}

		log.Println("Text handler state", user.State)
		switch user.State {
		case IDLE:
			{
				b.Send(m.Sender, "Не понимаю. Жми /start")
				return
			}
		case AddHolde:
			{
				holdeID, err := strconv.Atoi(m.Text)
				if err != nil {
					b.Send(m.Sender, "Неправильынй ноvер поместья")
					return
				}
				if holdeID < 0 || holdeID > HoldesNumber {
					b.Send(m.Sender, "Неправильынй нмоер поместья")
					return
				}
				holde, err := holdeStorage.Get(holdeID)
				if err != nil {
					b.Send(m.Sender, "Неправильынй нмоер поместья")
					return
				}
				user.CurrHolde = holdeID
				user.Save()
				users.Update(id, user)

				b.Send(m.Sender, holde.ResponseText(), addHoldeMenuKeyboard)

			}

		case EnterDice:
			{
				dice, err := strconv.Atoi(m.Text)
				if err != nil {
					b.Send(m.Sender, "Неправильынй бросок. Введите 1- 10")
					return
				}
				if dice < 1 || dice > 10 {
					b.Send(m.Sender, "Неправильынй бросок. Введите 1- 10")
					return
				}
				request := HoldeRequestItem{
					HoldeID: user.CurrHolde,
					Dice:    dice,
				}

				user.CurrPlayer.Request.Holdes = append(user.CurrPlayer.Request.Holdes, request)
				playerStorage.Update(*user.CurrPlayer)
				log.Println(user.CurrPlayer.Request)
				b.Send(m.Sender, "Поместье добавлено! Хотите добавить еще? Введите номер поместья илли обсчитайте поместья", addNewHoldeMenuKeyboard)
				user.SetState(AddHolde)
			}

		case EnterPlayerName:
			{
				log.Println("Text handler state", "EnterPlayerName:")
				playerName := m.Text
				player, created := playerStorage.GetOrCreate(playerName)
				if created {
					b.Send(m.Sender, "Свежее мясо!")
				} else {
					b.Send(m.Sender, "Знакомые всё лица!")
				}
				user.CurrPlayer = &player
				user.State = AddHolde

				b.Send(m.Sender, "Очень хорошо. Введи номер поместья")
				user.SetState(AddHolde)
			}

		case EnterUserName:
			{
				log.Println("Text handler state", "EnterUserName:")
				user.Name = m.Text

				user.State = UserSettings
				b.Send(m.Sender, "Приятно познакомиться, "+user.Name)
				b.Send(m.Sender, "Твой профиль в игре \n Имя"+user.Name+"\n Локация: "+user.Location+"\nыбери, что изменить? ", menuSettingsKbd)
				user.SetState(AddHolde)
			}
		}
		users.Update(id, user)

	})

	b.Handle(&calcHoldeButton, func(m *tb.Message) {
		log.Println("calcHoldeButton", "EnterPlayerName:")
		id := UserID(m.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(m.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}
		resp, err := user.CurrPlayer.HandleReq()
		users.Update(id, user)
		b.Send(m.Sender, resp.Show())

	})
	// On reply button pressed (message)
	// b.Handle(&btnHelp, func(m *tb.Message) {
	// 	log.Println("User", m.Sender.ID, m.Sender.FirstName, m.Sender.LastName)
	// 	log.Println("Message", m.Text)
	// 	b.Send(m.Sender, "Hello World!", selector)
	// })

	b.Handle(&btnCalculator, func(m *tb.Message) {
		log.Println("calcHoldeButton", "EnterPlayerName:")
		id := UserID(m.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(m.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}
		user.SetState(EnterPlayerName)
		users.Update(id, user)

		b.Send(m.Sender, "Сообщи, пожалуйста, мне имя игрока")
	})

	b.Handle(&addHoldeButton, func(c *tb.Callback) {
		id := UserID(c.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}

		user.SetState(EnterDice)
		users.Update(id, user)

		b.Send(c.Sender, "Попроси игрока бросить кубик d10")
		b.Respond(c, &tb.CallbackResponse{

			Text:      "Поместье добавлено",
			ShowAlert: false,
		})
	})

	b.Handle(&addHoldeMoreButton, func(c *tb.Callback) {
		id := UserID(c.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}

		user.SetState(AddHolde)
		users.Update(id, user)

		b.Send(c.Sender, "Введи номер поместья")
		b.Respond(c, &tb.CallbackResponse{
			Text:      "Еще поместий",
			ShowAlert: false,
		})
	})
	// calcHoldeButton
	b.Handle(&calcHoldeButton, func(c *tb.Callback) {
		id := UserID(c.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}
		user.SetState(HoldeCalc)

		resp, err := user.CurrPlayer.HandleReq()
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}

		users.Update(id, user)
		b.Send(c.Sender, resp.Show())
		b.Respond(c, &tb.CallbackResponse{
			Text:      "Еще поместий",
			ShowAlert: false,
		})

	})

	b.Handle(&btnSettings, func(c *tb.Callback) {
		id := UserID(c.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}
		user.SetState(UserSettings)
		users.Update(id, user)

		b.Send(c.Sender, "Твой профиль в игре \n Имя"+user.Name+"\n Локация: "+user.Location+"\nыбери, что изменить? ", menuSettingsKbd)

		b.Respond(c, &tb.CallbackResponse{
			Text:      "Настройки",
			ShowAlert: false,
		})
	})

	b.Handle(&changeUsernameButton, func(c *tb.Callback) {
		id := UserID(c.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}
		user.SetState(EnterUserName)
		users.Update(id, user)

		b.Send(c.Sender, "Введи свое имя")

		b.Respond(c, &tb.CallbackResponse{
			Text:      "",
			ShowAlert: false,
		})
	})

	//-----------
	locSelectKbd := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	locationNum := len(gameSettings.Locations)
	locationSelectButtons := make([]tb.Btn, locationNum)
	for i, loc := range gameSettings.Locations {
		locationSelectButtons[i] = locSelectKbd.Data(loc, "loc_sel_"+strconv.Itoa(i), loc)
		b.Handle(&locationSelectButtons[i], func(c *tb.Callback) {
			id := UserID(c.Sender.ID)
			user, err := users.Get(id)
			if err != nil {
				b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
				return
			}

			user.Location = c.Data // i m not sure what this is correct
			log.Println(*c, c.Data)
			user.SetState(UserSettings)
			users.Update(id, user)
			b.Send(c.Sender, "Твой профиль в игре \n Имя"+user.Name+"\n Локация: "+user.Location+"\nыбери, что изменить? ", menuSettingsKbd)
		})

	}
	locationSelectButtonsRows := make([]tb.Row, locationNum)
	for i, _ := range locationSelectButtons {
		locationSelectButtonsRows[i] = locSelectKbd.Row(locationSelectButtons[i])
	}

	locSelectKbd.Inline(locationSelectButtonsRows...)

	b.Handle(&changeUserLocButton, func(c *tb.Callback) {
		id := UserID(c.Sender.ID)
		user, err := users.Get(id)
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}
		user.SetState(EnterUserLoc)
		users.Update(id, user)

		b.Send(c.Sender, "Выбери локацию", locSelectKbd)

		// add

		b.Respond(c, &tb.CallbackResponse{
			Text:      "",
			ShowAlert: false,
		})
	})

	// Send hand

	// // On inline button pressed (callback)
	// b.Handle(&btnPrev, func(c *tb.Callback) {
	// 	// ...
	// 	// Always respond!
	// 	b.Respond(c, &tb.CallbackResponse{
	// 		Text: c.Message.Text,
	// 	})
	// })

	// type DialogNode struct {
	// 	// Текст ии дугое сообщение оторжаео при переходе на данный узел
	// 	Content DialogContent
	// 	// Клавиатура
	// 	Keyboard *tb.ReplyMarkup
	// 	//
	// }

	// On reply button pressed (message)
	// b.Handle(&btnMainMenu, func(m *tb.Message) {

	// 	log.Println("main menu", m.Text)

	// 	b.Send(m.Sender, "Выберите, что бы вы хотели сделать", mainMenu)

	// })

	b.Start()

}
