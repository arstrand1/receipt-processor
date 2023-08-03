package main

import (
	"errors"   // error handling
	"net/http" // handles HTTP status codes
	"strings"  // used to trim item desc
	"time"     // handles date/time formatting
	"unicode"  // used to find alphanumeric chars in retailer

	"github.com/gin-gonic/gin"      // handles router
	"github.com/google/uuid"        // generates uuid
	"github.com/shopspring/decimal" // used for determining points from prices
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate" time_format:"2006-01-02"`
	PurchaseTime string  `json:"purchaseTime" time_format:"10:10"`
	Total        string  `json:"total"`
	Items        []*Item `json:"items" binding:"required"`
	Points       int64   `json:"points"`
}

// Nested structs validation: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price" binding:"required,oneof=post nested"`
}

// Map declaration
var Receipts map[string]Receipt

// getPoints looks up the receipt by the ID and
// returns an object specifying the points awarded
func getPoints(c *gin.Context) {
	id := c.Param("id")
	receipt, valid := Receipts[id]
	if !valid {
		c.IndentedJSON(http.StatusNotFound, map[string]interface{}{"description": "No receipt found for that id"})
		return
	}
	c.IndentedJSON(http.StatusOK, map[string]interface{}{"points": receipt.Points})
}

// postReceipts takes a JSON receipt and returns a
// JSON object mapped to its assigned ID
func postReceipts(c *gin.Context) {
	var newReceipt Receipt
	err := c.BindJSON(&newReceipt)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"description": "The receipt is invalid"})
		return
	}

	// Determine Valid Receipt Points
	err = validateReceiptPoints(&newReceipt)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"description": "The receipt is invalid"})
		return
	}

	// Create ID - Generate UUID using google/uuid
	id := uuid.New()

	// Add the new receipt to map
	Receipts[id.String()] = newReceipt
	c.IndentedJSON(http.StatusOK, map[string]interface{}{"id": id.String()})
}

// validateReceiptPoints checks to see how many points
// are assigned to a valid receipt
func validateReceiptPoints(r *Receipt) error {
	var points int64 = 0

	// CHECK RETAILER
	if r.Retailer == "" {
		return errors.New("empty name")
	}
	//**Rule 1: One point per alphanumeric char**
	for _, character := range r.Retailer {
		if unicode.IsLetter(character) || unicode.IsNumber(character) {
			points++
		}
	}

	// CHECK DATE
	if r.PurchaseDate == "" {
		return errors.New("empty date")
	}
	dateLayout := "2006-01-02"
	datePurchased, err := time.Parse(dateLayout, r.PurchaseDate)
	if err != nil {
		return err
	}
	//**Rule 2: 6 points if day purchased is odd**
	if datePurchased.Day()%2 == 1 {
		points += 6
	}

	// CHECK TIME
	if r.PurchaseTime == "" {
		return errors.New("empty time")
	}
	timeLayout := "15:04"
	timePurchased, err := time.Parse(timeLayout, r.PurchaseTime)
	if err != nil {
		return err
	}
	//**Rule 3: 10 points if time purchased 2:00pm - 4:00pm**
	startTime, _ := time.Parse(timeLayout, "14:00")
	endTime, _ := time.Parse(timeLayout, "16:00")
	if timePurchased.After(startTime) && timePurchased.Before(endTime) {
		points += 10
	}

	// CHECK TOTAL
	total, err := decimal.NewFromString(r.Total)
	if err != nil {
		return err
	}
	//**Rule 4: 50 points if total ends in .00**
	if total.IsInteger() {
		points += 50
	}
	//**Rule 5: 25 points if total is multiple of .25**
	if total.Mod(decimal.NewFromFloat(.25)).IsZero() {
		points += 25
	}

	// CHECK ITEMS
	if len(r.Items) < 1 {
		return errors.New("minimum 1 item")
	}
	for i, item := range r.Items {
		//**Rule 6: 5 points per 2 items**
		if i%2 == 1 {
			points += 5
		}

		//**Rule 7: Ceiling(.2*price) points where length
		//  of item description is multiple of 3**
		trimmed := strings.Trim(item.ShortDescription, " ")
		size := len(trimmed)
		if size%3 == 0 {
			price, err := decimal.NewFromString(item.Price)
			if err != nil {
				return err
			}
			// meaning: price = .2 * price
			price = decimal.NewFromFloat(.2).Mul(price)
			// meaning: points += Ceiling(price).toInt()
			points += price.Ceil().IntPart()
		}
	}

	r.Points = points
	return nil
}

func main() {
	// Initialize map
	Receipts = make(map[string]Receipt)

	// Initialize router and link GET/POST methods
	router := gin.Default()
	router.GET("/receipts/:id/points", getPoints)
	router.POST("/receipts/process", postReceipts)

	// Start API server
	router.Run("localhost:8080")
}
