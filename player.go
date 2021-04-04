package main

import "errors"

type Holde struct {
	Name string
	ID   int
}

type Player struct {
	Name   string
	Holdes []Holde
}

func NewPlayer(name string) Player {
	return Player{Name: name}
}

type PlayerID string

type PlayerStorage map[PlayerID]Player

func NewPlayerStorage() PlayerStorage {
	//TODO : add load from database
	players := make(PlayerStorage)
	return players
}

func (ps PlayerStorage) Get(id PlayerID) (Player, bool) {
	player, exist := ps[id]
	return player, exist
}

func (ps *PlayerStorage) GetOrCreate(id PlayerID, name string) (Player, bool) {
	player, exist := (*ps)[id]
	if !exist {
		player = NewPlayer(name)
		(*ps)[id] = player
	}
	return player, exist
}
func (ps *PlayerStorage) Update(id PlayerID, holdes []Holde) error {
	player, exist := (*ps)[id]
	if !exist {
		return errors.New("Player not found")
	}
	player.Holdes = holdes
	return nil
}
