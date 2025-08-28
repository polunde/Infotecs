package repository

import (
	"fmt"
	"infotecs/internal/entity"

	"gorm.io/gorm"
)

type WalletRepository interface {
	Transfer(senderAddress string, receiverAddress string, amount float64) error
	GetLastTransactions(count int) ([]entity.Transaction, error)
	GetBalanceByAddress(address string) (float64, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Transfer(senderAddress string, receiverAddress string, amount float64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var sender, receiver entity.Wallet

		if err := tx.Where("address = ?", senderAddress).First(&sender).Error; err != nil {
			return fmt.Errorf("could not find the sender: %w", err)
		}
		if err := tx.Where("address = ?", receiverAddress).First(&receiver).Error; err != nil {
			return fmt.Errorf("could not find the receiver: %w", err)
		}

		if sender.Balance < amount {
			return fmt.Errorf("not enough funds")
		}

		if err := tx.Model(&sender).Update("balance", sender.Balance-amount).Error; err != nil {
			return err
		}
		if err := tx.Model(&receiver).Update("balance", receiver.Balance+amount).Error; err != nil {
			return err
		}

		t := entity.Transaction{
			SenderAddress:   sender.Address,
			ReceiverAddress: receiver.Address,
			Amount:          amount,
		}

		if err := tx.Create(&t).Error; err != nil {
			return err
		}

		return nil

	})
}

func (r *walletRepository) GetLastTransactions(count int) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Preload("Sender").Preload("Receiver").Order("created_at desc").Limit(count).Find(&transactions).Error
	return transactions, err
}

func (r *walletRepository) GetBalanceByAddress(address string) (float64, error) {
	var wallet entity.Wallet
	err := r.db.Select("balance").Where("address = ?", address).First(&wallet).Error
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}
