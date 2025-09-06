package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	mocks_service "cli-inventory/internal/mocks/service"
	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddStockCmd(t *testing.T) {
	// Save original stockService
	originalStockService := stockService
	defer func() {
		stockService = originalStockService
	}()

	// Create mock repositories and service
	mockProductRepo := mocks_service.NewMockProductRepositoryInterface(t)
	mockLocationRepo := mocks_service.NewMockLocationRepositoryInterface(t)
	mockStockRepo := mocks_service.NewMockStockRepositoryInterface(t)
	mockMovementRepo := mocks_service.NewMockStockMovementRepositoryInterface(t)
	
	// Create a mock database pool (can be nil for our tests)
	var mockDB *pgxpool.Pool
	
	stockService = service.NewStockService(mockProductRepo, mockLocationRepo, mockStockRepo, mockMovementRepo, mockDB)

	t.Run("Successful stock addition", func(t *testing.T) {
		expectedStock := &models.Stock{
			ID:         1,
			ProductID:  1,
			LocationID: 1,
			Quantity:   100,
		}

		// Set up expectations
		mockProductRepo.EXPECT().GetByID(mock.Anything, 1).Return(&models.Product{}, nil)
		mockLocationRepo.EXPECT().GetByID(mock.Anything, 1).Return(&models.Location{}, nil)
		mockStockRepo.EXPECT().AddStock(mock.Anything, 1, 1, 100).Return(expectedStock, nil)
		mockMovementRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.StockMovement")).Return(&models.StockMovement{}, nil)

		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "add-stock",
			Short: "Add stock for a product at a specific location",
			Long: `Add stock quantity for a specific product at a given location.
This will increase the stock level for the product at the specified location.`,
			Args: cobra.ExactArgs(3),
			Run:  addStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "1", "100"})

		// Capture output by redirecting os.Stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := testCmd.Execute()
		assert.NoError(t, err)

		// Close the write end and restore stdout
		w.Close()
		os.Stdout = old

		// Read the output
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Check output
		assert.Contains(t, output, "Stock added successfully")
		assert.Contains(t, output, "Product ID: 1")
		assert.Contains(t, output, "Location ID: 1")
		assert.Contains(t, output, "New Quantity: 100")
	})

	t.Run("Invalid product ID", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "add-stock",
			Short: "Add stock for a product at a specific location",
			Long: `Add stock quantity for a specific product at a given location.
This will increase the stock level for the product at the specified location.`,
			Args: cobra.ExactArgs(3),
			Run:  addStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"invalid", "1", "100"})

		// Capture output by redirecting os.Stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := testCmd.Execute()
		assert.NoError(t, err) // Command should not return error, just print to stdout

		// Close the write end and restore stdout
		w.Close()
		os.Stdout = old

		// Read the output
		var buf bytes.Buffer
		io.Copy(&buf, r)
		output := buf.String()

		// Check output
		assert.Contains(t, output, "Error: Invalid product ID")
	})
}