package utils

import (
	"infotecs/internal/entity"

	"github.com/google/uuid"
)

func GenerateWallets(n int, initialBalance float64) []entity.Wallet {
	wallets := make([]entity.Wallet, n)
	for i := 0; i < n; i++ {
		wallets[i] = entity.Wallet{
			Address: uuid.NewString(),
			Balance: initialBalance,
		}
	}
	return wallets
}
