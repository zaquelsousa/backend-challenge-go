package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func validWagerInput() NewWagerTransactionInput {
	return NewWagerTransactionInput{
		ProviderID:            "provider-1",
		ExternalTransactionID: "external-tx-1",
		IdempotencyKey:        "idempotency-1",
		PayloadHash:           "hash-1",
		WalletID:              uuid.New(),
		PlayerID:              uuid.New(),
		RoundID:               "round-1",
		GameID:                "game-1",
		Kind:                  WagerBet,
		Money: Money{
			Amount:   1000,
			Currency: "BRL",
		},
	}
}

func TestNewWagerTransaction_Valid(t *testing.T) {
	tests := []struct {
		name string
		kind WagerKind
		money Money
		ref  string
	}{
		{
			name: "BET",
			kind: WagerBet,
			money: Money{Amount: 1000, Currency: "BRL"},
		},
		{
			name: "WIN",
			kind: WagerWin,
			money: Money{Amount: 1000, Currency: "BRL"},
		},
		{
			name: "LOSS",
			kind: WagerLoss,
			money: Money{Amount: 0, Currency: "BRL"},
		},
		{
			name: "REFUND",
			kind: WagerRefund,
			money: Money{Amount: 1000, Currency: "BRL"},
			ref:  "original-tx-1",
		},
		{
			name: "ROLLBACK",
			kind: WagerRollback,
			money: Money{Amount: 1000, Currency: "BRL"},
			ref:  "original-tx-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validWagerInput()
			in.Kind = tt.kind
			in.Money = tt.money
			in.ReferenceExternalID = tt.ref

			tx, err := NewWagerTransaction(in)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tx == nil {
				t.Fatal("expected transaction, got nil")
			}

			if tx.id == uuid.Nil {
				t.Error("expected generated transaction ID")
			}

			if tx.status != WagerPending {
				t.Errorf("expected status %q, got %q",
					WagerPending, tx.status)
			}

			if tx.kind != tt.kind {
				t.Errorf("expected kind %q, got %q",
					tt.kind, tx.kind)
			}

			if !tx.createdAt.Equal(tx.updatedAt) {
				t.Error("expected createdAt and updatedAt to be equal")
			}

			if tx.createdAt.IsZero() {
				t.Error("expected createdAt to be set")
			}

			if tx.createdAt.Location() != time.UTC {
				t.Error("expected createdAt to use UTC")
			}
		})
	}
}

func TestNewWagerTransaction_RequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*NewWagerTransactionInput)
		wantErr string
	}{
		{
			name: "missing provider ID",
			mutate: func(in *NewWagerTransactionInput) {
				in.ProviderID = ""
			},
			wantErr: "provider ID is required",
		},
		{
			name: "missing external transaction ID",
			mutate: func(in *NewWagerTransactionInput) {
				in.ExternalTransactionID = ""
			},
			wantErr: "external transaction ID is required",
		},
		{
			name: "missing idempotency key",
			mutate: func(in *NewWagerTransactionInput) {
				in.IdempotencyKey = ""
			},
			wantErr: "idempotency key is required",
		},
		{
			name: "missing payload hash",
			mutate: func(in *NewWagerTransactionInput) {
				in.PayloadHash = ""
			},
			wantErr: "payload hash is required",
		},
		{
			name: "missing wallet ID",
			mutate: func(in *NewWagerTransactionInput) {
				in.WalletID = uuid.Nil
			},
			wantErr: "wallet ID is required",
		},
		{
			name: "missing player ID",
			mutate: func(in *NewWagerTransactionInput) {
				in.PlayerID = uuid.Nil
			},
			wantErr: "player ID is required",
		},
		{
			name: "missing round ID",
			mutate: func(in *NewWagerTransactionInput) {
				in.RoundID = "  "
			},
			wantErr: "round ID is required",
		},
		{
			name: "missing game ID",
			mutate: func(in *NewWagerTransactionInput) {
				in.GameID = " "
			},
			wantErr: "game ID is required",
		},
		{
			name: "missing currency",
			mutate: func(in *NewWagerTransactionInput) {
				in.Money.Currency = ""
			},
			wantErr: "currency is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validWagerInput()
			tt.mutate(&in)

			tx, err := NewWagerTransaction(in)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tx != nil {
				t.Error("expected nil transaction on error")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %q",
					tt.wantErr, err.Error())
			}
		})
	}
}

func TestNewWagerTransaction_InvalidKind(t *testing.T) {
	tests := []struct {
		name string
		kind WagerKind
	}{
		{name: "unknown kind", kind: WagerKind("UNKNOWN")},
		{name: "opening kind", kind: WagerOpening},
		{name: "empty kind", kind: WagerKind("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validWagerInput()
			in.Kind = tt.kind

			tx, err := NewWagerTransaction(in)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tx != nil {
				t.Error("expected nil transaction on error")
			}
		})
	}
}

func TestNewWagerTransaction_InvalidCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
	}{
		{name: "too short", currency: "BR"},
		{name: "too long", currency: "BRLX"},
		{name: "lowercase", currency: "brl"},
		{name: "mixed case", currency: "Brl"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validWagerInput()
			in.Money.Currency = tt.currency

			tx, err := NewWagerTransaction(in)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tx != nil {
				t.Error("expected nil transaction on error")
			}
		})
	}
}

func TestNewWagerTransaction_AmountRules(t *testing.T) {
	tests := []struct {
		name      string
		kind      WagerKind
		amount    int64
		wantError bool
	}{
		{"BET positive", WagerBet, 100, false},
		{"BET zero", WagerBet, 0, true},
		{"BET negative", WagerBet, -100, true},
		{"WIN positive", WagerWin, 100, false},
		{"WIN zero", WagerWin, 0, true},
		{"LOSS zero", WagerLoss, 0, false},
		{"LOSS positive", WagerLoss, 100, true},
		{"LOSS negative", WagerLoss, -100, true},
		{"REFUND positive", WagerRefund, 100, false},
		{"REFUND zero", WagerRefund, 0, true},
		{"ROLLBACK positive", WagerRollback, 100, false},
		{"ROLLBACK zero", WagerRollback, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validWagerInput()
			in.Kind = tt.kind
			in.Money.Amount = tt.amount

			if tt.kind == WagerRefund || tt.kind == WagerRollback {
				in.ReferenceExternalID = "original-tx-1"
			}

			tx, err := NewWagerTransaction(in)

			if tt.wantError && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if err != nil && tx != nil {
				t.Error("expected nil transaction on error")
			}
		})
	}
}

func TestNewWagerTransaction_ReferenceRequired(t *testing.T) {
	tests := []struct {
		name string
		kind WagerKind
	}{
		{name: "refund", kind: WagerRefund},
		{name: "rollback", kind: WagerRollback},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validWagerInput()
			in.Kind = tt.kind
			in.ReferenceExternalID = "  "

			tx, err := NewWagerTransaction(in)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tx != nil {
				t.Error("expected nil transaction on error")
			}

			if !strings.Contains(err.Error(), "requires a reference transaction ID") {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

