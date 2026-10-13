package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)


type LedgerDirection string

const (
	LedgerCredit LedgerDirection = "Credit"
	LedgerDebit LedgerDirection = "Debit"
)

type LedgerEntry struct {
	ID uuid.UUID
	WalletID uuid.UUID
	TransactionID uuid.UUID
	Direction LedgerDirection
	Money Money
	BalanceBefore Money
	BalanceAfter Money
	CreatedAt time.Time
}

func NewLedgerEntry(id uuid.UUID, walletID uuid.UUID, transactionID uuid.UUID,
	direction LedgerDirection, money Money, balanceBefore Money, balanceAfter Money, createdAt time.Time)(LedgerEntry, error){
	if id == uuid.Nil || walletID == uuid.Nil || transactionID == uuid.Nil {
		return LedgerEntry{}, fmt.Errorf("ledger IDs cannot be empty")
	}

	if direction != LedgerCredit && direction != LedgerDebit {
		return LedgerEntry{}, fmt.Errorf("invalid ledger direction")
	}

	if money.IsNegative() {
		return LedgerEntry{}, fmt.Errorf("ledger amount must be positive")
	}

	if money.Currency != balanceBefore.Currency || money.Currency != balanceAfter.Currency {
		return LedgerEntry{}, fmt.Errorf("ledger currencies must match")
	}

	var expectedBalance int64

	switch direction {
		case LedgerCredit:
			if balanceBefore.Amount > math.MaxInt64-money.Amount {
				return LedgerEntry{}, fmt.Errorf("invalid ledger balance transition")
			}
			expectedBalance = balanceBefore.Amount + money.Amount

		case LedgerDebit:
			if balanceBefore.Amount < money.Amount {
				return LedgerEntry{}, fmt.Errorf("invalid ledger balance transition")
			}
	}

	if balanceAfter.Amount != expectedBalance {
		return LedgerEntry{}, fmt.Errorf("invalid ledger balance transition")
	}

	if createdAt.IsZero() {
		return LedgerEntry{}, fmt.Errorf("ledger creation time cannot be zero")
	}

		
	return LedgerEntry{
		ID:            id,
		WalletID:      walletID,
		TransactionID: transactionID,
		Direction:     direction,
		Money:         money,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		CreatedAt:     createdAt,
	}, nil
}
