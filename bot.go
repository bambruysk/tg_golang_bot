package main

import (
	"fmt"
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

var gspread = NewGspreadHoldes()
var gameSettings = gspread.ReadSettings()

var holdeStorage = NewHoldeStorage()
var playerStorage = NewPlayerStorage(holdeStorage)

var (

	// Окно приветсвтия - регистрация  -

	// кнопка вызова главного меню
	//btnMainMenu = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("В главное меню")
	// Главное меню - Настройки | Считать
	menuMain      = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	btnSettings   = menuMain.Data("Настройки", "settings")
	btnCalculator = menuMain.Data("Поместья", "holdes")

	// Настройки - Oпция 1 | Oпция 2 | Oпция 3
	// Считать  -  Имя игрока - Добавить поместья - Снять кэш - Улучшить поместье.
	//                             |      |
	//								> ----^

	// Universal markup builders.
	// selector = &tb.ReplyMarkup{}

	addHoldeMenuKeyboard = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}

	addHoldeButton       = addHoldeMenuKeyboard.Data("Добавить это поместье", "add_holde")
	addHoldeCancelButton = addHoldeMenuKeyboard.Data("Нет, другое", "cancel_add_holde")

	addNewHoldeMenuKeyboard = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}

	calcHoldeButton    = addNewHoldeMenuKeyboard.Data("Обсчитать поместья", "calc_all_holde")
	addHoldeMoreButton = addNewHoldeMenuKeyboard.Data("Добавить еще поместий", "add_more_holde")

	upgradeKeyboard    = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	letUpgradeButton   = upgradeKeyboard.Data("Улушить поместье", "let_upgrade")
	endCalculateButton = upgradeKeyboard.Data("Закочнить расчеты", "end_caluclate")

	locSelectKbd = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}

	menuSettingsKbd      = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	changeUsernameButton = menuSettingsKbd.Data("Изменить имя", "change_user_name")
	changeUserLocButton  = menuSettingsKbd.Data("Изменить локацию", "change_user_loc")
	backToMainMenuButton = (&tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}).Data("В главное меню", "back_main_menu")

	holdeUpgradeKbd        = &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	holdeUpgradeApproveBtn = holdeUpgradeKbd.Data("Улучшить", "upgrade_approve")
	holdeUpgradeCancelBtn  = holdeUpgradeKbd.Data("Отмена", " upgrade_cancel")
)
var users UserStorager

func AddBotLogic(b *tb.Bot) {

	// gspread := NewGspreadHoldes()

	// gameSettings := gspread.ReadSettings()

	users = NewUsers()

	log.Println(holdeStorage)
	log.Println(playerStorage)
	log.Println(users)

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
		menuSettingsKbd.Row(backToMainMenuButton),
	)

	upgradeKeyboard.Inline(
		upgradeKeyboard.Row(letUpgradeButton),
		upgradeKeyboard.Row(endCalculateButton),
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
			b.Send(m.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
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
					b.Send(m.Sender, "Неправильынй номер поместья")
					return
				}
				if holdeID < 0 || holdeID > HoldesNumber {
					b.Send(m.Sender, "Неправильынй номер поместья")
					return
				}
				holde, err := holdeStorage.Get(holdeID)
				if err != nil {
					b.Send(m.Sender, "Неправильынй номер поместья")
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
					b.Send(m.Sender, "Неправильынй бросок. Введите 1-10")
					return
				}
				if dice < 1 || dice > 10 {
					b.Send(m.Sender, "Неправильынй бросок. Введите 1-10")
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
					b.Send(m.Sender, "Знакомые все лица!")
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

	b.Handle(&btnCalculator, func(c *tb.Callback) {
		log.Println("btnCalculator")
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}
		user.SetState(EnterPlayerName)
		users.Update(id, user)

		b.Send(c.Sender, "Сообщи, пожалуйста, мне имя игрока")
	})

	b.Handle(&addHoldeButton, func(c *tb.Callback) {
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}
		user.SetState(EnterDice)
		users.Update(id, user)

		b.Send(c.Sender, "Попроси игрока бросить кубик d10")
		b.Respond(c, &tb.CallbackResponse{

			Text: "Поместье добавлено",
		})
	})

	b.Handle(&addHoldeMoreButton, func(c *tb.Callback) {
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}

		user.SetState(AddHolde)

		b.Send(c.Sender, "Введи номер поместья")
		b.Respond(c, &tb.CallbackResponse{
			Text: "Еще поместий",
		})
		users.Update(id, user)
	})
	// calcHoldeButton
	b.Handle(&calcHoldeButton, func(c *tb.Callback) {
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}
		user.SetState(HoldeCalc)

		resp, err := user.CurrPlayer.HandleReq()
		if err != nil {
			b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
			return
		}

		b.Send(c.Sender, resp.Show(), upgradeKeyboard)
		b.Respond(c, &tb.CallbackResponse{
			Text: "Еще поместий",
		})
		users.Update(id, user)

	})

	b.Handle(&btnSettings, func(c *tb.Callback) {

		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}

		user.SetState(UserSettings)
		b.Send(c.Sender, user.ShowProfile(), menuSettingsKbd)

		b.Respond(c, &tb.CallbackResponse{
			Text: "Настройки",
		})

		users.Update(id, user)
	})

	b.Handle(&changeUsernameButton, func(c *tb.Callback) {
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}
		user.SetState(EnterUserName)

		b.Send(c.Sender, "Введи свое имя")

		b.Respond(c, &tb.CallbackResponse{
			Text: "",
		})
		users.Update(id, user)
	})

	//-----------
	//locSelectKbd := &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	locationNum := len(gameSettings.Locations)
	locationSelectButtons := make([]tb.Btn, locationNum)
	for i, loc := range gameSettings.Locations {
		locationSelectButtons[i] = locSelectKbd.Data(loc, "loc_sel_"+strconv.Itoa(i), loc)
		b.Handle(&locationSelectButtons[i], func(c *tb.Callback) {
			id, user, notFound := GetUserFromCallbackByID(c, users, b)
			if notFound {
				return
			}

			b.Send(c.Sender, "Введи свое имя")

			user.Location = c.Data // i m not sure what this is correct
			log.Println(*c, c.Data)
			user.SetState(UserSettings)
			b.Send(c.Sender, "Твой профиль в игре \nИмя: \t"+user.Name+"\n Локация: \t"+user.Location+"\n Выбери, что изменить? ", menuSettingsKbd)
			users.Update(id, user)
		})

	}
	locationSelectButtonsRows := make([]tb.Row, locationNum)
	for i, _ := range locationSelectButtons {
		locationSelectButtonsRows[i] = locSelectKbd.Row(locationSelectButtons[i])
	}

	locSelectKbd.Inline(locationSelectButtonsRows...)

	holdeUpgradeKbd.Inline(
		holdeUpgradeKbd.Row(holdeUpgradeApproveBtn),
		holdeUpgradeKbd.Row(holdeUpgradeCancelBtn),
	)

	b.Handle(&letUpgradeButton, LetUpgradeButtonHandler)

	b.Handle(&holdeUpgradeApproveBtn, HoldeUpgradeApproveButtonHandler)

	b.Handle(&holdeUpgradeCancelBtn, HoldeUpgradeCancelBtnHandler)

	b.Handle(&endCalculateButton, EndCalculateButtonHandler)

	b.Handle(&changeUserLocButton, func(c *tb.Callback) {
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}
		if user.LastMessage != nil {
			b.Delete(user.LastMessage)
		}

		var err error
		user.SetState(EnterUserLoc)

		b.Send(c.Sender, "Выбери локацию", locSelectKbd)
		if err != nil {
			log.Println("Send messeage fail: ", err)
		}

		b.Respond(c, &tb.CallbackResponse{
			Text: "",
		})
		users.Update(id, user)
	})

	b.Handle(&backToMainMenuButton, func(c *tb.Callback) {
		id, user, notFound := GetUserFromCallbackByID(c, users, b)
		if notFound {
			return
		}
		if user.LastMessage != nil {
			b.Delete(user.LastMessage)
		}

		var err error
		user.SetState(MainMenu)

		b.Send(c.Sender, "Выбери действие", menuMain)
		if err != nil {
			log.Println("Send messeage fail: ", err)
		}
		users.Update(id, user)
	})

	b.Start()

}

func GetUserFromCallbackByID(c *tb.Callback, users UserStorager, b *tb.Bot) (UserID, User, bool) {
	id := UserID(c.Sender.ID)
	if c.Message != nil {
		b.Delete(c.Message)
	}
	user, err := users.Get(id)
	if err != nil {
		b.Send(c.Sender, "Случилась какая то ошибка. давай начнем заново. Жми /start")
		return 0, User{}, true
	}

	return id, user, false
}

// func initConfig() error {
// 	viper.AddConfigPath("configs")
// 	viper.SetConfigName("config.yml")

// 	return viper.ReadConfig()
// }

func LetUpgradeButtonHandler(c *tb.Callback) {
	id, user, notFound := GetUserFromCallbackByID(c, users, b)
	if notFound {
		return
	}
	user.SetState(ChooseForUpgrade)
	if c.Message != nil {
		b.Delete(c.Message)
	}

	holdesForUpgrade := make([]*Holde, 0)

	for _, h := range user.CurrPlayer.Request.Holdes {
		hol, err := holdeStorage.Get(h.HoldeID)
		if err != nil {
			log.Println(err)
		}
		if hol.Level < gameSettings.HoldeMaxLevel {
			holdesForUpgrade = append(holdesForUpgrade, hol)
		}
	}
	if len(holdesForUpgrade) == 0 {
		user.SetState(MainMenu)
		users.Update(id, user)
		b.Send(c.Sender, "Нет поместий для улучшения", menuMain)
		return
	}

	log.Println("Holdes  for upgrade", holdesForUpgrade)
	holdeSelKbd := &tb.ReplyMarkup{ResizeReplyKeyboard: true, OneTimeKeyboard: true}
	holdeSelectionButtons := make([]tb.Btn, len(holdesForUpgrade)+1)

	for i, holde := range holdesForUpgrade {

		holdeSelectionButtons[i] = holdeSelKbd.Data(
			fmt.Sprintf("%d %s Ур.%d", holde.ID, holde.Name, holde.Level), "holdesel_"+strconv.Itoa(holde.ID), strconv.Itoa(holde.ID))

		b.Handle(&holdeSelectionButtons[i], func(c *tb.Callback) {
			//upgrade kbd
			id, user, notFound := GetUserFromCallbackByID(c, users, b)
			if notFound {
				return
			}
			if c.Message != nil {
				b.Delete(c.Message)
			}
			holdeId, err := strconv.Atoi(c.Data)

			user.CurrHolde = holdeId

			holde, err := holdeStorage.Get(holdeId)
			log.Println(holde)
			if err != nil {
				b.Send(c.Sender, "Случилась какая-то ошибка. давай начнем заново. Жми /start")
				return
			}
			msg := fmt.Sprintf("Улучшение поместья стоит %s монет", gameSettings.HoldeLevelUpgradeCost[holde.Level-1].String())
			b.Send(c.Sender, msg, holdeUpgradeKbd)
			user.SetState(UpgradeHolde)
			users.Update(id, user)

		})

	}
	holdeSelectionButtons[len(holdesForUpgrade)] = holdeSelKbd.Data("Отмена", "upgrade_holde_cancel_button")

	holdeSelectionButtons[len(holdesForUpgrade)] = endCalculateButton
	holdeSelectButtonsRows := make([]tb.Row, len(holdesForUpgrade)+1)
	for i, _ := range holdeSelectionButtons {
		holdeSelectButtonsRows[i] = holdeSelKbd.Row(holdeSelectionButtons[i])
	}
	holdeSelKbd.Inline(holdeSelectButtonsRows...)
	users.Update(id, user)
	b.Send(c.Sender, "Выбери поместье для улучшения", holdeSelKbd)

}

func HoldeUpgradeApproveButtonHandler(c *tb.Callback) {
	id, user, notFound := GetUserFromCallbackByID(c, users, b)
	if notFound {
		return
	}
	user.SetState(ChooseForUpgrade)
	log.Println("User curr holde  ", user.CurrHolde)
	holde, err := holdeStorage.Get(user.CurrHolde)
	if err != nil {
		b.Send(c.Sender, "Случилась какая-то ошибка. давай начнем заново. Жми /start")
		return
	}

	holde.Upgrade()

	users.Update(id, user)

	holdeStorage.Update(holde)

	b.Send(c.Sender, fmt.Sprintf("Поместье %d %s улучшено до %d уровня", holde.ID, holde.Name, holde.Level))
	b.Send(c.Sender, "Выберите следующее поместье для улучшения", upgradeKeyboard)
}

func HoldeUpgradeCancelBtnHandler(c *tb.Callback) {
	id, user, notFound := GetUserFromCallbackByID(c, users, b)
	if notFound {
		return
	}
	user.SetState(UpgradeHolde)
	users.Update(id, user)
	b.Send(c.Sender, "Выбери опцию", upgradeKeyboard)

}

// endCalculateButton

func EndCalculateButtonHandler(c *tb.Callback) {
	id, user, notFound := GetUserFromCallbackByID(c, users, b)
	if notFound {
		return
	}

	b.Send(c.Sender, "До свидания")

	user.CurrPlayer = nil
	user.CurrHolde = -1

	user.SetState(UpgradeHolde)
	users.Update(id, user)
	b.Send(c.Sender, "Выбери опцию", menuMain)

}
