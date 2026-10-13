package dto

import (
	"github.com/google/uuid"
)

type MoneyResponse struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type WalletResponse struct {
	ID       uuid.UUID     `json:"id"`
	PlayerID uuid.UUID     `json:"playerId"`
	Balance  MoneyResponse `json:"balance"`
	Version  int64         `json:"version"`
}


