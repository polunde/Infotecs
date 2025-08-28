package utils

import (
	"infotecs/internal/entity"
	"log"
	"os"

	"github.com/google/uuid"
)

func GenerateWallets(n int, initialBalance float64) []entity.Wallet {
	wallets := make([]entity.Wallet, n)
	file, err := os.Create("wallets_addresses.txt")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer file.Close()

	for i := 0; i < n; i++ {
		address := uuid.NewString()
		wallets[i] = entity.Wallet{
			Address: address,
			Balance: initialBalance,
		}

		if _, err := file.WriteString(address + "\n"); err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}
	}
	return wallets
}
