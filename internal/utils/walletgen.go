package utils

import (
	"bufio"
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

func LoadWalletsFromFile(filename string, initialBalance float64) []entity.Wallet {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	var wallets []entity.Wallet
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		address := scanner.Text()
		if address != "" {
			wallets = append(wallets, entity.Wallet{
				Address: address,
				Balance: initialBalance,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("reading file error: %v", err)
	}

	return wallets

}
