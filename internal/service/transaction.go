package service

import (
	"fmt"
	"infotecs/internal/entity"
	"infotecs/internal/repository"
)

type WalletService interface {
	Transfer(senderAddress string, receiverAddress string, amount float64) error
	GetLastTransactions(count int) ([]entity.Transaction, error)
	GetBalanceByAddress(address string) (float64, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) Transfer(senderAddress string, receiverAddress string, amount float64) error {
	if senderAddress == "" || receiverAddress == "" {
		return fmt.Errorf("sender and recipient addresses cannot be empty")
	}

	if amount <= 0 {
		return fmt.Errorf("the transfer amount must be greater than zero")
	}

	if senderAddress == receiverAddress {
		return fmt.Errorf("the sender and recipient addresses cannot be the same")
	}

	err := s.repo.Transfer(senderAddress, receiverAddress, amount)
	if err != nil {
		return fmt.Errorf("error while transferring funds: %w", err)
	}

	return nil
}

func (s *walletService) GetLastTransactions(count int) ([]entity.Transaction, error) {
	return s.repo.GetLastTransactions(count)
}

func (s *walletService) GetBalanceByAddress(address string) (float64, error) {
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}

	balance, err := s.repo.GetBalanceByAddress(address)
	if err != nil {
		return 0, fmt.Errorf("failed to get balance for address %s: %w", address, err)
	}
	return balance, err
}
