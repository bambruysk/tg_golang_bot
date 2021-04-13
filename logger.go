package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotLog struct {
	Chats []tb.ChatID
	Bot   *tb.Bot
}

func NewBotLog(bot *tb.Bot) BotLog {
	chats := make([]tb.ChatID, 1)
	chats[0] = tb.ChatID(342270809)

	return BotLog{
		Chats: chats,
		Bot:   bot,
	}
}

func (b *BotLog) AddSubscriber(id tb.ChatID) {
	b.Chats = append(b.Chats, id)
}

func (b *BotLog) Log(message string) {
	for _, chat := range b.Chats {
		b.Bot.Send(chat, message)
	}
}
