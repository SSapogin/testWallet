package server

import (
	"github.com/gin-gonic/gin"
	"wallet-test/internal/controllers"
)

func NewRouter(ctrl controllers.WalletController) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/wallet", ctrl.ChangeBalance)
		v1.GET("/wallets/:id", ctrl.GetBalance)
	}

	return r
}
