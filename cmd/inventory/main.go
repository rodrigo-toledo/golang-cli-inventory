package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rodrigotoledo/cli-inventory/internal/database"
	"github.com/rodrigotoledo/cli-inventory/internal/db"
	"github.com/rodrigotoledo/cli-inventory/internal/models"
	"github.com/rodrigotoledo/cli-inventory/internal/repository"
	"github.com/rodrigotoledo/cli-inventory/internal/service"
)

func main() {
	// Initialize database connection
	database.InitDB()
	defer database.DB.Close()

	// Initialize queries
	queries := db.New(database.DB)

	// Initialize repositories
	productRepo := repository.NewProductRepository(queries)
	locationRepo := repository.NewLocationRepository(queries)
	stockRepo := repository.NewStockRepository(queries)
	movementRepo := repository.NewStockMovementRepository(queries)

	// Initialize services
	productService := service.NewProductService(productRepo)
	// locationService := service.NewLocationService(locationRepo)  // Commented out as it's not currently used
	stockService := service.NewStockService(productRepo, locationRepo, stockRepo, movementRepo, database.DB)

	// Create a context
	ctx := context.Background()

	// Simple CLI interface
	if len(os.Args) < 2 {
		fmt.Println("Usage: inventory [command]")
		fmt.Println("Available commands: add-product, add-stock, find-product, move-stock, generate-report")
		return
	}

	command := os.Args[1]

	switch command {
	case "add-product":
		if len(os.Args) < 6 {
			fmt.Println("Usage: inventory add-product <sku> <name> <description> <price>")
			return
		}
		sku := os.Args[2]
		name := os.Args[3]
		description := os.Args[4]
		price := 0.0
		fmt.Sscanf(os.Args[5], "%f", &price)

		req := &models.CreateProductRequest{
			SKU:         sku,
			Name:        name,
			Description: description,
			Price:       price,
		}

		product, err := productService.CreateProduct(ctx, req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Product created successfully with ID: %d\n", product.ID)

	case "add-stock":
		if len(os.Args) < 5 {
			fmt.Println("Usage: inventory add-stock <product-id> <location-id> <quantity>")
			return
		}
		productID := 0
		locationID := 0
		quantity := 0
		fmt.Sscanf(os.Args[2], "%d", &productID)
		fmt.Sscanf(os.Args[3], "%d", &locationID)
		fmt.Sscanf(os.Args[4], "%d", &quantity)

		req := &models.AddStockRequest{
			ProductID:  productID,
			LocationID: locationID,
			Quantity:   quantity,
		}

		stock, err := stockService.AddStock(ctx, req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Stock added successfully. New quantity: %d\n", stock.Quantity)

	case "find-product":
		if len(os.Args) < 3 {
			fmt.Println("Usage: inventory find-product <sku>")
			return
		}
		sku := os.Args[2]

		product, err := productService.GetProductBySKU(ctx, sku)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Product found:\nID: %d\nSKU: %s\nName: %s\nDescription: %s\nPrice: %.2f\n",
			product.ID, product.SKU, product.Name, product.Description, product.Price)

	case "move-stock":
		if len(os.Args) < 6 {
			fmt.Println("Usage: inventory move-stock <product-id> <from-location-id> <to-location-id> <quantity>")
			return
		}
		productID := 0
		fromLocationID := 0
		toLocationID := 0
		quantity := 0
		fmt.Sscanf(os.Args[2], "%d", &productID)
		fmt.Sscanf(os.Args[3], "%d", &fromLocationID)
		fmt.Sscanf(os.Args[4], "%d", &toLocationID)
		fmt.Sscanf(os.Args[5], "%d", &quantity)

		req := &models.MoveStockRequest{
			ProductID:      productID,
			FromLocationID: fromLocationID,
			ToLocationID:   toLocationID,
			Quantity:       quantity,
		}

		stock, err := stockService.MoveStock(ctx, req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Stock moved successfully. New quantity at destination: %d\n", stock.Quantity)

	case "generate-report":
		if len(os.Args) < 3 {
			fmt.Println("Usage: inventory generate-report low-stock <threshold>")
			return
		}
		reportType := os.Args[2]

		switch reportType {
		case "low-stock":
			threshold := 10
			if len(os.Args) >= 4 {
				fmt.Sscanf(os.Args[3], "%d", &threshold)
			}

			stocks, err := stockService.GetLowStockReport(ctx, threshold)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("Low Stock Report (Threshold: %d):\n", threshold)
			fmt.Printf("ID\tProductID\tLocationID\tQuantity\n")
			for _, stock := range stocks {
				fmt.Printf("%d\t%d\t\t%d\t\t%d\n", stock.ID, stock.ProductID, stock.LocationID, stock.Quantity)
			}

		default:
			fmt.Printf("Unknown report type: %s\n", reportType)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: add-product, add-stock, find-product, move-stock, generate-report")
	}
}
