package handlers

import (
	"encoding/json"
	"infotecs/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type TransferRequest struct {
	SenderAddress   string  `json:"from"`
	ReceiverAddress string  `json:"to"`
	Amount          float64 `json:"amount"`
}

type TransferResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type BalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type WalletHandler struct {
	walletService service.WalletService
}

func NewWalletHandler(s service.WalletService) *WalletHandler {
	return &WalletHandler{walletService: s}
}

func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.walletService.Transfer(req.SenderAddress, req.ReceiverAddress, req.Amount)
	if err != nil {
		log.Printf("transfer error: %v", err)
		resp := TransferResponse{Success: false, Error: err.Error()}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := TransferResponse{Success: true}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *WalletHandler) GetLastTransactions(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		count = 10
	}

	transactions, err := h.walletService.GetLastTransactions(count)
	if err != nil {
		log.Printf("error getting transactions: %v", err)
		http.Error(w, "failed to get transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("error encoding transactions: %v", err)
		http.Error(w, "failed to encode transactions", http.StatusInternalServerError)
	}
}

func (h *WalletHandler) GetBalanceByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	if address == "" {
		http.Error(w, "address parameter is required", http.StatusBadRequest)
		return
	}

	balance, err := h.walletService.GetBalanceByAddress(address)
	if err != nil {
		log.Printf("failed to get balance for %s: %v", address, err)
		http.Error(w, "could not get wallet balance", http.StatusInternalServerError)
		return
	}

	resp := BalanceResponse{
		Address: address,
		Balance: balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
