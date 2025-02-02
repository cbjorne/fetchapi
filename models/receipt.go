package models

import (
	"strconv"
	"time"

	"errors"
	"strings"

	"github.com/google/uuid"
)

func GetDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

func GetTime(s string) (int, int, error) {
	i := strings.IndexByte(s, ':')

	if i <= -1 {
		return -1, -1, errors.New("invalid time")
	}

	hours, hoursErr := strconv.Atoi(s[:i])
	minutes, minutesErr := strconv.Atoi(s[i+1:])

	if hoursErr != nil || minutesErr != nil || hours > 24 || hours < 0 || minutes > 59 || minutes < 0 {
		return -1, -1, errors.New("invalid time")
	}

	return hours, minutes, nil

}

func GetTotal(t string) (float64, error) {
	i := strings.IndexByte(t, '.')

	if i <= -1 || len(t)-i-1 != 2 {
		return -1, errors.New("invalid decimal")
	}

	return strconv.ParseFloat(t, 64)
}

func GetItems(itemsRequest []ItemRequest) ([]Item, error) {
	var items []Item

	for i := 0; i < len(itemsRequest); i++ {
		desc := itemsRequest[i].ShortDescription
		price, err := GetTotal(itemsRequest[i].Price)

		if err != nil {
			return items, err
		}

		items = append(items, Item{ShortDescription: desc, Price: price})
	}

	return items, nil
}

type ReceiptLink struct {
	Id     uuid.UUID `json:"id"`
	Points int       `json:"points"`
}

type Receipt struct {
	Retailer            string
	PurchaseDate        time.Time
	PurchaseTimeHours   int
	PurchaseTimeMinutes int
	Items               []Item
	Total               float64
}

type Item struct {
	ShortDescription string
	Price            float64
}

type ReceiptRequest struct {
	Retailer     string        `json:"retailer"`
	PurchaseDate string        `json:"purchaseDate"`
	PurchaseTime string        `json:"purchaseTime"`
	Items        []ItemRequest `json:"items"`
	Total        string        `json:"total"`
}

type ItemRequest struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
