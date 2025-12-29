package service

import (
	"fmt"
	"log"

	"github.com/Demon2050502/TG_Bot/adapter"
)

type TgAPIAdapter interface {
	GetMessageChannel() <-chan adapter.IncomingMessage
	GetCommandChannel() <-chan adapter.IncomingCommand
	SendTaskList(adapter.SendTaskListCommand) error
}

type Service struct {
	adapter	TgAPIAdapter
}

func NewService(a TgAPIAdapter) *Service {
	return &Service{
		adapter: a,
	}
}

func (s *Service) StartServe() error {
	log.Println("Service started serving!")

	for {
		select {
		case msg := <-s.adapter.GetMessageChannel():
			// Обработка входящего сообщения
			fmt.Println("Received message:", msg.Text)
			

			s.adapter.SendTaskList(adapter.SendTaskListCommand{Text: msg.Text, ChatID: msg.ChatID})

		case cmd := <-s.adapter.GetCommandChannel():
			// Обработка входящей команды
			fmt.Println("Received command:", cmd)
		}
	}
	
	return nil
}

