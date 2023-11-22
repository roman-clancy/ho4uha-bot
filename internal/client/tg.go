package client

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/roman-clancy/ho4uha-bot/internal/model/messages"
	"github.com/roman-clancy/ho4uha-bot/internal/model/types"
)

type TgClient struct {
	client      *tgbotapi.BotAPI
	handlerFunc HandlerFunc
}

type HandlerFunc func(update tgbotapi.Update, c *TgClient, m *messages.BotModel)

func New(token string, handlerFunc HandlerFunc) (*TgClient, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TgClient{
		client:      client,
		handlerFunc: handlerFunc,
	}, nil
}

func (c *TgClient) SendMessage(userId int64, text string) error {
	message := tgbotapi.NewMessage(userId, text)
	_, err := c.client.Send(message)
	return err
}

func (c *TgClient) ListenUpdates(botModel *messages.BotModel) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updatesChan := c.client.GetUpdatesChan(u)
	for update := range updatesChan {
		c.handlerFunc(update, c, botModel)
	}
}

func (c *TgClient) ShowButtons(userId int64, text string, buttons []types.TgRowButtons) error {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, len(buttons))
	for i := 0; i < len(buttons); i++ {
		tgRowButtons := buttons[i]
		keyboard[i] = make([]tgbotapi.InlineKeyboardButton, len(tgRowButtons))
		for j := 0; j < len(tgRowButtons); j++ {
			tgInlineButton := tgRowButtons[j]
			keyboard[i][j] = tgbotapi.NewInlineKeyboardButtonData(tgInlineButton.DisplayName, tgInlineButton.Value)
		}
	}
	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	msg := tgbotapi.NewMessage(userId, text)
	msg.ReplyMarkup = numericKeyboard
	msg.ParseMode = "markdown"
	_, err := c.client.Send(msg)
	return err
}

func ProcessingMessage(update tgbotapi.Update, client *TgClient, botModel *messages.BotModel) {
	if update.Message != nil {
		err := botModel.OnMessage(messages.Message{
			Text:     update.Message.Text,
			UserID:   update.Message.From.ID,
			UserName: update.Message.From.UserName,
		})
		if err != nil {
			return
		}
	} else if update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
		if _, err := client.client.Request(callback); err != nil {
			return
		}
		if err := deleteInlineButtons(client, update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Message.Text); err != nil {
			return
		}
		err := botModel.OnMessage(messages.Message{
			Text:          update.CallbackQuery.Data,
			UserID:        update.CallbackQuery.From.ID,
			UserName:      update.CallbackQuery.From.UserName,
			IsCallback:    true,
			CallbackMsgID: update.CallbackQuery.ID,
		})
		if err != nil {
			return
		}
	}
}

func deleteInlineButtons(c *TgClient, userID int64, msgID int, sourceText string) error {
	msg := tgbotapi.NewEditMessageText(userID, msgID, sourceText)
	_, err := c.client.Send(msg)
	return err
}
