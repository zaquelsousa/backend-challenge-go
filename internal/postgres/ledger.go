package postgres

import (
	"context"
	"fmt"

	"backend-challange/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerStore struct {
	pool *pgxpool.Pool
}

func NewLedgerStore(pool *pgxpool.Pool) *LedgerStore {
	return &LedgerStore{
		pool: pool,
	}
}

func (s *LedgerStore) Create(ctx context.Context, entry domain.LedgerEntry) error {
    const query = `
        INSERT INTO wallet_ledger_entries (
            id,
            wallet_id,
            transaction_id,
            direction,
            amount_minor,
            currency,
            balance_before_minor,
            balance_after_minor,
            created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `

    _, err := s.pool.Exec(
        ctx,
        query,
        entry.ID,
        entry.WalletID,
        entry.TransactionID,
        entry.Direction,
        entry.Money.Amount,
        entry.Money.Currency,
        entry.BalanceBefore.Amount,
        entry.BalanceAfter.Amount,
        entry.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("create ledger entry: %w", err)
    }

    return nil
}

func (s *LedgerStore) ListByWallet(ctx context.Context, 
		walletID uuid.UUID, limit int,) ([]domain.LedgerEntry, error){
	if limit < 1 || limit > 100 {
		return nil, fmt.Errorf("limit must be between 1 and 100")
	}

	const query = `
		SELECT
			id,
			wallet_id,
			transaction_id,
			direction,
			amount_minor,
			currency,
			balance_before_minor,
			balance_after_minor,
			created_at
		FROM wallet_ledger_entries
		WHERE wallet_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2
	`

	rows, err := s.pool.Query(ctx, query, walletID, limit)
	if err != nil {
		return nil, fmt.Errorf("query wallet ledger: %w", err)
	}
	defer rows.Close()

	entries := make([]domain.LedgerEntry, 0, limit)

	for rows.Next() {
		var entry domain.LedgerEntry

		err := rows.Scan(
			&entry.ID,
			&entry.WalletID,
			&entry.TransactionID,
			&entry.Direction,
			&entry.Money.Amount,
			&entry.Money.Currency,
			&entry.BalanceBefore.Amount,
			&entry.BalanceAfter.Amount,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wallet ledger entry: %w", err)
		}

		entry.BalanceBefore.Currency = entry.Money.Currency
		entry.BalanceAfter.Currency = entry.Money.Currency

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wallet ledger: %w", err)
	}

	return entries, nil
}

