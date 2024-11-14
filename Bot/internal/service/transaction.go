package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/justIGreK/MoneyKeeper/Bot/internal/models"
	"github.com/justIGreK/MoneyKeeper/Bot/pkg/util"
)

const (
	TimeFormat     string = "15:04"
	Dateformat     string = "2006-01-02"
	DateTimeformat string = "2006-01-02T15:04"
)

func (s *Service) AddTransaction(args []string, chatID string) (string, error) {
	if len(args) < 2 || len(args) > 5 {
		return "", fmt.Errorf("incorrect input format")
	}

	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	cost, err := strconv.ParseFloat(args[1], 32)
	if err != nil {
		return "", errors.New("invalid amount")
	}
	var category string
	var dateTime time.Time
	dateProvided := false
	timeProvided := false

	for _, part := range args[2:] {
		if util.IsDate(part) {
			if dateProvided {
				return "", fmt.Errorf("date is given several times")
			}
			dateTime, err = util.ParseDate(part)
			if err != nil {
				return "", fmt.Errorf("Incorrect date format")
			}
			dateProvided = true
		} else if util.IsTime(part) {
			if timeProvided {
				return "", fmt.Errorf("time is indicated several times")
			}
			dateTime, err = util.ParseTime(part, dateTime)
			if err != nil {
				return "", fmt.Errorf("Incorrect time format")
			}
			timeProvided = true
		} else if category == "" {
			category = part
		} else {
			return "", fmt.Errorf("Invalid request format: too many parameters")
		}
	}
	if dateTime.IsZero() {
		dateTime = time.Now()
	}
	date := dateTime.Format(DateTimeformat)
	createTx := models.CreateTransaction{
		UserID:   id,
		Name:     args[0],
		Cost:     float32(cost),
		Date:     &date,
		Category: category,
	}
	txid, err := s.tx.AddTransaction(s.ctx, createTx)
	if err != nil {
		return "", err
	}
	message := fmt.Sprintf("transaction is added %v", txid)

	return message, nil
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

func (s *Service) DeleteTransaction(args []string, chatID string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("incorrect input format")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	err = s.tx.DeleteTransaction(s.ctx, id, args[0])
	if err != nil {
		return "", err
	}
	
	return "Transaction is deleted", nil
}
func (s *Service) UpdateTransaction(args []string, chatID string) (string, error) {
	if len(args) < 2 || len(args) > 6 {
		return "", fmt.Errorf("incorrect input format")
	}
	id, err := s.GetUserID(chatID)
	if err != nil {
		return "", err
	}
	updates := make(map[string]interface{})
	categoryProvided := false
	costProvided := false
	nameProvided := false
	dateProvided := false
	timeProvided := false
	for _, part := range args[1:] {
		key, value, ok := util.ParseKeyValue(part)
		if !ok {
			return "", fmt.Errorf("incorrect format parameter: %s", part)
		}
		switch key {
		case "category":
			if categoryProvided {
				return "", fmt.Errorf("category is given several times")
			}
			updates["category"] = value
			categoryProvided = true
		case "name":
			if nameProvided {
				return "", fmt.Errorf("name is given several times")
			}
			updates["name"] = value
			nameProvided = true
		case "cost":
			if costProvided {
				return "", fmt.Errorf("cost is given several times")
			}
			cost, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return "", fmt.Errorf("incorrect cost: %s", value)
			}
			updates["cost"] = cost
			costProvided = true
		case "date":
			if dateProvided {
				return "", fmt.Errorf("date is given several times")
			}
			date, err := util.ParseDate(value)
			if err != nil {
				return "", fmt.Errorf("incorrect date format: %s", value)
			}
			updates["date"] = date
			dateProvided = true
		case "time":
			if timeProvided {
				return "", fmt.Errorf("time is given several times")
			}
			parsedTime, err := time.Parse("15:04", value)
			if err != nil {
				return "", fmt.Errorf("incorrect time format: %s", value)
			}
			updates["time"] = parsedTime
			timeProvided = true
		default:
			return "", fmt.Errorf("unknown parameter: %s", key)
		}
	}

	txUpdates := s.prepareTxForUpdate(updates)
	txUpdates.ID, txUpdates.UserID = args[0], id
	txs, err := s.tx.UpdateTransaction(s.ctx, txUpdates)
	if err != nil {
		return "", err
	}
	message := "Updated transaction:"
	txMessage := s.PrepareTransactions([]models.Transaction{*txs})
	message = fmt.Sprintf("%s%s", message, txMessage)
	return message, nil
}

func (s *Service) prepareTxForUpdate(updates map[string]interface{}) models.UpdateTransaction {
	var txUpdates models.UpdateTransaction
	for field, value := range updates {
		switch field {
		case "name":
			name := value.(string)
			txUpdates.Name = &name
		case "category":
			category := value.(string)
			txUpdates.Category = &category
		case "cost":
			cost := value.(float64)
			txUpdates.Cost = &cost
		case "time":
			times := value.(time.Time)
			timestr := times.Format(TimeFormat)
			txUpdates.Time = &timestr
		case "date":
			date := value.(time.Time)
			datestr := date.Format(Dateformat)
			txUpdates.Date = &datestr
		}
	}
	return txUpdates
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
