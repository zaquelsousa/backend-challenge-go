package httpapi

import (
	"backend-challange/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)


type moneyRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type createWalletRequest struct {
	PlayerID       string      `json:"playerId"`
	InitialBalance moneyRequest `json:"initialBalance"`
}

type moneyResponse struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type createWalletResponse struct {
	ID        string        `json:"id"`
	PlayerID  string        `json:"playerId"`
	Currency  string        `json:"currency"`
	Balance   moneyResponse `json:"balance"`
	Version   int64         `json:"version"`
}

//handler POST /wallets
func HandlerWallet(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	
	//to prevent the client to send large payload
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	var request createWalletRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	 if err := decoder.Decode(&request); err  != nil {
		http.Error(w, "invalid JSON request", http.StatusBadRequest)
		return
	 }

	 if strings.TrimSpace(request.PlayerID) == "" ||
		strings.TrimSpace(request.InitialBalance.Currency) == "" ||
		strings.TrimSpace(request.InitialBalance.Amount) == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	playerID, err := uuid.Parse(request.PlayerID)
	if err != nil {
		http.Error(w, "invalid playerId", http.StatusBadRequest)
		return
	}

	currency := request.InitialBalance.Currency

	inititalBalance, err := domain.ParseMoney(request.InitialBalance.Amount, currency)
	if err != nil {
		http.Error(w, "invalid initialBalance", http.StatusBadRequest)
		return
	}

	wallet, err := domain.NewWallet(playerID, currency, inititalBalance)
	if err != nil {
		http.Error(w, "could not create wallet", http.StatusBadRequest)
		return
	}

	balance := wallet.Balance()

	response := createWalletResponse{
		ID:       wallet.ID().String(),
		PlayerID: wallet.PlayerID().String(),
		Currency: wallet.Currency(),
		Balance: moneyResponse{
			Amount:   fmt.Sprintf("%d.%02d", balance.Amount/100, balance.Amount%100),
			Currency: balance.Currency,
		},
		Version: wallet.Version(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
