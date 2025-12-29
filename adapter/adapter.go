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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = a.Config.Timeout

	updates := a.bot.GetUpdatesChan(u)

	go func() {
		for upd := range updates {
			switch {
			case upd.Message.Text != "":
				// if upd.Message.Chat != nil {continue} // игнорируем сообщения из каналов
				msq := IncomingMessage{
					ChatID: upd.Message.Chat.ID,
					Text:   upd.Message.Text,
				}
				a.messageChan <- msq

			case upd.CallbackQuery != nil:

			}
		}
	}()

	return nil
}