package main

type DialogState struct { 
	state int 
}
const (
	// userState 
	IDLE int=  iota
	UserSettings 
	Menu 
	HoldeCalc
)


type User  struct {
	State  DialogState 
}

type Users map [string] User 