package service

import (
	"context"
	"fmt"
	"log"
	"time"

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

	loc, _ := time.LoadLocation("Europe/Moscow")
	if loc == nil {
		loc = time.Local
	}

	for {
		select {
			case msg := <-s.adapter.GetMessageChannel():
				// Обработка входящего сообщения
				fmt.Println("Received message:", msg.Text)
				s.adapter.SendTaskList(adapter.SendTaskListCommand{Text: msg.Text, ChatID: msg.ChatID})

			case cmd := <-s.adapter.GetCommandChannel():
				// Обработка входящей команды
				fmt.Println("Received command:", cmd)

				switch cmd.Name {
				case "start":
					_ = s.adapter.SendTaskList(adapter.SendTaskListCommand{
						ChatID: cmd.ChatID,
						Text:   "Привет! Команды:\n/schema — расписание релизов на неделю",
					})

				case "schema":
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					text, err := BuildSchemaText(ctx, loc)
					cancel()

					if err != nil {
						_ = s.adapter.SendTaskList(adapter.SendTaskListCommand{
							ChatID: cmd.ChatID,
							Text:   "Ошибка получения расписания: " + err.Error(),
						})
						continue
					}

					for _, part := range splitByLimit(text, 3800) {
						_ = s.adapter.SendTaskList(adapter.SendTaskListCommand{
							ChatID: cmd.ChatID,
							Text:   part,
						})
					}

				default:
					_ = s.adapter.SendTaskList(adapter.SendTaskListCommand{
						ChatID: cmd.ChatID,
						Text:   "Неизвестная команда: /" + cmd.Name,
					})
			}
		}
	}
	
	return nil
}

