package api

import (
	"fetch/models"

	"github.com/gin-gonic/gin"
)

var receiptLink = []models.ReceiptLink{}

var router = gin.Default()

func Run() {
	getRoutes()
	router.Run("localhost:8080")
}

func getRoutes() {
	receiptsGroup := router.Group("/receipts")

	addReceiptsRoutes(receiptsGroup)
}
