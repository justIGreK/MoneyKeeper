package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
)

func (s *Service) AddTransaction(args []string, chatID string) (string, error) {
	if len(args) < 2 {
		return "", errors.New("missing arguments")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	cost, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return "", errors.New("invalid amount")
	}
	createTx := models.CreateTransaction{
		UserID: id,
		Name:   args[0],
		Cost:   float32(cost),
	}
	if len(args) >= 3 {
		createTx.Category = args[2]
	}
	notification, err := s.tx.AddTransaction(s.ctx, createTx)
	if err != nil {
		return "", err
	}

	message := "transaction is added"
	if len(notification) != 0 {
		notifications := s.createNotificationsMessage(notification)
		message = fmt.Sprintf("%s%s", message, notifications)
	}
	return message, nil
}

func (s *Service) createNotificationsMessage(notifications []string) string {
	var sb strings.Builder
	sb.WriteString(", but:\n")

	for _, notification := range notifications {
		sb.WriteString(fmt.Sprintf("%s\n", notification))
	}

	return sb.String()
}

func (s *Service) GetTransaction(args []string, chatID string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("missing arguments")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}

	tx, err := s.tx.GetTransaction(s.ctx, id, args[0])
	if err != nil {
		return "", err
	}
	message := "Transaction:"
	txMessage := s.PrepareTransactions([]models.Transaction{*tx})
	message = fmt.Sprintf("%s%s", message, txMessage)
	return message, nil
}

func (s *Service) PrepareTransactions(txs []models.Transaction) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for _, tx := range txs {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("ID: %s\n", tx.ID))
		sb.WriteString(fmt.Sprintf("Name: %s\n", tx.Name))
		sb.WriteString(fmt.Sprintf("Category: %s\n", tx.Category))
		sb.WriteString(fmt.Sprintf("Cost: %.2f\n", tx.Cost))
		sb.WriteString(fmt.Sprintf("Date: %s\n", tx.Date))
	}
	return sb.String()
}

func (s *Service) GetTransactionList(args []string, chatID string) (string, error) {
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}

	txs, err := s.tx.GetTransactionList(s.ctx, id)
	if err != nil {
		return "", err
	}
	message := "Transactions:"
	txMessage := s.PrepareTransactions(txs)
	message = fmt.Sprintf("%s%s", message, txMessage)
	return message, nil
}

func (s *Service) GetTXByTimeFrame(args []string, chatID string) (string, error) {
	start, end, err := s.parseTimeFrameOption(args)
	if err != nil {
		return "", err
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	txs, err := s.tx.GetTXByTimeFrame(s.ctx, id, start, end)
	if err != nil {
		return "", err
	}
	message := "Transactions:"
	txMessage := s.PrepareTransactions(txs)
	message = fmt.Sprintf("%s%s", message, txMessage)
	return message, nil
}

const (
	startPeriod = "from"
	endPeriod   = "to"
)

func (s *Service) parseTimeFrameOption(args []string) (string, string, error) {
	start, end := "", ""
	length := len(args)
	switch length {
	case 2:
		if args[0] == startPeriod {
			_, err := time.Parse(Dateformat, args[1])
			if err != nil {
				return "", "", err
			}
			start = args[1]
		} else if args[0] == endPeriod {
			_, err := time.Parse(Dateformat, args[1])
			if err != nil {
				return "", "", err
			}
			end = args[1]
		} else {
			return "", "", errors.New("invalid timeframe")
		}
	case 4:
		if args[0] == startPeriod && args[2] == endPeriod {
			_, err := time.Parse(Dateformat, args[1])
			if err != nil {
				return "", "", err
			}
			start = args[1]
			_, err = time.Parse(Dateformat, args[3])
			if err != nil {
				return "", "", err
			}
			end = args[3]
		} else {
			return "", "", errors.New("invalid timeframe")
		}
	default:
		return "", "", errors.New("invalid timeframe")
	}
	return start, end, nil
}
