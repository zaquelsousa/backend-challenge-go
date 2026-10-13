package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)


type WagerKind string

const (
    WagerOpening  WagerKind = "OPENING"
    WagerBet      WagerKind = "BET"
    WagerWin      WagerKind = "WIN"
    WagerLoss     WagerKind = "LOSS"
    WagerRefund   WagerKind = "REFUND"
    WagerRollback WagerKind = "ROLLBACK"
)

type WagerStatus string

const (
    WagerPending          WagerStatus = "PENDING"
    WagerPendingReference WagerStatus = "PENDING_REFERENCE"
    WagerProcessed        WagerStatus = "PROCESSED"
    WagerRejected         WagerStatus = "REJECTED"
    WagerFailed           WagerStatus = "FAILED"
)

type WagerTransaction struct {
    id                          uuid.UUID
    providerID                  string
    externalTransactionID       string
    idempotencyKey              string
    payloadHash                 string

    walletID                    uuid.UUID
    playerID                    uuid.UUID
    roundID                     string
    gameID                      string

    kind                        WagerKind
    money                       Money
    referenceExternalID         string

    status                      WagerStatus
    failureCode                 string

    createdAt                   time.Time
    updatedAt                   time.Time
}


//to build the construtor
type NewWagerTransactionInput struct {
	ProviderID                  string
	ExternalTransactionID       string
	IdempotencyKey              string
	PayloadHash                 string
	WalletID                    uuid.UUID
	PlayerID                    uuid.UUID
	RoundID                     string
	GameID                      string
	Kind                        WagerKind
	Money                       Money
	ReferenceExternalID         string
}

func NewWagerTransaction(in NewWagerTransactionInput) (*WagerTransaction, error){
	if in.ProviderID == "" {
		return nil, fmt.Errorf("provider ID is required")
	}

	if in.ExternalTransactionID == "" {
		return nil, fmt.Errorf("external transaction ID is required")
	}

	if in.IdempotencyKey == "" {
		return nil, fmt.Errorf("idempotency key is required")
	}

	if in.PayloadHash == "" {
		return nil, fmt.Errorf("payload hash is required")
	}

	if in.WalletID == uuid.Nil {
		return nil, fmt.Errorf("wallet ID is required")
	}

	if in.PlayerID == uuid.Nil {
		return nil, fmt.Errorf("player ID is required")
	}

	if strings.TrimSpace(in.RoundID) == "" {
		return nil, fmt.Errorf("round ID is required")
	}

	if strings.TrimSpace(in.GameID) == "" {
		return nil, fmt.Errorf("game ID is required")
	}

	switch in.Kind {
	case WagerBet, WagerWin, WagerLoss,
		WagerRefund, WagerRollback:
		// Valid external operation.
	case WagerOpening:
		return nil, fmt.Errorf(
			"OPENING is reserved for internal wallet creation",
		)
	default:
		return nil, fmt.Errorf("invalid wager kind: %q", in.Kind)
	}
	

	if in.Money.Currency == "" {
		return nil, fmt.Errorf("currency is required")
	}

	if len(in.Money.Currency) != 3 ||
		in.Money.Currency != strings.ToUpper(in.Money.Currency) {
		return nil, fmt.Errorf("invalid currency code: %q", in.Money.Currency)
	}

	switch in.Kind {
	case WagerLoss:
		if in.Money.Amount != 0 {
			return nil, fmt.Errorf("LOSS must have an amount of zero")
		}
	case WagerBet, WagerWin, WagerRefund, WagerRollback:
		if in.Money.Amount <= 0 {
			return nil, fmt.Errorf(
				"%s requires a positive amount",
				in.Kind,
			)
		}
	}

	switch in.Kind {
	case WagerRefund, WagerRollback:
		if strings.TrimSpace(in.ReferenceExternalID) == "" {
			return nil, fmt.Errorf(
				"%s requires a reference transaction ID",
				in.Kind,
			)
		}
	}

	now := time.Now().UTC()
	
	return &WagerTransaction{
		id:                    uuid.New(),
		providerID:            in.ProviderID,
		externalTransactionID: in.ExternalTransactionID,
		idempotencyKey:        in.IdempotencyKey,
		payloadHash:           in.PayloadHash,
		walletID:              in.WalletID,
		playerID:              in.PlayerID,
		roundID:               in.RoundID,
		gameID:                in.GameID,
		kind:                  in.Kind,
		money:                 in.Money,
		referenceExternalID:   in.ReferenceExternalID,
		status:                WagerPending,
		createdAt:             now,
		updatedAt:             now,
	}, nil
}
