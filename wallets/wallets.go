package wallets

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const defaultBalance = 100

type Wallet struct {
	ID      uuid.UUID `bson:"_id" json:"id"`
	Balance int64     `bson:"balance" json:"balance"`
}

const TransactionHistoryPrefix = "history_"

type Transaction struct {
	Time   time.Time `bson:"time" json:"time"`
	From   string    `bson:"from" json:"from"`
	To     string    `bson:"to" json:"to"`
	Amount int64     `bson:"amount" json:"amount"`
}

type CreateTransactionParams struct {
	From   uuid.UUID
	To     uuid.UUID
	Amount int64
}

type Interface interface {
	CreateWallet(ctx context.Context) (Wallet, error)
	CreateTransaction(ctx context.Context, params CreateTransactionParams) error
	GetHistory(ctx context.Context, id uuid.UUID) ([]Transaction, error)
	GetWallet(ctx context.Context, id uuid.UUID) (Wallet, error)
}
