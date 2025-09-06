package testutils

import (
	"context"
	"math/rand"
	"time"

	"cli-inventory/internal/models"
)

// CreateTestProduct creates a test product with random data
func CreateTestProduct() *models.CreateProductRequest {
	return &models.CreateProductRequest{
		SKU:         generateRandomSKU(),
		Name:        generateRandomName("Product"),
		Description: generateRandomDescription(),
		Price:       generateRandomPrice(),
	}
}

// CreateTestLocation creates a test location with random data
func CreateTestLocation() *models.CreateLocationRequest {
	return &models.CreateLocationRequest{
		Name: generateRandomName("Location"),
	}
}

// CreateTestStockRequest creates a test stock request
func CreateTestStockRequest(productID, locationID int) *models.AddStockRequest {
	return &models.AddStockRequest{
		ProductID:  productID,
		LocationID: locationID,
		Quantity:   generateRandomQuantity(),
	}
}

// CreateTestMoveStockRequest creates a test move stock request
func CreateTestMoveStockRequest(productID, fromLocationID, toLocationID int) *models.MoveStockRequest {
	return &models.MoveStockRequest{
		ProductID:      productID,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
		Quantity:       generateRandomQuantity(),
	}
}

// generateRandomSKU generates a random SKU
func generateRandomSKU() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return "TEST-" + string(b)
}

// generateRandomName generates a random name with a prefix
func generateRandomName(prefix string) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return prefix + "-" + string(b)
}

// generateRandomDescription generates a random description
func generateRandomDescription() string {
	descriptions := []string{
		"A high-quality test product",
		"Durable and reliable item",
		"Premium quality material",
		"Essential for testing purposes",
		"Manufactured with care",
		"Test item with great features",
		"Reliable and efficient",
		"Perfect for demonstration",
	}
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return descriptions[seededRand.Intn(len(descriptions))]
}

// generateRandomPrice generates a random price between 1.00 and 999.99
func generateRandomPrice() float64 {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return float64(seededRand.Intn(99900)+100) / 100
}

// generateRandomQuantity generates a random quantity between 1 and 100
func generateRandomQuantity() int {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return seededRand.Intn(100) + 1
}

// ProductTestHelper provides helper methods for product testing
type ProductTestHelper struct {
	ctx context.Context
}

func NewProductTestHelper(ctx context.Context) *ProductTestHelper {
	return &ProductTestHelper{ctx: ctx}
}

// LocationTestHelper provides helper methods for location testing
type LocationTestHelper struct {
	ctx context.Context
}

func NewLocationTestHelper(ctx context.Context) *LocationTestHelper {
	return &LocationTestHelper{ctx: ctx}
}

// StockTestHelper provides helper methods for stock testing
type StockTestHelper struct {
	ctx context.Context
}

func NewStockTestHelper(ctx context.Context) *StockTestHelper {
	return &StockTestHelper{ctx: ctx}
}
