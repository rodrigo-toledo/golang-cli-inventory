package repository

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"

	pgtype "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDBTXForProducts is a mock implementation of the DBTX interface
type MockDBTXForProducts struct {
	mock.Mock
}

func (m *MockDBTXForProducts) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsCalled := m.Called(ctx, query, args)
	return argsCalled.Get(0).(pgconn.CommandTag), argsCalled.Error(1)
}

func (m *MockDBTXForProducts) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	argsCalled := m.Called(ctx, query, args)
	return argsCalled.Get(0).(pgx.Rows), argsCalled.Error(1)
}

func (m *MockDBTXForProducts) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	argsCalled := m.Called(ctx, query, args)
	return argsCalled.Get(0).(pgx.Row)
}

// MockRowForProducts is a mock implementation of the pgx.Row interface
type MockRowForProducts struct {
	mock.Mock
}

func (m *MockRowForProducts) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

// MockRowsForProducts is a mock implementation of the pgx.Rows interface
type MockRowsForProducts struct {
	mock.Mock
	currentIndex int
	rows         []map[string]interface{}
}

func (m *MockRowsForProducts) Close() {
	m.Called()
}

func (m *MockRowsForProducts) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRowsForProducts) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRowsForProducts) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRowsForProducts) Next() bool {
	args := m.Called()
	// If we have a specific return value, use it
	if args.Get(0) != nil {
		return args.Get(0).(bool)
	}
	// Otherwise, use our internal state
	if m.currentIndex < len(m.rows) {
		m.currentIndex++
		return true
	}
	return false
}

func (m *MockRowsForProducts) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRowsForProducts) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRowsForProducts) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

func (m *MockRowsForProducts) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

// MockQuerierForProducts is a mock implementation of the db.Querier interface for products
type MockQuerierForProducts struct {
	mock.Mock
}

func (m *MockQuerierForProducts) AddStock(ctx context.Context, arg db.AddStockParams) (db.Stock, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Stock), args.Error(1)
}

func (m *MockQuerierForProducts) CreateLocation(ctx context.Context, name string) (db.Location, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(db.Location), args.Error(1)
}

func (m *MockQuerierForProducts) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Product), args.Error(1)
}

func (m *MockQuerierForProducts) CreateStock(ctx context.Context, arg db.CreateStockParams) (db.Stock, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Stock), args.Error(1)
}

func (m *MockQuerierForProducts) CreateStockMovement(ctx context.Context, arg db.CreateStockMovementParams) (db.StockMovement, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.StockMovement), args.Error(1)
}

func (m *MockQuerierForProducts) DeleteLocation(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerierForProducts) DeleteProduct(ctx context.Context, id int32) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerierForProducts) DeleteStock(ctx context.Context, arg db.DeleteStockParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerierForProducts) GetLocationByID(ctx context.Context, id int32) (db.Location, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Location), args.Error(1)
}

func (m *MockQuerierForProducts) GetLocationByName(ctx context.Context, name string) (db.Location, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(db.Location), args.Error(1)
}

func (m *MockQuerierForProducts) GetProductByID(ctx context.Context, id int32) (db.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.Product), args.Error(1)
}

func (m *MockQuerierForProducts) GetProductBySKU(ctx context.Context, sku string) (db.Product, error) {
	args := m.Called(ctx, sku)
	return args.Get(0).(db.Product), args.Error(1)
}

func (m *MockQuerierForProducts) GetStockByLocation(ctx context.Context, locationID int32) ([]db.Stock, error) {
	args := m.Called(ctx, locationID)
	return args.Get(0).([]db.Stock), args.Error(1)
}

func (m *MockQuerierForProducts) GetStockByProduct(ctx context.Context, productID int32) ([]db.Stock, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]db.Stock), args.Error(1)
}

func (m *MockQuerierForProducts) GetStockByProductAndLocation(ctx context.Context, arg db.GetStockByProductAndLocationParams) (db.Stock, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Stock), args.Error(1)
}

func (m *MockQuerierForProducts) GetStockMovementsByLocation(ctx context.Context, fromLocationID pgtype.Int4) ([]db.StockMovement, error) {
	args := m.Called(ctx, fromLocationID)
	return args.Get(0).([]db.StockMovement), args.Error(1)
}

func (m *MockQuerierForProducts) GetStockMovementsByProduct(ctx context.Context, productID int32) ([]db.StockMovement, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]db.StockMovement), args.Error(1)
}

func (m *MockQuerierForProducts) ListLocations(ctx context.Context) ([]db.Location, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Location), args.Error(1)
}

func (m *MockQuerierForProducts) ListProducts(ctx context.Context) ([]db.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.Product), args.Error(1)
}

func (m *MockQuerierForProducts) ListStockMovements(ctx context.Context) ([]db.StockMovement, error) {
	args := m.Called(ctx)
	return args.Get(0).([]db.StockMovement), args.Error(1)
}

func (m *MockQuerierForProducts) RemoveStock(ctx context.Context, arg db.RemoveStockParams) (db.Stock, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Stock), args.Error(1)
}

func (m *MockQuerierForProducts) UpdateLocation(ctx context.Context, arg db.UpdateLocationParams) (db.Location, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Location), args.Error(1)
}

func (m *MockQuerierForProducts) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Product), args.Error(1)
}

func (m *MockQuerierForProducts) UpdateStock(ctx context.Context, arg db.UpdateStockParams) (db.Stock, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(db.Stock), args.Error(1)
}

// MockQueriesForProducts is a wrapper to make our mock compatible with *db.Queries
type MockQueriesForProducts struct {
	*MockQuerierForProducts
	db.Queries
}

func TestProductRepository_Create(t *testing.T) {
	// Create a pgtype.Numeric with a float64 value
	price1 := pgtype.Numeric{}
	price1.Scan(9.99)

	tests := []struct {
		name          string
		productReq    *models.CreateProductRequest
		mockProduct   db.Product
		mockError     error
		expectedError string
	}{
		{
			name: "successful creation",
			productReq: &models.CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			},
			mockProduct: db.Product{
				ID:          1,
				Sku:         "TEST001",
				Name:        "Test Product",
				Description: pgtype.Text{String: "A test product", Valid: true},
				Price:       price1,
				CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name: "database error",
			productReq: &models.CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			},
			mockProduct:   db.Product{},
			mockError:     errors.New("database error"),
			expectedError: "failed to create product: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForProducts)
			queries := db.New(mockDB)
			repo := NewProductRepository(queries)

			// Set up mock expectations for the database call
			mockRow := new(MockRowForProducts)
			mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "INSERT INTO products")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRow)
			
			// Set up mock expectations for row scanning
			if tt.mockError != nil {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(tt.mockError)
			} else {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
					// Set the values that would be scanned
					*(args.Get(0).(*int32)) = tt.mockProduct.ID
					*(args.Get(1).(*string)) = tt.mockProduct.Sku
					*(args.Get(2).(*string)) = tt.mockProduct.Name
					*(args.Get(3).(*pgtype.Text)) = tt.mockProduct.Description
					*(args.Get(4).(*pgtype.Numeric)) = tt.mockProduct.Price
					*(args.Get(5).(*pgtype.Timestamptz)) = tt.mockProduct.CreatedAt
				})
			}

			// Execute the method
			result, err := repo.Create(context.Background(), tt.productReq)

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tt.mockProduct.ID), result.ID)
				assert.Equal(t, tt.mockProduct.Sku, result.SKU)
				assert.Equal(t, tt.mockProduct.Name, result.Name)
				assert.Equal(t, tt.mockProduct.Description.String, result.Description)
				
				// Convert the price to float64 for comparison
				floatPrice, _ := tt.mockProduct.Price.Float64Value()
				assert.Equal(t, floatPrice.Float64, result.Price)
				assert.Equal(t, tt.mockProduct.CreatedAt.Time, result.CreatedAt)
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

func TestProductRepository_GetBySKU(t *testing.T) {
	// Create a pgtype.Numeric with a float64 value
	price := pgtype.Numeric{}
	price.Scan(9.99)

	tests := []struct {
		name          string
		productSKU    string
		mockProduct   db.Product
		mockError     error
		expectedError string
	}{
		{
			name:       "successful retrieval",
			productSKU: "TEST001",
			mockProduct: db.Product{
				ID:          1,
				Sku:         "TEST001",
				Name:        "Test Product",
				Description: pgtype.Text{String: "A test product", Valid: true},
				Price:       price,
				CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "product not found",
			productSKU:    "NONEXISTENT",
			mockProduct:   db.Product{},
			mockError:     errors.New("product not found"),
			expectedError: "failed to get product by SKU: product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForProducts)
			queries := db.New(mockDB)
			repo := NewProductRepository(queries)

			// Set up mock expectations for the database call
			mockRow := new(MockRowForProducts)
			mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT id, sku, name, description, price, created_at FROM products WHERE sku = $1")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRow)
			
			// Set up mock expectations for row scanning
			if tt.mockError != nil {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(tt.mockError)
			} else {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
					// Set the values that would be scanned
					*(args.Get(0).(*int32)) = tt.mockProduct.ID
					*(args.Get(1).(*string)) = tt.mockProduct.Sku
					*(args.Get(2).(*string)) = tt.mockProduct.Name
					*(args.Get(3).(*pgtype.Text)) = tt.mockProduct.Description
					*(args.Get(4).(*pgtype.Numeric)) = tt.mockProduct.Price
					*(args.Get(5).(*pgtype.Timestamptz)) = tt.mockProduct.CreatedAt
				})
			}

			// Execute the method
			result, err := repo.GetBySKU(context.Background(), tt.productSKU)

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tt.mockProduct.ID), result.ID)
				assert.Equal(t, tt.mockProduct.Sku, result.SKU)
				assert.Equal(t, tt.mockProduct.Name, result.Name)
				assert.Equal(t, tt.mockProduct.Description.String, result.Description)
				
				// Convert the price to float64 for comparison
				floatPrice, _ := tt.mockProduct.Price.Float64Value()
				assert.Equal(t, floatPrice.Float64, result.Price)
				assert.Equal(t, tt.mockProduct.CreatedAt.Time, result.CreatedAt)
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

func TestProductRepository_GetByID(t *testing.T) {
	// Create a pgtype.Numeric with a float64 value
	price := pgtype.Numeric{}
	price.Scan(9.99)

	tests := []struct {
		name          string
		productID     int
		mockProduct   db.Product
		mockError     error
		expectedError string
	}{
		{
			name:      "successful retrieval",
			productID: 1,
			mockProduct: db.Product{
				ID:          1,
				Sku:         "TEST001",
				Name:        "Test Product",
				Description: pgtype.Text{String: "A test product", Valid: true},
				Price:       price,
				CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "product not found",
			productID:     999,
			mockProduct:   db.Product{},
			mockError:     errors.New("product not found"),
			expectedError: "failed to get product by ID: product not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForProducts)
			queries := db.New(mockDB)
			repo := NewProductRepository(queries)

			// Set up mock expectations for the database call
			mockRow := new(MockRowForProducts)
			mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT id, sku, name, description, price, created_at FROM products WHERE id = $1")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRow)
			
			// Set up mock expectations for row scanning
			if tt.mockError != nil {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(tt.mockError)
			} else {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
					// Set the values that would be scanned
					*(args.Get(0).(*int32)) = tt.mockProduct.ID
					*(args.Get(1).(*string)) = tt.mockProduct.Sku
					*(args.Get(2).(*string)) = tt.mockProduct.Name
					*(args.Get(3).(*pgtype.Text)) = tt.mockProduct.Description
					*(args.Get(4).(*pgtype.Numeric)) = tt.mockProduct.Price
					*(args.Get(5).(*pgtype.Timestamptz)) = tt.mockProduct.CreatedAt
				})
			}

			// Execute the method
			result, err := repo.GetByID(context.Background(), tt.productID)

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tt.mockProduct.ID), result.ID)
				assert.Equal(t, tt.mockProduct.Sku, result.SKU)
				assert.Equal(t, tt.mockProduct.Name, result.Name)
				assert.Equal(t, tt.mockProduct.Description.String, result.Description)
				
				// Convert the price to float64 for comparison
				floatPrice, _ := tt.mockProduct.Price.Float64Value()
				assert.Equal(t, floatPrice.Float64, result.Price)
				assert.Equal(t, tt.mockProduct.CreatedAt.Time, result.CreatedAt)
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

func TestProductRepository_List(t *testing.T) {
	// Create pgtype.Numeric values with float64 values
	price1 := pgtype.Numeric{}
	price1.Scan(9.99)

	price2 := pgtype.Numeric{}
	price2.Scan(19.99)

	tests := []struct {
		name           string
		mockProducts   []db.Product
		mockError      error
		expectedError  string
	}{
		{
			name: "successful list",
			mockProducts: []db.Product{
				{
					ID:          1,
					Sku:         "TEST001",
					Name:        "Test Product 1",
					Description: pgtype.Text{String: "A test product 1", Valid: true},
					Price:       price1,
					CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
				{
					ID:          2,
					Sku:         "TEST002",
					Name:        "Test Product 2",
					Description: pgtype.Text{String: "A test product 2", Valid: true},
					Price:       price2,
					CreatedAt:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:           "database error",
			mockProducts:   nil,
			mockError:      errors.New("database error"),
			expectedError:  "failed to list products: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForProducts)
			queries := db.New(mockDB)
			repo := NewProductRepository(queries)

			// Set up mock expectations for the database call
			mockRows := new(MockRowsForProducts)
			mockDB.On("Query", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT id, sku, name, description, price, created_at FROM products")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRows, tt.mockError)
			
			if tt.mockError == nil {
				// Set up mock expectations for rows iteration
				mockRows.On("Next").Return(true).Times(len(tt.mockProducts))
				mockRows.On("Next").Return(false).Once()
				
				// Set up mock expectations for row scanning
				for _, prod := range tt.mockProducts {
					mockRows.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Text"), mock.AnythingOfType("*pgtype.Numeric"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
						// Set the values that would be scanned
						*(args.Get(0).(*int32)) = prod.ID
						*(args.Get(1).(*string)) = prod.Sku
						*(args.Get(2).(*string)) = prod.Name
						*(args.Get(3).(*pgtype.Text)) = prod.Description
						*(args.Get(4).(*pgtype.Numeric)) = prod.Price
						*(args.Get(5).(*pgtype.Timestamptz)) = prod.CreatedAt
					}).Once()
				}
				
				// Set up mock expectations for error checking and closing
				mockRows.On("Err").Return(nil)
				mockRows.On("Close").Return()
			} else {
				// When there's an error, Close should still be called
				mockRows.On("Close").Maybe().Return()
			}

			// Execute the method
			result, err := repo.List(context.Background())

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, len(tt.mockProducts))
				for i, prod := range tt.mockProducts {
					assert.Equal(t, int(prod.ID), result[i].ID)
					assert.Equal(t, prod.Sku, result[i].SKU)
					assert.Equal(t, prod.Name, result[i].Name)
					assert.Equal(t, prod.Description.String, result[i].Description)
					
					// Convert the price to float64 for comparison
					floatPrice, _ := prod.Price.Float64Value()
					assert.Equal(t, floatPrice.Float64, result[i].Price)
					assert.Equal(t, prod.CreatedAt.Time, result[i].CreatedAt)
				}
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRows.AssertExpectations(t)
		})
	}
}