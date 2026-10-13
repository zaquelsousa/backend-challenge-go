package postgres

import (
	dto "backend-challange/internal/DTO"
	"backend-challange/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)



type WalletStore struct {
	db *pgxpool.Pool
}

func NewWalletStore(db *pgxpool.Pool) *WalletStore {
	return &WalletStore{db: db}
}

func formatMinorUnits(amount int64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	return fmt.Sprintf("%s%d.%02d", sign, amount/100, amount%100)
}

func (s *WalletStore) Create(ctx context.Context, w *domain.Wallet) error {
	balance := w.Balance()

	_, err := s.db.Exec(ctx, `
		INSERT INTO wallets (
			id, player_id, currency, balance, version, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		w.ID(),
		w.PlayerID(),
		w.Currency(),
		balance.Amount,
		w.Version(),
		w.CreatedAt(),
		w.UpdateAt(),
	)

	if err != nil {
		return fmt.Errorf("insert wallet: %w", err)
	}

	return nil
}

func (s *WalletStore) GetWallet(ctx context.Context, walletID uuid.UUID,) (dto.WalletResponse, error) {
	var wallet dto.WalletResponse
	var amountMinor int64
	var currency string

	err := s.db.QueryRow(ctx, `
		SELECT id, player_id, currency, balance, version
		FROM wallets
		WHERE id = $1
	`, walletID).Scan(
		&wallet.ID,
		&wallet.PlayerID,
		&currency,
		&amountMinor,
		&wallet.Version,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.WalletResponse{}, fmt.Errorf("record not found")
	}
	if err != nil {
		return dto.WalletResponse{}, err
	}

	wallet.Balance = dto.MoneyResponse{
		Amount:   formatMinorUnits(amountMinor),
		Currency: currency,
	}

	return wallet, nil
}
