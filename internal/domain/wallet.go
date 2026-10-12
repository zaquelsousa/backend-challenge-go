package domain

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)

type Wallet struct {
	id uuid.UUID
	playerID uuid.UUID
	currency string
	balance Money
	version int64
	createdAt time.Time
	updatedAt time.Time
}

func newWallet(playerID uuid.UUID, currency string, initialBalance Money) (*Wallet, error){
	if initialBalance.IsTheCurrency() != currency {
		return nil, fmt.Errorf("Currency mismatch")
	}

	if initialBalance.IsNegative(){
		return nil, fmt.Errorf("Invalid amount")
	}

	now := time.Now().UTC()
	return &Wallet{
		id:        uuid.New(),
		playerID:  playerID,
		currency:  currency,
		balance:   initialBalance,
		version:   1,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (w *Wallet) Debit(amount Money)error{
	if amount.IsNegative() || amount.IsZero() {
		return fmt.Errorf("debit amount must be positive")
	}

	if amount.IsTheCurrency() != w.currency {
		return fmt.Errorf("Currency mismatch")
	}

	newBalance, err := w.balance.Subtract(amount)
	if err != nil {
		return err
	}

	if newBalance.IsNegative(){
		return fmt.Errorf("insufficient funds")
	}

	w.commit(newBalance)
	return nil
}
func (w *Wallet) Credit(amount Money)error {
	if amount.IsNegative() || amount.IsZero() {
		return fmt.Errorf("credit amount must be positive")
	}
	if amount.IsTheCurrency() != w.currency {
		return fmt.Errorf("currency mismatch")
	}

	newBalance, err := w.balance.Add(amount)
	if err != nil {
		return err
	}

	w.commit(newBalance)
	return nil

}

func (w *Wallet) commit(newBalance Money) {
	w.balance = newBalance
	w.version++
	w.updatedAt = time.Now().UTC()
}

