package cli

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	mocks_service "cli-inventory/internal/mocks/service"
	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddProductCmd(t *testing.T) {
	// Save original productService
	originalProductService := productService
	defer func() {
		productService = originalProductService
	}()

	// Create mock repositories and service
	mockProductRepo := mocks_service.NewMockProductRepositoryInterface(t)
	productService = service.NewProductService(mockProductRepo)

	t.Run("Successful product creation", func(t *testing.T) {
		expectedProduct := &models.Product{
			ID:          1,
			SKU:         "TEST001",
			Name:        "Test Product",
			Description: "A test product",
			Price:       99.99,
		}

		// Mock the GetBySKU call to return an error (product not found)
		mockProductRepo.EXPECT().GetBySKU(mock.Anything, "TEST001").Return((*models.Product)(nil), errors.New("product not found"))
		
		// Mock the Create call
		mockProductRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(req *models.CreateProductRequest) bool {
			return req.SKU == "TEST001" && req.Name == "Test Product" && req.Description == "A test product" && req.Price == 99.99
		})).Return(expectedProduct, nil)

		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "add-product",
			Short: "Add a new product to the inventory",
			Long: `Add a new product to the inventory system with SKU, name, description, and price.
The SKU must be unique across all products.`,
			Args: cobra.ExactArgs(4),
			Run:  addProductCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"TEST001", "Test Product", "A test product", "99.99"})

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
		assert.Contains(t, output, "Product created successfully")
		assert.Contains(t, output, "ID: 1")
		assert.Contains(t, output, "SKU: TEST001")
		assert.Contains(t, output, "Name: Test Product")
		assert.Contains(t, output, "$99.99")
	})

	t.Run("Invalid price format", func(t *testing.T) {
		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "add-product",
			Short: "Add a new product to the inventory",
			Long: `Add a new product to the inventory system with SKU, name, description, and price.
The SKU must be unique across all products.`,
			Args: cobra.ExactArgs(4),
			Run:  addProductCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"TEST001", "Test Product", "A test product", "invalid"})

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
		assert.Contains(t, output, "Error: Invalid price format")
	})
}

func TestFindProductCmd(t *testing.T) {
	// Save original productService
	originalProductService := productService
	defer func() {
		productService = originalProductService
	}()

	// Create mock repositories and service
	mockProductRepo := mocks_service.NewMockProductRepositoryInterface(t)
	productService = service.NewProductService(mockProductRepo)

	t.Run("Successful product retrieval", func(t *testing.T) {
		expectedProduct := &models.Product{
			ID:          1,
			SKU:         "TEST001",
			Name:        "Test Product",
			Description: "A test product",
			Price:       99.99,
		}

		mockProductRepo.EXPECT().GetBySKU(mock.Anything, "TEST001").Return(expectedProduct, nil)

		// Create a test command with the same Run function as the original
		testCmd := &cobra.Command{
			Use:   "find-product",
			Short: "Find a product by SKU",
			Long: `Search for a product in the inventory using its SKU (Stock Keeping Unit).
This will display all product details if found.`,
			Args: cobra.ExactArgs(1),
			Run:  findProductCmd.Run, // Use the original Run function
		}
		testCmd.SetArgs([]string{"TEST001"})

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
		assert.Contains(t, output, "Product found")
		assert.Contains(t, output, "ID: 1")
		assert.Contains(t, output, "SKU: TEST001")
		assert.Contains(t, output, "Name: Test Product")
		assert.Contains(t, output, "Description: A test product")
		assert.Contains(t, output, "$99.99")
	})
}