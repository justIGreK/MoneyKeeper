package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TelegramHandler struct {
	bot *telego.Bot
}

func NewTelegramHandler(botToken string) *TelegramHandler {
	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%v/setWebhook?url=%v/webhook", os.Getenv("BOT_TOKEN"), os.Getenv("WEBHOOK_URL"))
	bot.SetWebhook(&telego.SetWebhookParams{
		URL: url,
	})
	return &TelegramHandler{bot: bot}
}

func (h *TelegramHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var update telego.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("Error decoding webhook update: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	h.processUpdate(update)

	w.WriteHeader(http.StatusOK)
}

func (h *TelegramHandler) processUpdate(update telego.Update) {
	if update.Message != nil {
		log.Printf("Received message: %s", update.Message.Text)
		chatID := tu.ID(update.Message.Chat.ID)
		replyText := update.Message.Text

		message := &telego.SendMessageParams{
			Text:   replyText,
			ChatID: chatID,
		}
		_, err := h.bot.SendMessage(message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
