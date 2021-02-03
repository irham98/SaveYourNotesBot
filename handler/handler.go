package handler

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/irham/agung/notes-bot/gdrive"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type (
	handler struct {
		bot *tgbotapi.BotAPI
		log *logrus.Logger
	}
)

func New(bot *tgbotapi.BotAPI, log *logrus.Logger) *handler {
	return &handler{
		bot: bot,
		log: log,
	}
}

func (h *handler) FindImage() {
}

func (h *handler) MessageStart(update *tgbotapi.Update) {
	_, err := h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf(MessageStartResponse, update.Message.Chat.FirstName)))
	if err != nil {
		h.log.Errorf("MessageStart error: %v", err)
	}
}

func (h *handler) MessageAuth(config *oauth2.Config, update *tgbotapi.Update) {
	_, err := gdrive.FindToken(config)
	if err != nil {
		authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, MessageMissingAuthResponse)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonURL(AuthenticationUrlString, authURL),
			),
		)

		_, err := h.bot.Send(msg)
		if err != nil {
			h.log.Errorf("MessageAuth error: %v", err)
			return
		}
	}

}

func (h *handler) MessageUnknown(update *tgbotapi.Update) {
	_, err := h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, MessageUnknownResponse))
	if err != nil {
		h.log.Errorf("MessageUnknown error: %v", err)
	}
}
