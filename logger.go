package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotLog struct {
	channels [] tb.Chat 
	users [] tb.Message 
}

func (b * BotLog) AddSubscriber {
	
}