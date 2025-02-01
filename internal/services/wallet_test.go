package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"wallet-test/internal/repositories"
	"wallet-test/internal/services"
)

type mockRepo struct {
	balances map[uuid.UUID]int64
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		balances: make(map[uuid.UUID]int64),
	}
}

func (m *mockRepo) ChangeBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	current := m.balances[walletID]
	newBalance := current + amount
	if newBalance < 0 {
		return repositories.ErrNotEnoughFunds
	}
	m.balances[walletID] = newBalance
	return nil
}

func (m *mockRepo) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	return m.balances[walletID], nil
}

func TestWalletService(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	srv := services.NewWalletService(repo)

	wID := uuid.New()

	// Начальный баланс — 0
	bal, err := srv.GetBalance(ctx, wID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), bal)

	// Депозит 1000
	err = srv.Deposit(ctx, wID, 1000)
	assert.NoError(t, err)

	// Баланс -> 1000
	bal, err = srv.GetBalance(ctx, wID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1000), bal)

	// Снимаем 300
	err = srv.Withdraw(ctx, wID, 300)
	assert.NoError(t, err)

	// Баланс -> 700
	bal, _ = srv.GetBalance(ctx, wID)
	assert.Equal(t, int64(700), bal)

	// Невалидная сумма (<= 0)
	err = srv.Withdraw(ctx, wID, 0)
	assert.Error(t, err)

	// Пробуем снять 800 (превысить 700)
	err = srv.Withdraw(ctx, wID, 800)
	assert.Error(t, err)

	// Баланс остаётся 700
	bal, _ = srv.GetBalance(ctx, wID)
	assert.Equal(t, int64(700), bal)
}
