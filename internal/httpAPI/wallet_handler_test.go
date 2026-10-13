package httpapi

import (
	"backend-challange/internal/postgres"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestHandlerWallet_CreatesWallet(t *testing.T) {
	body := `{
		"playerId": "550e8400-e29b-41d4-a716-446655440000",
		"initialBalance": {
			"amount": "1000.00",
			"currency": "BRL"
		}
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/wallets",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()


	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	store := postgres.NewWalletStore(pool)
	handler := NewWalletHandler(store)

handler.HandlerCreateWallet(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var response createWalletResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.PlayerID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("unexpected player ID: %s", response.PlayerID)
	}

	if response.Balance.Amount != "1000.00" {
		t.Errorf("expected balance 1000.00, got %s", response.Balance.Amount)
	}

	if response.Balance.Currency != "BRL" {
		t.Errorf("expected currency BRL, got %s", response.Balance.Currency)
	}

	if response.Version != 1 {
		t.Errorf("expected version 1, got %d", response.Version)
	}
}

