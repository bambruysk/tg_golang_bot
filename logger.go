package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotLog struct {
	Chats []tb.ChatID
	Bot   *tb.Bot
}

//NewBotLog create ne l0gger instanse with deafult subscribers 
func NewBotLog(bot *tb.Bot) BotLog {
	chats := make([]tb.ChatID, 1)
	chats[0] = tb.ChatID(-1001241689222)
	return BotLog{
		Chats: chats,
		Bot:   bot,
	}
}

// AddSubscriber return true if subscriber added
func (b *BotLog) AddSubscriber(id tb.ChatID) bool {
	// repalce to index& 
	for _,c :=  range b.Chats {
		if c ==  id {
			return false
		}
	}
	b.Chats = append(b.Chats, id)
	return true
}

//Log send message to subscribers
func (b *BotLog) Log(message string) {
	for _, chat := range b.Chats {
		b.Bot.Send(chat, message)
	}
}
