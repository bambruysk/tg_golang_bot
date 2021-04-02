package main

import (
	"errors"

	tb "gopkg.in/tucnak/telebot.v2"
)

type DialogState DialogNode

/*
const (
	// userState
	IDLE DialogState = iota
	UserSettings
	MainMenu
	HoldeCalc
)

*/
// Пользователь
type User struct {
	State DialogState
}

func (u *User) SetState(state DialogState) {
	u.State = state
}

func NewUser() User {
	return User{
		State: CreateDeaultState(),
	}
}




/// Хранилище пользователей, потом заменим на BD

type UserID int

type UserStorager interface {
	// Get user
	Get(id UserID) (User, error)
	// Geto or create user w default state
	GetOrCreate(id UserID) (User, bool)

	Update(id UserID, user User) error

	Create(id UserID, user User)
}

type Users map[UserID]User

func NewUsers() Users {
	users := make(Users)
	return users
}

func (u Users) Get(id UserID) (User, error) {
	user, exist := u[id]
	if !exist {
		return User{}, errors.New("User not found")
	}
	return user, nil
}

func (u Users) GetOrCreate(id UserID) (User, bool) {
	user, exist := u[id]
	if !exist {
		user = NewUser()
		u.Create(id, user)
		return user, true
	}
	return user, false
}

func (u Users) Create(id UserID, user User) {
	u[id] = user
}

func (u Users) Update(id UserID, user User) error {
	_, exist := u[id]
	if !exist {
		return errors.New("User not found ")
	}
	u[id] = user
	return nil
}

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
func (dn DialogNode)  Send(*tg.Bot, tg.Recipient, *tg.SendOptions) (*tg.Message, error) {
	// not implemented yet

	return &tg.Message{}, errors.New("Not implemented yet")

}


func (dn DialogNode) GetHandler () interface{} {
	return func 


}

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
)

var users UserStorager

func CreateDialog(bot *tb.Bot, enterPoint interface{}) {

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
	*bot.Handle(enterPoint, func(m *tb.Message) {
		
		user_id := UserID(m.Sender.ID)

		user, created := users.GetOrCreate(user_id)

		*bot.Send(m.Sender, start_dn.DialogContent.Message, start_dn.Keyboard)
	} )

	// Setiings
	
	*bot.Handle(&btnSettings, func(m *tb.Message) {
		
		id := UserID(m.Sender.ID)
		user, err := users.Get(id)
		if 
		user.SetState(settings_dn)
		users.Update(id, user)

		

		*bot.Send(m.Sender, start_dn.DialogContent.Message, start_dn.Keyboard)
	} )



}

func CreateDeaultState() DialogState {
	return DialogState(NewDialogNode(DialogContent{
		Message: "Добро пожаловать",
		Media:   nil,
	}, btnSettings, btnCalculator))

}
