package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wallet-test/internal/entities"
)

var ErrNotEnoughFunds = errors.New("Недостаточно средств")

// WalletRepository — интерфейс, описывающий работу с кошельком
type WalletRepository interface {
	// ChangeBalance — изменить баланс (amount > 0 => пополнение; amount < 0 => списание).
	ChangeBalance(ctx context.Context, walletID uuid.UUID, amount int64) error

	// GetBalance — вернуть текущий баланс кошелька (или 0, если нет записи).
	GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error)
}

// walletRepo — реализация интерфейса
type walletRepo struct {
	db *gorm.DB
}

// NewWalletRepository — конструктор репозитория
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepo{db: db}
}

// ChangeBalance — меняем баланс.
func (r *walletRepo) ChangeBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var w entities.Wallet

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&w, "id = ?", walletID).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Если записи нет, создаём её с нулевым балансом.
				w.ID = walletID
				w.Balance = 0
				if err2 := tx.Create(&w).Error; err2 != nil {
					return err2
				}
			} else {
				return err
			}
		}

		newBalance := w.Balance + amount
		if newBalance < 0 {
			return ErrNotEnoughFunds
		}

		// Обновляем баланс
		w.Balance = newBalance
		if err2 := tx.Save(&w).Error; err2 != nil {
			return err2
		}

		return nil
	})
}

// GetBalance — возвращает текущий баланс кошелька.
func (r *walletRepo) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	var w entities.Wallet

	err := r.db.WithContext(ctx).
		First(&w, "id = ?", walletID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return w.Balance, nil
}
