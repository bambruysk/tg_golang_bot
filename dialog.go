package main

import (
	"errors"

	tb "gopkg.in/tucnak/telebot.v2"
)

type DialogContent struct {
	Message string
	Media   interface{}
}

type Handler func(interface{})

// Каждая нода диалога является состоянием пользвателя
type DialogNode struct {
	// Текст ии дугое сообщение оторжаео при переходе на данный узел
	Content DialogContent
	// Клавиатура
	Keyboard *tb.ReplyMarkup
	// Принимает текст в качетсве
	IsReceiveText bool
	//
	HasKeyboard bool
}

func NewDialogNode(content DialogContent, buttons ...tb.Btn) DialogNode {
	keyboard := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

	// для простоты будет выстраовиать все кнопки в столбик, то есть каждая кнопка
	rows := make([]tb.Row, len(buttons))
	for i, btn := range buttons {
		rows[i] = keyboard.Row(btn)
	}
	keyboard.Reply(rows...)

	return DialogNode{
		Content:  content,
		Keyboard: keyboard,
	}
}

/// Implement sendable inteface of telebot
func (dn DialogNode) Send(*tb.Bot, tb.Recipient, *tb.SendOptions) (*tb.Message, error) {
	// not implemented yet

	return &tb.Message{}, errors.New("Not implemented yet")

}

// func (dn DialogNode) GetHandler () interface{} {
// 	return func
// }

// Next compute next DialogNode from update for user
func (dn *DialogNode) Next(update interface{}) DialogNode {
	// switch tu:= update.(type) {
	// case *tb.Message: {

	// 	tu.
	// }
	// }
	return DialogNode{}
}

// my dialog buttons

var (
	btnSettings   = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("Настройки")
	btnCalculator = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("Поместья")
	btnMainMenu   = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("Главное меню")
	btnAddHolde   = (&tb.ReplyMarkup{ResizeReplyKeyboard: true}).Text("Добавить поместье")
)

var users UserStorager

/*
func TextHandler(m *tb.Message) {
	user_id := UserID(m.Sender.ID)
	user, exist := users.Get(user_id)
	text := m.Text
	if !exist {
		return
	}
	switch user.State {
	case PlayerName:
		{
			pl_name := text

		}
	case AddHolde:
		holde_num, err := strconv.Atoi(m.Text)
		if err != nil {
			//  send er tor bot
		}
		// add holde to hold request

	}

}

func CreateDialog(bot *tb.Bot, enterPoint interface{}) {
	/*
	   	start_dn := NewDialogNode(DialogContent{
	   		Message: "Добро пожаловать",
	   		Media:   nil,
	   	}, btnSettings, btnCalculator)

	   	start_dn.HasKeyboard = true

	   	settings_dn := NewDialogNode(DialogContent{
	   		Message: "Выберите настройку",
	   		Media:   nil,
	   	}, btnSettings, btnCalculator, btnMainMenu)

	   	// Add starting point
	   /*	*bot.Handle(enterPoint, func(m *tb.Message) {

	   		user_id := UserID(m.Sender.ID)

	   		user, created := users.GetOrCreate(user_id)

	   		*bot.Send(m.Sender, start_dn.DialogContent.Message, start_dn.Keyboard)
	   	})
*/
// Setiings
/*
		*bot.Handle(&btnSettings, func(m *tb.Message) {

			id := UserID(m.Sender.ID)
			user, err := users.Get(id)
			//if
			user.SetState(settings_dn)
			users.Update(id, user)

			*bot.Send(m.Sender, start_dn.DialogContent.Message, start_dn.Keyboard)
		})
	}
*/

func CreateDeaultState() DialogState {
	return IDLE
}
