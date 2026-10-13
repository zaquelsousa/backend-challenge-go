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

func NewWallet(playerID uuid.UUID, currency string) (*Wallet, error){
	if len(currency) != 3 {
		return nil, fmt.Errorf("invalid currency")
	}

	now := time.Now().UTC()
	return &Wallet{
		id:        uuid.New(),
		playerID:  playerID,
		currency:  currency,
		balance:   Money{Amount: 0, Currency: currency},
		version:   1,
		createdAt: now,
		updatedAt: now,
	}, nil
}

//getters
func (w *Wallet) ID() uuid.UUID { return w.id }
func (w *Wallet) PlayerID() uuid.UUID { return w.playerID }
func (w *Wallet) Currency() string { return w.currency }
func (w *Wallet) Balance() Money { return w.balance }
func (w *Wallet) Version() int64 { return w.version }
func (w *Wallet) CreatedAt() time.Time { return w.createdAt }
func (w *Wallet) UpdateAt() time.Time { return w.updatedAt }


func (w *Wallet) Debit(amount Money, transactionID uuid.UUID)( LedgerEntry, error){
	if amount.IsNegative() || amount.IsZero() {
		return LedgerEntry{}, fmt.Errorf("debit amount must be positive")
	}

	if amount.IsTheCurrency() != w.currency {
		return LedgerEntry{}, fmt.Errorf("Currency mismatch")
	}

	newBalance, err := w.balance.Subtract(amount)
	if err != nil {
		return LedgerEntry{}, err
	}

	if newBalance.IsNegative(){
		return LedgerEntry{}, fmt.Errorf("insufficient funds")
	}
	
	//transctionID jaja eu crio a transacao por hora so para testes locais
	entry, err := NewLedgerEntry(w.id, transactionID, LedgerDebit, amount, w.balance, newBalance)

	w.commit(newBalance)
	return entry, nil
}
func (w *Wallet) Credit(amount Money, transactionID uuid.UUID)(LedgerEntry, error) {
	if amount.IsNegative() || amount.IsZero() {
		return LedgerEntry{}, fmt.Errorf("debit amount must be positive")

	}
	if amount.IsTheCurrency() != w.currency {
		return LedgerEntry{}, fmt.Errorf("Currency mismatch")

	}

	newBalance, err := w.balance.Add(amount)
	if err != nil {
		return LedgerEntry{}, fmt.Errorf("insufficient funds")
	}

	//transctionID jaja eu crio a transacao por hora so para testes locais
	entry, err := NewLedgerEntry(w.id, transactionID, LedgerCredit, amount, w.balance, newBalance)

	w.commit(newBalance)
	return entry ,nil

}

func (w *Wallet) commit(newBalance Money) {
	w.balance = newBalance
	w.version++
	w.updatedAt = time.Now().UTC()
}

