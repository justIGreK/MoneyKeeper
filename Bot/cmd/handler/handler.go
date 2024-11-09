package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/service"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type TelegramHandler struct {
	bot     *telego.Bot
	service *service.Service
}

type CommandHandler func(args []string, chatID string) (string, error)

var (
	commandMapInstance map[string]CommandHandler
	once               sync.Once
)

func NewTelegramHandler(botToken string, srv *service.Service) *TelegramHandler {
	bot, err := telego.NewBot(botToken)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%v/setWebhook?url=%v/webhook", os.Getenv("BOT_TOKEN"), os.Getenv("WEBHOOK_URL"))
	err = bot.SetWebhook(&telego.SetWebhookParams{
		URL: url,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &TelegramHandler{bot: bot, service: srv}
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
		var message telego.SendMessageParams
		chatID := tu.ID(update.Message.Chat.ID)
		log.Printf("Received message: %s", update.Message.Text)
		text, err := h.handleMessage(update.Message.Text, update.Message.Chat.ChatID().String())
		if err != nil {
			message.Text = err.Error()
		} else {
			message.Text = text
		}
		message.ChatID = chatID
		_, err = h.bot.SendMessage(&message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

func (h *TelegramHandler) handleMessage(message, chatID string) (string, error) {
	h.CommandMap()
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty message")
	}

	command := parts[0]
	args := parts[1:]

	handler, exists := commandMapInstance[command]
	if !exists {
		return "", fmt.Errorf("unknown command: %s", command)
	}
	return handler(args, chatID)
}

func (h *TelegramHandler) CommandMap() map[string]CommandHandler {
	once.Do(func() {
		commandMapInstance = map[string]CommandHandler{
			"/addBudget": h.service.AddBudget,
			"/getBudgets": h.service.GetBudgetList,
			"/addTx": h.service.AddTransaction,
			"/getTx": h.service.GetTransaction,
			"/getTxList": h.service.GetTransactionList,
			"/getTxsByDates": h.service.GetTXByTimeFrame,
			"/getSummary": h.service.GetSummary,
			"/getBudgetReport": h.service.GetBudgetReport,
		}
	})
	return commandMapInstance
}
