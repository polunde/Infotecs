package entity

import (
	"time"
)

type Wallet struct {
	Address              string  `gorm:"primaryKey;unique;not null"`
	Balance              float64 `gorm:"not null;default:100"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	SentTransactions     []Transaction `gorm:"foreignKey:SenderAddress;references:Address"`
	ReceivedTransactions []Transaction `gorm:"foreignKey:ReceiverAddress;references:Address"`
}

type Transaction struct {
	ID              uint    `gorm:"primaryKey"`
	SenderAddress   string  `gorm:"not null"`
	ReceiverAddress string  `gorm:"not null"`
	Amount          float64 `gorm:"not null"`
	CreatedAt       time.Time

	Sender   Wallet `gorm:"foreignKey:SenderAddress;references:Address"`
	Receiver Wallet `gorm:"foreignKey:ReceiverAddress;references:Address"`
}
