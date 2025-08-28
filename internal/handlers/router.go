package handlers

import (
	"infotecs/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(walletService service.WalletService) *chi.Mux {
	r := chi.NewRouter()

	walletHandler := NewWalletHandler(walletService)

	r.Post("/api/send", walletHandler.Transfer)
	r.Get("/api/transactions", walletHandler.GetLastTransactions)
	r.Get("/api/wallet/{address}/balance", walletHandler.GetBalanceByAddress)

	return r
}
