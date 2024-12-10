package methodsForTransaction

import (
	"cachManagerApp/app/db/models"
	"cachManagerApp/app/pkg/logger"
	"cachManagerApp/database"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

func (transactions *TransactionsMethod) PostExpense(update tgbotapi.Update, category string) error {
	log := logger.GetLogger()

	msg := strings.Split(update.Message.Text, ", ")

	sum, err := strconv.Atoi(msg[0])
	if err != nil {
		log.Error("Ошибка преобразования суммы: %v", err)
		return err
	}

	transaction := models.Transactions{
		TelegramID:    uint64(update.Message.Chat.ID),
		CreatedAt:     time.Now(),
		OperationType: false,
		Quantities:    uint64(sum),
		Category:      category,
		Description:   msg[1],
	}
	var transactionExist models.Transactions
	res := database.DB.Where("telegram_id = ? AND created_at = ?", transaction.TelegramID, transaction.CreatedAt).First(&transactionExist).Error
	if res == nil {
		log.Debug("Транзакция существует")
		return errors.New("transaction already exists")
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		log.Error("Ошибка добавления новой транзакции: %v", err)
		return err

	}
	log.Println("Транзакция успешно добавлена")
	return nil
}