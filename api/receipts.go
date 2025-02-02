package api

import (
	"fetch/models"
	"fmt"

	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
X One point for every alphanumeric character in the retailer name.
X 50 points if the total is a round dollar amount with no cents.
X 25 points if the total is a multiple of 0.25.
X 5 points for every two items on the receipt.
X If the trimmed length of the item description is a multiple of 3,

	multiply the price by 0.2 and round up to the nearest integer.
	The result is the number of points earned.

X 6 points if the day in the purchase date is odd.
10 points if the time of purchase is after 2:00pm and before 4:00pm.
*/

func evaluatePoints(r models.Receipt) int {
	alphanumeric := "abcdefghijklmnopqrstuvwxyz0123456789"
	points := 0

	for i := 0; i < len(r.Retailer); i++ {
		if strings.Contains(alphanumeric, string(r.Retailer[i])) {
			points += 1
		}
	}
	fmt.Println(points)

	if r.Total == math.Trunc(r.Total) {
		points += 50
	}

	if math.Mod(r.Total, 0.25) == 0 {
		points += 25
	}

	points += 5 * (len(r.Items) / 2)

	for i := 0; i < len(r.Items); i++ {
		desc := r.Items[i].ShortDescription

		if len(strings.Trim(desc, " "))%3 == 0 {
			points += int(math.Ceil(r.Items[i].Price * 0.2))
		}
	}

	if r.PurchaseDate.Day()%2 != 0 {
		points += 6
	}

	if validTime(r.PurchaseTimeHours, r.PurchaseTimeMinutes) {
		points += 10
	}

	return points
}

func validTime(hours int, minutes int) bool {
	return (hours == 15) || (hours == 14 && minutes >= 1) || (hours == 16 && minutes >= 1)
}

func mapReceipt(receiptRequest models.ReceiptRequest) (models.Receipt, error) {
	purchaseDate, err := models.GetDate(receiptRequest.PurchaseDate)
	if err != nil {
		return models.Receipt{}, err
	}

	purchaseTimeHours, purchaseTimeMinutes, err := models.GetTime(receiptRequest.PurchaseTime)
	if err != nil {
		return models.Receipt{}, err
	}

	items, err := models.GetItems(receiptRequest.Items)
	if err != nil {
		return models.Receipt{}, err
	}

	total, err := models.GetTotal(receiptRequest.Total)
	if err != nil {
		return models.Receipt{}, err
	}

	return models.Receipt{
		Retailer:            strings.ToLower(receiptRequest.Retailer),
		PurchaseDate:        purchaseDate,
		PurchaseTimeHours:   purchaseTimeHours,
		PurchaseTimeMinutes: purchaseTimeMinutes,
		Items:               items,
		Total:               total,
	}, nil
}

func processReceipt(c *gin.Context) {
	var receiptRequest models.ReceiptRequest

	if err := c.BindJSON(&receiptRequest); err != nil {
		c.JSON(http.StatusBadRequest, "The receipt is invalid.")
		return
	}

	var receipt, err = mapReceipt(receiptRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, "The receipt is invalid.")
		return
	}

	newLink := models.ReceiptLink{
		Id:     uuid.New(),
		Points: evaluatePoints(receipt),
	}

	receiptLink = append(receiptLink, newLink)

	c.JSON(http.StatusOK, newLink)
}

func getPoints(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusNotFound, "No receipt found for that ID.")
		return
	}

	for i := 0; i < len(receiptLink); i++ {
		if receiptLink[i].Id == id {
			c.JSON(http.StatusOK, gin.H{"points": receiptLink[i].Points})
			return
		}
	}

	c.JSON(http.StatusNotFound, "No receipt found for that ID.")
}

func addReceiptsRoutes(rg *gin.RouterGroup) {
	receiptsRoute := rg.Group("/")

	receiptsRoute.POST("/process", processReceipt)
	receiptsRoute.GET("/:id/points", getPoints)
}
