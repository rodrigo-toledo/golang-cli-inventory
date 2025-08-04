//go:build integration

package repository

import (
	"context"
	"testing"

	"cli-inventory/internal/models"
	"cli-inventory/internal/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	db := testutils.SetupTestDatabase(t)
	defer testutils.TeardownTestDatabase(t)

	// Cleanup database before each test
	testutils.CleanupTestDatabase(t, db)

	// Create queries instance
	queries := testutils.GetTestQueries(db)
	repo := NewProductRepository(queries)

	ctx := context.Background()

	t.Run("Create and Get Product", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Create a product
		createReq := &models.CreateProductRequest{
			SKU:         "TEST001",
			Name:        "Test Product",
			Description: "A test product",
			Price:       9.99,
		}

		created, err := repo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, createReq.SKU, created.SKU)
		assert.Equal(t, createReq.Name, created.Name)
		assert.Equal(t, createReq.Description, created.Description)
		assert.Equal(t, createReq.Price, created.Price)
		assert.NotZero(t, created.CreatedAt)

		// Retrieve the product by SKU
		retrieved, err := repo.GetBySKU(ctx, createReq.SKU)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.SKU, retrieved.SKU)
		assert.Equal(t, created.Name, retrieved.Name)
		assert.Equal(t, created.Description, retrieved.Description)
		assert.Equal(t, created.Price, retrieved.Price)

		// Retrieve the product by ID
		retrievedByID, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrievedByID.ID)
		assert.Equal(t, created.SKU, retrievedByID.SKU)
	})

	t.Run("Get Non-Existent Product", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Try to get a product that doesn't exist
		product, err := repo.GetBySKU(ctx, "NONEXISTENT")
		require.NoError(t, err)
		assert.Nil(t, product)

		product, err = repo.GetByID(ctx, 99999)
		require.NoError(t, err)
		assert.Nil(t, product)
	})

	t.Run("Create Duplicate SKU", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		createReq := &models.CreateProductRequest{
			SKU:         "DUPLICATE001",
			Name:        "First Product",
			Description: "First product",
			Price:       10.00,
		}

		// Create first product
		_, err := repo.Create(ctx, createReq)
		require.NoError(t, err)

		// Try to create product with same SKU
		_, err = repo.Create(ctx, createReq)
		assert.Error(t, err)
	})

	t.Run("List Products", func(t *testing.T) {
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
			_, err := repo.Create(ctx, p)
			require.NoError(t, err)
		}

		// List all products
		retrieved, err := repo.List(ctx)
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
}

func TestLocationRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	db := testutils.SetupTestDatabase(t)
	defer testutils.TeardownTestDatabase(t)

	// Cleanup database before each test
	testutils.CleanupTestDatabase(t, db)

	// Create queries instance
	queries := testutils.GetTestQueries(db)
	repo := NewLocationRepository(queries)

	ctx := context.Background()

	t.Run("Create and Get Location", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Create a location
		createReq := &models.CreateLocationRequest{
			Name: "Test Warehouse",
		}

		created, err := repo.Create(ctx, createReq)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, createReq.Name, created.Name)
		assert.NotZero(t, created.CreatedAt)

		// Retrieve the location by name
		retrieved, err := repo.GetByName(ctx, createReq.Name)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, created.Name, retrieved.Name)

		// Retrieve the location by ID
		retrievedByID, err := repo.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrievedByID.ID)
		assert.Equal(t, created.Name, retrievedByID.Name)
	})

	t.Run("Get Non-Existent Location", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Try to get a location that doesn't exist
		location, err := repo.GetByName(ctx, "NONEXISTENT")
		require.NoError(t, err)
		assert.Nil(t, location)

		location, err = repo.GetByID(ctx, 99999)
		require.NoError(t, err)
		assert.Nil(t, location)
	})

	t.Run("Create Duplicate Name", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		createReq := &models.CreateLocationRequest{
			Name: "Duplicate Location",
		}

		// Create first location
		_, err := repo.Create(ctx, createReq)
		require.NoError(t, err)

		// Try to create location with same name
		_, err = repo.Create(ctx, createReq)
		assert.Error(t, err)
	})

	t.Run("List Locations", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Create multiple locations
		locations := []*models.CreateLocationRequest{
			{Name: "Warehouse A"},
			{Name: "Warehouse B"},
			{Name: "Store Main"},
		}

		for _, l := range locations {
			_, err := repo.Create(ctx, l)
			require.NoError(t, err)
		}

		// List all locations
		retrieved, err := repo.List(ctx)
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
}

func TestStockRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup test database
	db := testutils.SetupTestDatabase(t)
	defer testutils.TeardownTestDatabase(t)

	// Cleanup database before each test
	testutils.CleanupTestDatabase(t, db)

	// Create queries instance
	queries := testutils.GetTestQueries(db)
	productRepo := NewProductRepository(queries)
	locationRepo := NewLocationRepository(queries)
	stockRepo := NewStockRepository(queries)

	ctx := context.Background()

	// Setup test data
	product := &models.CreateProductRequest{
		SKU:         "STOCK001",
		Name:        "Stock Test Product",
		Description: "Product for stock testing",
		Price:       19.99,
	}

	_, err := productRepo.Create(ctx, product)
	require.NoError(t, err)

	location := &models.CreateLocationRequest{
		Name: "Stock Test Location",
	}

	_, err = locationRepo.Create(ctx, location)
	require.NoError(t, err)

	t.Run("Add and Get Stock", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Recreate test data
		createdProduct, err := productRepo.Create(ctx, product)
		require.NoError(t, err)

		createdLocation, err := locationRepo.Create(ctx, location)
		require.NoError(t, err)

		// Add stock
		quantity := 50
		stock, err := stockRepo.AddStock(ctx, createdProduct.ID, createdLocation.ID, quantity)
		require.NoError(t, err)
		assert.Equal(t, createdProduct.ID, stock.ProductID)
		assert.Equal(t, createdLocation.ID, stock.LocationID)
		assert.Equal(t, quantity, stock.Quantity)
		assert.NotZero(t, stock.CreatedAt)
		assert.NotZero(t, stock.UpdatedAt)

		// Get stock by product and location
		retrieved, err := stockRepo.GetByProductAndLocation(ctx, createdProduct.ID, createdLocation.ID)
		require.NoError(t, err)
		assert.Equal(t, stock.ID, retrieved.ID)
		assert.Equal(t, stock.ProductID, retrieved.ProductID)
		assert.Equal(t, stock.LocationID, retrieved.LocationID)
		assert.Equal(t, stock.Quantity, retrieved.Quantity)
	})

	t.Run("Add More Stock to Existing", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Recreate test data
		createdProduct, err := productRepo.Create(ctx, product)
		require.NoError(t, err)

		createdLocation, err := locationRepo.Create(ctx, location)
		require.NoError(t, err)

		// Add initial stock
		stockRepo.AddStock(ctx, createdProduct.ID, createdLocation.ID, 30)

		// Add more stock
		stock, err := stockRepo.AddStock(ctx, createdProduct.ID, createdLocation.ID, 20)
		require.NoError(t, err)
		assert.Equal(t, 50, stock.Quantity) // Should be cumulative
	})

	t.Run("Remove Stock", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Recreate test data
		createdProduct, err := productRepo.Create(ctx, product)
		require.NoError(t, err)

		createdLocation, err := locationRepo.Create(ctx, location)
		require.NoError(t, err)

		// Add initial stock
		stockRepo.AddStock(ctx, createdProduct.ID, createdLocation.ID, 100)

		// Remove some stock
		stock, err := stockRepo.RemoveStock(ctx, createdProduct.ID, createdLocation.ID, 30)
		require.NoError(t, err)
		assert.Equal(t, 70, stock.Quantity)
	})

	t.Run("Remove More Stock Than Available", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Recreate test data
		createdProduct, err := productRepo.Create(ctx, product)
		require.NoError(t, err)

		createdLocation, err := locationRepo.Create(ctx, location)
		require.NoError(t, err)

		// Add initial stock
		stockRepo.AddStock(ctx, createdProduct.ID, createdLocation.ID, 50)

		// Try to remove more stock than available
		stock, err := stockRepo.RemoveStock(ctx, createdProduct.ID, createdLocation.ID, 100)
		require.NoError(t, err)
		assert.Equal(t, 0, stock.Quantity) // Should not go below zero
	})

	t.Run("Get Low Stock", func(t *testing.T) {
		testutils.CleanupTestDatabase(t, db)

		// Create multiple products and locations
		products := []*models.CreateProductRequest{
			{SKU: "LOW1", Name: "Low Stock Product 1", Price: 10.00},
			{SKU: "LOW2", Name: "Low Stock Product 2", Price: 20.00},
			{SKU: "HIGH1", Name: "High Stock Product 1", Price: 30.00},
		}

		locations := []*models.CreateLocationRequest{
			{Name: "Location A"},
			{Name: "Location B"},
		}

		var createdProducts []*models.Product
		for _, p := range products {
			cp, err := productRepo.Create(ctx, p)
			require.NoError(t, err)
			createdProducts = append(createdProducts, cp)
		}

		var createdLocations []*models.Location
		for _, l := range locations {
			cl, err := locationRepo.Create(ctx, l)
			require.NoError(t, err)
			createdLocations = append(createdLocations, cl)
		}

		// Add stock with different quantities
		stockRepo.AddStock(ctx, createdProducts[0].ID, createdLocations[0].ID, 5)  // Low stock
		stockRepo.AddStock(ctx, createdProducts[1].ID, createdLocations[0].ID, 8)  // Low stock
		stockRepo.AddStock(ctx, createdProducts[2].ID, createdLocations[0].ID, 50) // High stock
		stockRepo.AddStock(ctx, createdProducts[0].ID, createdLocations[1].ID, 15) // High stock

		// Get low stock with threshold of 10
		lowStock, err := stockRepo.GetLowStock(ctx, 10)
		require.NoError(t, err)
		assert.Len(t, lowStock, 2) // Should find 2 items with low stock

		// Verify the correct items are returned
		stockMap := make(map[[2]int]int) // key: [productID, locationID], value: quantity
		for _, s := range lowStock {
			stockMap[[2]int{s.ProductID, s.LocationID}] = s.Quantity
		}

		assert.Equal(t, 5, stockMap[[2]int{createdProducts[0].ID, createdLocations[0].ID}])
		assert.Equal(t, 8, stockMap[[2]int{createdProducts[1].ID, createdLocations[0].ID}])
	})
}
