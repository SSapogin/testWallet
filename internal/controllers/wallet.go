package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"wallet-test/internal/services"
)

// WalletController — интерфейс
type WalletController interface {
	// ChangeBalance — изменяет баланс
	ChangeBalance(ctx *gin.Context)

	// GetBalance — получает текущий баланс
	GetBalance(ctx *gin.Context)
}

// WalletController — контроллер для операций с кошельком
type walletController struct {
	service services.WalletService
}

// NewWalletController — конструктор для контроллера, принимает WalletService
func NewWalletController(srv services.WalletService) WalletController {
	return &walletController{
		service: srv,
	}
}

// processOperationDTO — модель для входящего JSON при POST /wallet
type processOperationDTO struct {
	WalletID      string `json:"walletId"`
	OperationType string `json:"operationType"` // "DEPOSIT" / "WITHDRAW"
	Amount        int64  `json:"amount"`
}

// ChangeBalance — обрабатывает POST
func (c *walletController) ChangeBalance(ctx *gin.Context) {
	var dto processOperationDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Невалидный JSON: " + err.Error()})
		return
	}

	walletUUID, err := uuid.Parse(dto.WalletID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверное значение walletId"})
		return
	}

	switch dto.OperationType {
	case "DEPOSIT":
		if err := c.service.Deposit(ctx, walletUUID, dto.Amount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	case "WITHDRAW":
		if err := c.service.Withdraw(ctx, walletUUID, dto.Amount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неизвестное значение operationType (используйте DEPOSIT или WITHDRAW)"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetBalance — обрабатывает GET /api/v1/wallets/:id
func (c *walletController) GetBalance(ctx *gin.Context) {
	idParam := ctx.Param("id")
	walletUUID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка парсинга UUID"})
		return
	}

	balance, err := c.service.GetBalance(ctx, walletUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"walletId": walletUUID.String(),
		"balance":  balance,
	})
}
