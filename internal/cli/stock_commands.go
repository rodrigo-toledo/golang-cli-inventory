// Package cli provides the command-line interface for the inventory management system.
package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/spf13/cobra"
)

// addStockCmd represents the add-stock command
var addStockCmd = &cobra.Command{
	Use:   "add-stock",
	Short: "Add stock for a product at a specific location",
	Long: `Add stock quantity for a specific product at a given location.
This will increase the stock level for the product at the specified location.`,
	Args: cobra.ExactArgs(3),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initDatabase(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		productID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Invalid product ID. Please provide a valid number.\n")
			return
		}

		locationID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error: Invalid location ID. Please provide a valid number.\n")
			return
		}

		quantity, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Error: Invalid quantity. Please provide a valid number.\n")
			return
		}

		if quantity <= 0 {
			fmt.Printf("Error: Quantity must be greater than 0.\n")
			return
		}

		req := &models.AddStockRequest{
			ProductID:  productID,
			LocationID: locationID,
			Quantity:   quantity,
		}

		stock, err := stockService.AddStock(context.Background(), req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Stock added successfully!\n")
		fmt.Printf("   Product ID: %d\n", stock.ProductID)
		fmt.Printf("   Location ID: %d\n", stock.LocationID)
		fmt.Printf("   New Quantity: %d\n", stock.Quantity)
	},
	Example: "inventory add-stock 1 1 50",
}

// moveStockCmd represents the move-stock command
var moveStockCmd = &cobra.Command{
	Use:   "move-stock",
	Short: "Move stock between locations",
	Long: `Move a specified quantity of a product from one location to another.
This operation is performed atomically to ensure data consistency.`,
	Args: cobra.ExactArgs(4),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initDatabase(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		productID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error: Invalid product ID. Please provide a valid number.\n")
			return
		}

		fromLocationID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error: Invalid source location ID. Please provide a valid number.\n")
			return
		}

		toLocationID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Printf("Error: Invalid destination location ID. Please provide a valid number.\n")
			return
		}

		quantity, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Printf("Error: Invalid quantity. Please provide a valid number.\n")
			return
		}

		if quantity <= 0 {
			fmt.Printf("Error: Quantity must be greater than 0.\n")
			return
		}

		if fromLocationID == toLocationID {
			fmt.Printf("Error: Source and destination locations cannot be the same.\n")
			return
		}

		req := &models.MoveStockRequest{
			ProductID:      productID,
			FromLocationID: fromLocationID,
			ToLocationID:   toLocationID,
			Quantity:       quantity,
		}

		stock, err := stockService.MoveStock(context.Background(), req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Stock moved successfully!\n")
		fmt.Printf("   Product ID: %d\n", stock.ProductID)
		fmt.Printf("   From Location: %d → To Location: %d\n", fromLocationID, toLocationID)
		fmt.Printf("   Quantity Moved: %d\n", quantity)
		fmt.Printf("   New Quantity at Destination: %d\n", stock.Quantity)
	},
	Example: "inventory move-stock 1 1 2 10",
}

// generateReportCmd represents the generate-report command
var generateReportCmd = &cobra.Command{
	Use:   "generate-report",
	Short: "Generate inventory reports",
	Long: `Generate various types of inventory reports.
Currently supports low-stock reports with customizable thresholds.`,
	Args: cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initDatabase(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		reportType := args[0]

		switch reportType {
		case "low-stock":
			threshold := 10 // Default threshold
			if len(args) > 1 {
				var err error
				threshold, err = strconv.Atoi(args[1])
				if err != nil {
					fmt.Printf("Error: Invalid threshold. Please provide a valid number.\n")
					return
				}
				if threshold < 0 {
					fmt.Printf("Error: Threshold cannot be negative.\n")
					return
				}
			}

			stocks, err := stockService.GetLowStockReport(context.Background(), threshold)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if len(stocks) == 0 {
				fmt.Printf("📊 No products found with stock below threshold %d.\n", threshold)
				return
			}

			fmt.Printf("📊 Low Stock Report (Threshold: %d items)\n", threshold)
			fmt.Printf("%-6s %-12s %-12s %-10s\n", "ID", "Product", "Location", "Quantity")
			fmt.Printf("%-6s %-12s %-12s %-10s\n", "------", "------------", "------------", "----------")

			for _, stock := range stocks {
				fmt.Printf("%-6d %-12d %-12d %-10d\n", stock.ID, stock.ProductID, stock.LocationID, stock.Quantity)
			}

		default:
			fmt.Printf("❌ Unknown report type: %s\n", reportType)
			fmt.Println("Available report types:")
			fmt.Println("  low-stock [threshold] - Show products with stock below threshold")
		}
	},
	Example: "inventory generate-report low-stock 20",
}

// InitStockCommands initializes the stock-related commands with the required service
func InitStockCommands(ss *service.StockService) {
	stockService = ss
}
