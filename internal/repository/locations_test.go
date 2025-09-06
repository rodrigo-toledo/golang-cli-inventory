package repository

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDBTXForLocations is a mock implementation of the DBTX interface
type MockDBTXForLocations struct {
	mock.Mock
}

func (m *MockDBTXForLocations) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	argsCalled := m.Called(ctx, query, args)
	return argsCalled.Get(0).(pgconn.CommandTag), argsCalled.Error(1)
}

func (m *MockDBTXForLocations) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	argsCalled := m.Called(ctx, query, args)
	return argsCalled.Get(0).(pgx.Rows), argsCalled.Error(1)
}

func (m *MockDBTXForLocations) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	argsCalled := m.Called(ctx, query, args)
	return argsCalled.Get(0).(pgx.Row)
}

// MockRow is a mock implementation of the pgx.Row interface
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

// MockRows is a mock implementation of the pgx.Rows interface
type MockRows struct {
	mock.Mock
	currentIndex int
	rows         []map[string]interface{}
}

func (m *MockRows) Close() {
	m.Called()
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	args := m.Called()
	return args.Get(0).(pgconn.CommandTag)
}

func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	args := m.Called()
	return args.Get(0).([]pgconn.FieldDescription)
}

func (m *MockRows) Next() bool {
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

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRows) RawValues() [][]byte {
	args := m.Called()
	return args.Get(0).([][]byte)
}

func (m *MockRows) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func TestLocationRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		locationReq   *models.CreateLocationRequest
		mockLocation  db.Location
		mockError     error
		expectedError string
	}{
		{
			name: "successful creation",
			locationReq: &models.CreateLocationRequest{
				Name: "Test Warehouse",
			},
			mockLocation: db.Location{
				ID:        1,
				Name:      "Test Warehouse",
				CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name: "database error",
			locationReq: &models.CreateLocationRequest{
				Name: "Test Warehouse",
			},
			mockLocation:  db.Location{},
			mockError:     errors.New("database error"),
			expectedError: "failed to create location: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForLocations)
			queries := db.New(mockDB)
			repo := NewLocationRepository(queries)

			// Set up mock expectations for the database call
			mockRow := new(MockRow)
			mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "INSERT INTO locations")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRow)
			
			// Set up mock expectations for row scanning
			if tt.mockError != nil {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(tt.mockError)
			} else {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
					// Set the values that would be scanned
					*(args.Get(0).(*int32)) = tt.mockLocation.ID
					*(args.Get(1).(*string)) = tt.mockLocation.Name
					*(args.Get(2).(*pgtype.Timestamptz)) = tt.mockLocation.CreatedAt
				})
			}

			// Execute the method
			result, err := repo.Create(context.Background(), tt.locationReq)

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tt.mockLocation.ID), result.ID)
				assert.Equal(t, tt.mockLocation.Name, result.Name)
				assert.Equal(t, tt.mockLocation.CreatedAt.Time, result.CreatedAt)
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

func TestLocationRepository_GetByName(t *testing.T) {
	tests := []struct {
		name          string
		locationName  string
		mockLocation  db.Location
		mockError     error
		expectedError string
	}{
		{
			name:         "successful retrieval",
			locationName: "Test Warehouse",
			mockLocation: db.Location{
				ID:        1,
				Name:      "Test Warehouse",
				CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "location not found",
			locationName:  "Non-existent Warehouse",
			mockLocation:  db.Location{},
			mockError:     errors.New("no rows in result set"),
			expectedError: "failed to get location by name: no rows in result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForLocations)
			queries := db.New(mockDB)
			repo := NewLocationRepository(queries)

			// Set up mock expectations for the database call
			mockRow := new(MockRow)
			mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT id, name, created_at FROM locations WHERE name = $1")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRow)
			
			// Set up mock expectations for row scanning
			if tt.mockError != nil {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(tt.mockError)
			} else {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
					// Set the values that would be scanned
					*(args.Get(0).(*int32)) = tt.mockLocation.ID
					*(args.Get(1).(*string)) = tt.mockLocation.Name
					*(args.Get(2).(*pgtype.Timestamptz)) = tt.mockLocation.CreatedAt
				})
			}

			// Execute the method
			result, err := repo.GetByName(context.Background(), tt.locationName)

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tt.mockLocation.ID), result.ID)
				assert.Equal(t, tt.mockLocation.Name, result.Name)
				assert.Equal(t, tt.mockLocation.CreatedAt.Time, result.CreatedAt)
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

func TestLocationRepository_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		locationID    int
		mockLocation  db.Location
		mockError     error
		expectedError string
	}{
		{
			name:       "successful retrieval",
			locationID: 1,
			mockLocation: db.Location{
				ID:        1,
				Name:      "Test Warehouse",
				CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "location not found",
			locationID:    999,
			mockLocation:  db.Location{},
			mockError:     errors.New("no rows in result set"),
			expectedError: "failed to get location by ID: no rows in result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForLocations)
			queries := db.New(mockDB)
			repo := NewLocationRepository(queries)

			// Set up mock expectations for the database call
			mockRow := new(MockRow)
			mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT id, name, created_at FROM locations WHERE id = $1")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRow)
			
			// Set up mock expectations for row scanning
			if tt.mockError != nil {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(tt.mockError)
			} else {
				mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
					// Set the values that would be scanned
					*(args.Get(0).(*int32)) = tt.mockLocation.ID
					*(args.Get(1).(*string)) = tt.mockLocation.Name
					*(args.Get(2).(*pgtype.Timestamptz)) = tt.mockLocation.CreatedAt
				})
			}

			// Execute the method
			result, err := repo.GetByID(context.Background(), tt.locationID)

			// Assert results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tt.mockLocation.ID), result.ID)
				assert.Equal(t, tt.mockLocation.Name, result.Name)
				assert.Equal(t, tt.mockLocation.CreatedAt.Time, result.CreatedAt)
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRow.AssertExpectations(t)
		})
	}
}

func TestLocationRepository_List(t *testing.T) {
	tests := []struct {
		name           string
		mockLocations  []db.Location
		mockError      error
		expectedError  string
	}{
		{
			name: "successful retrieval",
			mockLocations: []db.Location{
				{
					ID:        1,
					Name:      "Test Warehouse",
					CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
				{
					ID:        2,
					Name:      "Secondary Warehouse",
					CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
				},
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:           "database error",
			mockLocations:  nil,
			mockError:      errors.New("database error"),
			expectedError:  "failed to list locations: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDBTXForLocations)
			queries := db.New(mockDB)
			repo := NewLocationRepository(queries)

			// Set up mock expectations for the database call
			mockRows := new(MockRows)
			mockDB.On("Query", mock.Anything, mock.MatchedBy(func(query string) bool {
				return strings.Contains(query, "SELECT id, name, created_at FROM locations")
			}), mock.AnythingOfType("[]interface {}")).Return(mockRows, tt.mockError)
			
			if tt.mockError == nil {
				// Set up mock expectations for rows iteration
				mockRows.On("Next").Return(true).Times(len(tt.mockLocations))
				mockRows.On("Next").Return(false).Once()
				
				// Set up mock expectations for row scanning
				for _, loc := range tt.mockLocations {
					mockRows.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
						// Set the values that would be scanned
						*(args.Get(0).(*int32)) = loc.ID
						*(args.Get(1).(*string)) = loc.Name
						*(args.Get(2).(*pgtype.Timestamptz)) = loc.CreatedAt
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
				assert.Len(t, result, len(tt.mockLocations))
				for i, loc := range tt.mockLocations {
					assert.Equal(t, int(loc.ID), result[i].ID)
					assert.Equal(t, loc.Name, result[i].Name)
					assert.Equal(t, loc.CreatedAt.Time, result[i].CreatedAt)
				}
			}

			// Assert that the mock expectations were met
			mockDB.AssertExpectations(t)
			mockRows.AssertExpectations(t)
		})
	}
}