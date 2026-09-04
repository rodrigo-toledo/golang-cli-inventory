//go:build integration

package repository

import (
	"context"
	"testing"

	"cli-inventory/internal/models"
	"cli-inventory/internal/testutils"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Setup a single database for the entire integration test suite
	db := testutils.SetupTestDatabase(t)
	queries := testutils.GetTestQueries(db)

	// Create repositories once
	locationRepo := NewLocationRepository(queries)
	productRepo := NewProductRepository(queries)
	stockRepo := NewStockRepository(queries)
	stockMovementRepo := NewStockMovementRepository(queries)

	t.Run("Locations", func(t *testing.T) {
		// Clean the database before this test group runs
		testutils.CleanupTestDatabase(t, db)

		t.Run("Create and Get Location", func(t *testing.T) {
			// Create a location
			createReq := &models.CreateLocationRequest{
				Name: "Test Warehouse",
			}

			created, err := locationRepo.Create(ctx, createReq)
			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, createReq.Name, created.Name)
			assert.NotZero(t, created.CreatedAt)

			// Retrieve the location by name
			retrieved, err := locationRepo.GetByName(ctx, createReq.Name)
			require.NoError(t, err)
			assert.Equal(t, created.ID, retrieved.ID)
			assert.Equal(t, created.Name, retrieved.Name)

			// Retrieve the location by ID
			retrievedByID, err := locationRepo.GetByID(ctx, created.ID)
			require.NoError(t, err)
			assert.Equal(t, created.ID, retrievedByID.ID)
			assert.Equal(t, created.Name, retrievedByID.Name)
		})

		t.Run("Get Non-Existent Location", func(t *testing.T) {
			// Try to get a location that doesn't exist
			location, err := locationRepo.GetByName(ctx, "NONEXISTENT")
			require.NoError(t, err)
			assert.Nil(t, location)

			location, err = locationRepo.GetByID(ctx, 99999)
			require.NoError(t, err)
			assert.Nil(t, location)
		})

		t.Run("Create Duplicate Name", func(t *testing.T) {
			// Use a unique name for each test run
			createReq := &models.CreateLocationRequest{
				Name: "Duplicate Location",
			}

			// Create first location
			_, err := locationRepo.Create(ctx, createReq)
			require.NoError(t, err)

			// Try to create location with same name
			_, err = locationRepo.Create(ctx, createReq)
			assert.Error(t, err)
		})

		t.Run("List Locations", func(t *testing.T) {
			// Start from a clean database so the exact-count assertion below
			// does not depend on rows created by earlier subtests.
			testutils.CleanupTestDatabase(t, db)

			// Create multiple locations
			locations := []*models.CreateLocationRequest{
				{Name: "Warehouse A"},
				{Name: "Warehouse B"},
				{Name: "Store Main"},
			}

			for _, l := range locations {
				_, err := locationRepo.Create(ctx, l)
				require.NoError(t, err)
			}

			// List all locations
			retrieved, err := locationRepo.List(ctx)
			require.NoError(t, err)
			assert.Len(t, retrieved, 3)

			// Verify all locations are present
			names := make(map[string]bool)
			for _, l := range retrieved {
				names[l.Name] = true
			}

			for _, l := range locations {
				assert.True(t, names[l.Name], "Location with name %s should be in the list", l.Name)
			}
		})

		t.Run("Delete Location", func(t *testing.T) {
			// Create a location to be deleted
			createReq := &models.CreateLocationRequest{
				Name: "Location to Delete",
			}
			created, err := locationRepo.Create(ctx, createReq)
			require.NoError(t, err)

			// Delete the location
			err = locationRepo.Delete(ctx, created.ID)
			require.NoError(t, err)

			// Verify it's gone
			retrieved, err := locationRepo.GetByID(ctx, created.ID)
			require.NoError(t, err)
			assert.Nil(t, retrieved)
		})
	})

	t.Run("Products", func(t *testing.T) {
		// Clean the database before this test group runs
		testutils.CleanupTestDatabase(t, db)

		t.Run("Create and Get Product", func(t *testing.T) {
			// Create a product
			createReq := &models.CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			}

			created, err := productRepo.Create(ctx, createReq)
			require.NoError(t, err)
			assert.NotZero(t, created.ID)
			assert.Equal(t, createReq.SKU, created.SKU)
			assert.Equal(t, createReq.Name, created.Name)
			assert.Equal(t, createReq.Description, created.Description)
			assert.Equal(t, createReq.Price, created.Price)
			assert.NotZero(t, created.CreatedAt)

			// Retrieve the product by SKU
			retrieved, err := productRepo.GetBySKU(ctx, createReq.SKU)
			require.NoError(t, err)
			assert.Equal(t, created.ID, retrieved.ID)
			assert.Equal(t, created.SKU, retrieved.SKU)
			assert.Equal(t, created.Name, retrieved.Name)
			assert.Equal(t, created.Description, retrieved.Description)
			assert.Equal(t, createReq.Price, retrieved.Price)

			// Retrieve the product by ID
			retrievedByID, err := productRepo.GetByID(ctx, created.ID)
			require.NoError(t, err)
			assert.Equal(t, created.ID, retrievedByID.ID)
			assert.Equal(t, created.SKU, retrievedByID.SKU)
		})

		t.Run("Get Non-Existent Product", func(t *testing.T) {
			// Try to get a product that doesn't exist
			product, err := productRepo.GetBySKU(ctx, "NONEXISTENT")
			require.NoError(t, err)
			assert.Nil(t, product)

			product, err = productRepo.GetByID(ctx, 99999)
			require.NoError(t, err)
			assert.Nil(t, product)
		})

		t.Run("Create Duplicate SKU", func(t *testing.T) {
			// Use a unique SKU for each test run
			createReq := &models.CreateProductRequest{
				SKU:         "DUPLICATE-SKU",
				Name:        "First Product",
				Description: "First product",
				Price:       10.00,
			}

			// Create first product
			_, err := productRepo.Create(ctx, createReq)
			require.NoError(t, err)

			// Try to create product with same SKU
			_, err = productRepo.Create(ctx, createReq)
			assert.Error(t, err)
		})

		t.Run("List Products", func(t *testing.T) {
			// Start from a clean database so the exact-count assertion below
			// does not depend on rows created by earlier subtests.
			testutils.CleanupTestDatabase(t, db)

			// Create multiple products
			products := []*models.CreateProductRequest{
				{
					SKU:         "LIST001",
					Name:        "Product 1",
					Description: "First product",
					Price:       5.99,
				},
				{
					SKU:         "LIST002",
					Name:        "Product 2",
					Description: "Second product",
					Price:       15.99,
				},
				{
					SKU:         "LIST003",
					Name:        "Product 3",
					Description: "Third product",
					Price:       25.99,
				},
			}

			for _, p := range products {
				_, err := productRepo.Create(ctx, p)
				require.NoError(t, err)
			}

			// List all products
			retrieved, err := productRepo.List(ctx)
			require.NoError(t, err)
			assert.Len(t, retrieved, 3)

			// Verify all products are present
			skus := make(map[string]bool)
			for _, p := range retrieved {
				skus[p.SKU] = true
			}

			for _, p := range products {
				assert.True(t, skus[p.SKU], "Product with SKU %s should be in the list", p.SKU)
			}
		})

		t.Run("Delete Product", func(t *testing.T) {
			// Create a product to be deleted
			createReq := &models.CreateProductRequest{
				SKU:   "DELETE001",
				Name:  "Product to Delete",
				Price: 1.00,
			}
			created, err := productRepo.Create(ctx, createReq)
			require.NoError(t, err)

			// Delete the product
			err = productRepo.Delete(ctx, created.ID)
			require.NoError(t, err)

			// Verify it's gone
			retrieved, err := productRepo.GetByID(ctx, created.ID)
			require.NoError(t, err)
			assert.Nil(t, retrieved)
		})
	})

	t.Run("Stock and Movements", func(t *testing.T) {
		// newStockFixtures resets the database and returns a fresh product and
		// location so that every subtest is independent of execution order.
		newStockFixtures := func(t *testing.T) (*models.Product, *models.Location) {
			t.Helper()
			testutils.CleanupTestDatabase(t, db)

			product, err := productRepo.Create(ctx, &models.CreateProductRequest{
				SKU:  "STOCK-MV-PROD",
				Name: "Stock and Movement Test Product",
			})
			require.NoError(t, err)

			location, err := locationRepo.Create(ctx, &models.CreateLocationRequest{
				Name: "Stock and Movement Test Location",
			})
			require.NoError(t, err)

			return product, location
		}

		t.Run("Add and Get Stock", func(t *testing.T) {
			product, location := newStockFixtures(t)

			// Add stock
			quantity := 50
			stock, err := stockRepo.AddStock(ctx, product.ID, location.ID, quantity)
			require.NoError(t, err)
			assert.Equal(t, product.ID, stock.ProductID)
			assert.Equal(t, location.ID, stock.LocationID)
			assert.Equal(t, quantity, stock.Quantity)

			// Get stock by product and location
			retrieved, err := stockRepo.GetByProductAndLocation(ctx, product.ID, location.ID)
			require.NoError(t, err)
			assert.Equal(t, stock.ID, retrieved.ID)
		})

		t.Run("Add More Stock to Existing", func(t *testing.T) {
			product, location := newStockFixtures(t)

			stockRepo.AddStock(ctx, product.ID, location.ID, 30)
			stock, err := stockRepo.AddStock(ctx, product.ID, location.ID, 20)
			require.NoError(t, err)
			assert.Equal(t, 50, stock.Quantity)
		})

		t.Run("Remove Stock", func(t *testing.T) {
			product, location := newStockFixtures(t)

			stockRepo.AddStock(ctx, product.ID, location.ID, 100)
			stock, err := stockRepo.RemoveStock(ctx, product.ID, location.ID, 30)
			require.NoError(t, err)
			assert.Equal(t, 70, stock.Quantity)
		})

		t.Run("Remove More Stock Than Available", func(t *testing.T) {
			product, location := newStockFixtures(t)

			stockRepo.AddStock(ctx, product.ID, location.ID, 50)
			_, err := stockRepo.RemoveStock(ctx, product.ID, location.ID, 100)
			require.ErrorIs(t, err, pgx.ErrNoRows)

			// Quantity must remain unchanged
			stock, getErr := stockRepo.GetByProductAndLocation(ctx, product.ID, location.ID)
			require.NoError(t, getErr)
			assert.Equal(t, 50, stock.Quantity)
		})

		t.Run("Get Low Stock", func(t *testing.T) {
			// Start from a clean database so the exact-count assertion below
			// does not depend on rows created by earlier subtests.
			testutils.CleanupTestDatabase(t, db)

			// Setup
			p1, _ := productRepo.Create(ctx, &models.CreateProductRequest{SKU: "LOW1"})
			p2, _ := productRepo.Create(ctx, &models.CreateProductRequest{SKU: "LOW2"})
			p3, _ := productRepo.Create(ctx, &models.CreateProductRequest{SKU: "HIGH1"})
			l1, _ := locationRepo.Create(ctx, &models.CreateLocationRequest{Name: "Loc A"})
			l2, _ := locationRepo.Create(ctx, &models.CreateLocationRequest{Name: "Loc B"})
			stockRepo.AddStock(ctx, p1.ID, l1.ID, 5)
			stockRepo.AddStock(ctx, p2.ID, l1.ID, 8)
			stockRepo.AddStock(ctx, p3.ID, l1.ID, 50)
			stockRepo.AddStock(ctx, p1.ID, l2.ID, 15)

			// Test
			lowStock, err := stockRepo.GetLowStock(ctx, 10)
			require.NoError(t, err)
			assert.Len(t, lowStock, 2)
		})

		t.Run("Create and List Stock Movements", func(t *testing.T) {
			product, _ := newStockFixtures(t)

			source, _ := locationRepo.Create(ctx, &models.CreateLocationRequest{Name: "Source"})
			dest, _ := locationRepo.Create(ctx, &models.CreateLocationRequest{Name: "Destination"})
			stockRepo.AddStock(ctx, product.ID, source.ID, 100)

			// Test
			movement, err := stockMovementRepo.Create(ctx, &models.CreateStockMovementRequest{
				ProductID:        product.ID,
				SourceLocationID: source.ID,
				DestLocationID:   dest.ID,
				Quantity:         25,
				MovementType:     "MOVE",
			})
			require.NoError(t, err)
			assert.NotZero(t, movement.ID)
			assert.Equal(t, "MOVE", movement.MovementType)

			movements, err := stockMovementRepo.ListByProductID(ctx, product.ID)
			require.NoError(t, err)
			assert.Len(t, movements, 1)
		})
	})
}
