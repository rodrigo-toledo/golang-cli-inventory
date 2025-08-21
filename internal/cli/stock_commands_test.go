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

func TestMoveStockCmd(t *testing.T) {
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

	t.Run("Successful stock move", func(t *testing.T) {
		// Create mock repositories and service for this specific test case
		mockProductRepo := mocks_service.NewMockProductRepositoryInterface(t)
		mockLocationRepo := mocks_service.NewMockLocationRepositoryInterface(t)
		mockStockRepo := mocks_service.NewMockStockRepositoryInterface(t)
		mockMovementRepo := mocks_service.NewMockStockMovementRepositoryInterface(t)
		stockService = service.NewStockService(mockProductRepo, mockLocationRepo, mockStockRepo, mockMovementRepo, nil)

		expectedStock := &models.Stock{
			ID:         1,
			ProductID:  1,
			LocationID: 2,
			Quantity:   50,
		}

		// Set up expectations
		mockProductRepo.EXPECT().GetByID(mock.Anything, 1).Return(&models.Product{}, nil)
		mockLocationRepo.EXPECT().GetByID(mock.Anything, 1).Return(&models.Location{}, nil)
		mockLocationRepo.EXPECT().GetByID(mock.Anything, 2).Return(&models.Location{}, nil)
		mockStockRepo.EXPECT().GetByProductAndLocation(mock.Anything, 1, 1).Return(&models.Stock{Quantity: 100}, nil)
		mockStockRepo.EXPECT().RemoveStock(mock.Anything, 1, 1, 25).Return(&models.Stock{}, nil)
		mockStockRepo.EXPECT().AddStock(mock.Anything, 1, 2, 25).Return(expectedStock, nil)
		mockMovementRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*models.StockMovement")).Return(&models.StockMovement{}, nil)

		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "1", "2", "25"})

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
		assert.Contains(t, output, "Stock moved successfully")
		assert.Contains(t, output, "Product ID: 1")
		assert.Contains(t, output, "From Location: 1")
		assert.Contains(t, output, "To Location: 2")
		assert.Contains(t, output, "Quantity Moved: 25")
		assert.Contains(t, output, "New Quantity at Destination: 50")
	})

	t.Run("Invalid product ID", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"invalid", "1", "2", "25"})

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

	t.Run("Invalid source location ID", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "invalid", "2", "25"})

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
		assert.Contains(t, output, "Error: Invalid source location ID")
	})

	t.Run("Invalid destination location ID", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "1", "invalid", "25"})

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
		assert.Contains(t, output, "Error: Invalid destination location ID")
	})

	t.Run("Invalid quantity", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "1", "2", "invalid"})

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
		assert.Contains(t, output, "Error: Invalid quantity")
	})

	t.Run("Zero quantity", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "1", "2", "0"})

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
		assert.Contains(t, output, "Error: Quantity must be greater than 0")
	})

	t.Run("Same source and destination locations", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "move-stock",
			Short: "Move stock between locations",
			Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
			Args: cobra.ExactArgs(4),
			Run:  moveStockCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"1", "1", "1", "25"})

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
		assert.Contains(t, output, "Error: Source and destination locations cannot be the same")
	})
}

func TestGenerateReportCmd(t *testing.T) {
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

	t.Run("Successful low-stock report generation", func(t *testing.T) {
		expectedStocks := []models.Stock{
			{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   5,
			},
			{
				ID:         2,
				ProductID:  2,
				LocationID: 1,
				Quantity:   8,
			},
		}

		// Set up expectations
		mockStockRepo.EXPECT().GetLowStock(mock.Anything, 10).Return(expectedStocks, nil)

		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "generate-report",
			Short: "Generate inventory reports",
			Long: `Generate various types of inventory reports.
Currently supports low-stock reports with customizable thresholds.`,
			Args: cobra.MinimumNArgs(1),
			Run:  generateReportCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"low-stock", "10"})

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
		assert.Contains(t, output, "Low Stock Report")
		assert.Contains(t, output, "Threshold: 10")
		assert.Contains(t, output, "ID")
		assert.Contains(t, output, "Product")
		assert.Contains(t, output, "Location")
		assert.Contains(t, output, "Quantity")
		assert.Contains(t, output, "1")
		assert.Contains(t, output, "2")
	})

	t.Run("Low-stock report with no results", func(t *testing.T) {
		// Set up expectations
		mockStockRepo.EXPECT().GetLowStock(mock.Anything, 5).Return([]models.Stock{}, nil)

		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "generate-report",
			Short: "Generate inventory reports",
			Long: `Generate various types of inventory reports.
Currently supports low-stock reports with customizable thresholds.`,
			Args: cobra.MinimumNArgs(1),
			Run:  generateReportCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"low-stock", "5"})

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
		assert.Contains(t, output, "No products found with stock below threshold 5")
	})

	t.Run("Invalid threshold", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "generate-report",
			Short: "Generate inventory reports",
			Long: `Generate various types of inventory reports.
Currently supports low-stock reports with customizable thresholds.`,
			Args: cobra.MinimumNArgs(1),
			Run:  generateReportCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"low-stock", "invalid"})

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
		assert.Contains(t, output, "Error: Invalid threshold")
	})

	t.Run("Negative threshold", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "generate-report",
			Short: "Generate inventory reports",
			Long: `Generate various types of inventory reports.
Currently supports low-stock reports with customizable thresholds.`,
			Args: cobra.MinimumNArgs(1),
			Run:  generateReportCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"low-stock", "--", "-5"})

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
		assert.Contains(t, output, "Error: Threshold cannot be negative")
	})

	t.Run("Unknown report type", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "generate-report",
			Short: "Generate inventory reports",
			Long: `Generate various types of inventory reports.
Currently supports low-stock reports with customizable thresholds.`,
			Args: cobra.MinimumNArgs(1),
			Run:  generateReportCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"unknown-report"})

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
		assert.Contains(t, output, "Unknown report type: unknown-report")
	})
}
