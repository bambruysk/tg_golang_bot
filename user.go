package main

import "errors"

type DialogState int

const (
	// userState
	IDLE DialogState = iota
	UserSettings
	MainMenu
	HoldeCalc
	AddHolde
	EnterDice
	EnterPlayerName
	EnterUserName
	EnterUserLoc
	
)

// Пользователь
type User struct {
	State      DialogState
	CurrHolde  int
	CurrPlayer *Player
	Name       string
	Location   string
}

func (u *User) SetState(state DialogState) {
	u.State = state
}

func NewUser() User {
	return User{
		State:     IDLE,
		CurrHolde: -1,
	}
}

func (u *User) Save() {
	//TODO: add user  save
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
