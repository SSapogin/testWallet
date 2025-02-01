package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"wallet-test/internal/repositories"
)

// WalletService — интерфейс, описывающий работу с кошельком
type WalletService interface {
	// Deposit — вносит средства на кошелёк
	Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error

	// Withdraw — списывает средства с кошелька
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error

	// GetBalance — получить текущий баланс
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
}

type walletService struct {
	repo repositories.WalletRepository
}

// NewWalletService — конструктор сервиса.
func NewWalletService(r repositories.WalletRepository) WalletService {
	return &walletService{
		repo: r,
	}
}

// Deposit — пополнение кошелька
func (s *walletService) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма для депозита должна быть > 0 (получено %d)", amount)
	}
	return s.repo.ChangeBalance(ctx, walletID, amount)
}

// Withdraw — списание средств
func (s *walletService) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("сумма для списания должна быть > 0 (получено %d)", amount)
	}
	return s.repo.ChangeBalance(ctx, walletID, -amount)
}

// GetBalance — текущий баланса
func (s *walletService) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	return s.repo.GetBalance(ctx, walletID)
}
