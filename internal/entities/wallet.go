package entities

import (
	"time"

	"github.com/google/uuid"
)

// Wallet — модель
type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"` // Поле-ключ, UUID
	Balance   int64     `gorm:"not null;default:0"`   // Текущее значение баланса
	CreatedAt time.Time
	UpdatedAt time.Time
}
