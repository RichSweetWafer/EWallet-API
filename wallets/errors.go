package wallets

import (
	"fmt"

	"github.com/google/uuid"
)

type DuplicateKeyError struct {
	ID uuid.UUID
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("Duplicate wallet uuid: %v", e.ID)
}

type WalletNotFoundError struct{}

func (e *WalletNotFoundError) Error() string {
	return "Wallet not found"
}

type WalletBalanceLow struct{}

func (e *WalletBalanceLow) Error() string {
	return "Wallet balance is too low"
}
