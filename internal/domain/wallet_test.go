package domain

import (
	"testing"

	"github.com/google/uuid"
)

func newTestWallet(t *testing.T, amount int64, currency string) *Wallet {
	t.Helper()

	initialBalance, err := NewMoney(amount, currency)
	if err != nil {
		t.Fatalf("NewMoney() error = %v", err)
	}

	wallet, err := NewWallet(uuid.New(), currency, initialBalance)
	if err != nil {
		t.Fatalf("newWallet() error = %v", err)
	}

	return wallet
}

func TestWalletDebit(t *testing.T) {
	wallet := newTestWallet(t, 10000, "BRL")

	amount, err := NewMoney(2500, "BRL")
	if err != nil {
		t.Fatal(err)
	}

	err = wallet.Debit(amount)
	if err != nil {
		t.Fatalf("Debit() error = %v", err)
	}

	if wallet.balance.Amount != 7500 {
		t.Errorf("balance = %d, want 7500", wallet.balance.Amount)
	}

	if wallet.version != 2 {
		t.Errorf("version = %d, want 2", wallet.version)
	}
}

func TestWalletDebitInsufficientFunds(t *testing.T) {
	wallet := newTestWallet(t, 10000, "BRL")

	amount, err := NewMoney(12000, "BRL")
	if err != nil {
		t.Fatal(err)
	}

	originalBalance := wallet.balance
	originalVersion := wallet.version
	originalUpdatedAt := wallet.updatedAt

	err = wallet.Debit(amount)
	if err == nil {
		t.Fatal("Debit() expected an error for insufficient funds")
	}

	if wallet.balance != originalBalance {
		t.Errorf("balance changed after failed debit: got %+v, want %+v",
			wallet.balance, originalBalance)
	}

	if wallet.version != originalVersion {
		t.Errorf("version changed after failed debit: got %d, want %d",
			wallet.version, originalVersion)
	}

	if !wallet.updatedAt.Equal(originalUpdatedAt) {
		t.Error("updatedAt changed after failed debit")
	}
}

func TestWalletCredit(t *testing.T) {
	wallet := newTestWallet(t, 10000, "BRL")

	amount, err := NewMoney(2500, "BRL")
	if err != nil {
		t.Fatal(err)
	}

	err = wallet.Credit(amount)
	if err != nil {
		t.Fatalf("Credit() error = %v", err)
	}

	if wallet.balance.Amount != 12500 {
		t.Errorf("balance = %d, want 12500", wallet.balance.Amount)
	}

	if wallet.version != 2 {
		t.Errorf("version = %d, want 2", wallet.version)
	}
}

func TestWalletRejectsCurrencyMismatch(t *testing.T) {
	tests := []struct {
		name     string
		movement func(*Wallet, Money) error
	}{
		{
			name: "debit",
			movement: func(w *Wallet, amount Money) error {
				return w.Debit(amount)
			},
		},
		{
			name: "credit",
			movement: func(w *Wallet, amount Money) error {
				return w.Credit(amount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wallet := newTestWallet(t, 10000, "BRL")

			amount, err := NewMoney(1000, "USD")
			if err != nil {
				t.Fatal(err)
			}

			originalBalance := wallet.balance
			originalVersion := wallet.version

			err = tt.movement(wallet, amount)
			if err == nil {
				t.Fatal("expected currency mismatch error")
			}

			if wallet.balance != originalBalance {
				t.Errorf("balance changed: got %+v, want %+v",
					wallet.balance, originalBalance)
			}

			if wallet.version != originalVersion {
				t.Errorf("version changed: got %d, want %d",
					wallet.version, originalVersion)
			}
		})
	}
}

func TestWalletRejectsZeroOrNegativeMovement(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
	}{
		{name: "zero", amount: 0},
		{name: "negative", amount: -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, operation := range []struct {
				name     string
				movement func(*Wallet, Money) error
			}{
				{
					name: "debit",
					movement: func(w *Wallet, amount Money) error {
						return w.Debit(amount)
					},
				},
				{
					name: "credit",
					movement: func(w *Wallet, amount Money) error {
						return w.Credit(amount)
					},
				},
			} {
				t.Run(operation.name, func(t *testing.T) {
					wallet := newTestWallet(t, 10000, "BRL")

					amount, err := NewMoney(tt.amount, "BRL")
					if err != nil {
						t.Fatal(err)
					}

					originalBalance := wallet.balance
					originalVersion := wallet.version

					err = operation.movement(wallet, amount)
					if err == nil {
						t.Fatal("expected an error for zero or negative movement")
					}

					if wallet.balance != originalBalance {
						t.Errorf("balance changed: got %+v, want %+v",
							wallet.balance, originalBalance)
					}

					if wallet.version != originalVersion {
						t.Errorf("version changed: got %d, want %d",
							wallet.version, originalVersion)
					}
				})
			}
		})
	}
}

func TestNewWalletRejectsInvalidInitialBalance(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		amount   int64
	}{
		{
			name:     "currency mismatch",
			currency: "BRL",
			amount:   10000,
		},
		{
			name:     "negative balance",
			currency: "BRL",
			amount:   -100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moneyCurrency := tt.currency
			if tt.name == "currency mismatch" {
				moneyCurrency = "USD"
			}

			initialBalance, err := NewMoney(tt.amount, moneyCurrency)
			if err != nil {
				t.Fatal(err)
			}

			_, err = NewWallet(uuid.New(), tt.currency, initialBalance)
			if err == nil {
				t.Fatal("newWallet() expected an error")
			}
		})
	}
}
