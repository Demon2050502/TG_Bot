package adapter

import (
	"log"

	"github.com/Demon2050502/TG_Bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type IncomingMessage struct {
	ChatID int64
	Text string
}

type IncomingCommand struct {
	ChatID int64
	Name   string 
	Args   string
}

// type SendMessageCommand struct {
// }

type SendTaskListCommand struct {
	Text string
	ChatID int64
}

type Adapter struct {
	Config config.BotConfig
	bot *tgbotapi.BotAPI
	messageChan chan IncomingMessage
	commandChan chan IncomingCommand
}

func (a *Adapter)GetMessageChannel() <-chan IncomingMessage {
	return a.messageChan
}

func (a *Adapter)GetCommandChannel() <-chan IncomingCommand {
	return a.commandChan
}

func (a *Adapter) SendTaskList(c SendTaskListCommand) error {

	msg := tgbotapi.NewMessage(c.ChatID, c.Text)
	if _, err := a.bot.Send(msg); err != nil {
		log.Println(err)
		return err
	}

	return nil
}


func NewAdapter(cfg config.BotConfig) *Adapter {
	a := &Adapter{
		Config: cfg,
	}

	a.messageChan = make(chan IncomingMessage)
	a.commandChan = make(chan IncomingCommand)

	if err := a.startAdapter(); err != nil {
		panic(err)
	}

	return a
}

func (a *Adapter) startAdapter() error{
	var err error
	a.bot, err = tgbotapi.NewBotAPI(a.Config.APIKei)
	if err != nil {
		return err
	}

	commands := []tgbotapi.BotCommand{
	{Command: "start", Description: "Запуск и помощь"},
	{Command: "schema", Description: "Расписание релизов на неделю"},
}

	cfg := tgbotapi.NewSetMyCommands(commands...)
	if _, err := a.bot.Request(cfg); err != nil {
		log.Println("Не удалось установить команды:", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = a.Config.Timeout

	updates := a.bot.GetUpdatesChan(u)

	go func() {
		for upd := range updates {
			if upd.Message == nil || upd.Message.Chat == nil {
				continue
			}

			// Команды (/schema)
			if upd.Message.IsCommand() {
				a.commandChan <- IncomingCommand{
					ChatID: upd.Message.Chat.ID,
					Name:   upd.Message.Command(),
					Args:   upd.Message.CommandArguments(),
				}
				continue
			}

			// Обычный текст
			if upd.Message.Text != "" {
				a.messageChan <- IncomingMessage{
					ChatID: upd.Message.Chat.ID,
					Text:   upd.Message.Text,
				}
			}
		}
	}()

	return nil
}