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

// addProductCmd represents the add-product command
var addProductCmd = &cobra.Command{
	Use:   "add-product",
	Short: "Add a new product to the inventory",
	Long: `Add a new product to the inventory system with SKU, name, description, and price.
The SKU must be unique across all products.`,
	Args: cobra.ExactArgs(4),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initDatabase(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		sku := args[0]
		name := args[1]
		description := args[2]

		price, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Printf("Error: Invalid price format. Please provide a valid number.\n")
			return
		}

		req := &models.CreateProductRequest{
			SKU:         sku,
			Name:        name,
			Description: description,
			Price:       price,
		}

		product, err := productService.CreateProduct(context.Background(), req)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Product created successfully!\n")
		fmt.Printf("   ID: %d\n", product.ID)
		fmt.Printf("   SKU: %s\n", product.SKU)
		fmt.Printf("   Name: %s\n", product.Name)
		fmt.Printf("   Price: $%.2f\n", product.Price)
	},
	Example: "inventory add-product PROD001 \"Laptop\" \"High-performance laptop\" 1299.99",
}

// findProductCmd represents the find-product command
var findProductCmd = &cobra.Command{
	Use:   "find-product",
	Short: "Find a product by SKU",
	Long: `Search for a product in the inventory using its SKU (Stock Keeping Unit).
This will display all product details if found.`,
	Args: cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initDatabase(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		sku := args[0]

		product, err := productService.GetProductBySKU(context.Background(), sku)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("📦 Product found:\n")
		fmt.Printf("   ID: %d\n", product.ID)
		fmt.Printf("   SKU: %s\n", product.SKU)
		fmt.Printf("   Name: %s\n", product.Name)
		fmt.Printf("   Description: %s\n", product.Description)
		fmt.Printf("   Price: $%.2f\n", product.Price)
		fmt.Printf("   Created: %s\n", product.CreatedAt.Format("2006-01-02 15:04:05"))
	},
	Example: "inventory find-product PROD001",
}

// listProductsCmd represents the list-products command
var listProductsCmd = &cobra.Command{
	Use:   "list-products",
	Short: "List all products in the inventory",
	Long:  `Display a list of all products in the inventory system with their basic information.`,
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := initDatabase(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		products, err := productService.ListProducts(context.Background())
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(products) == 0 {
			fmt.Println("No products found in inventory.")
			return
		}

		fmt.Printf("📋 Products in Inventory (%d items):\n", len(products))
		fmt.Printf("%-6s %-15s %-30s %-10s\n", "ID", "SKU", "Name", "Price")
		fmt.Printf("%-6s %-15s %-30s %-10s\n", "------", "---------------", "------------------------------", "----------")

		for _, product := range products {
			fmt.Printf("%-6d %-15s %-30s $%-9.2f\n", product.ID, product.SKU, product.Name, product.Price)
		}
	},
	Example: "inventory list-products",
}

// InitProductCommands initializes the product-related commands with the required service
func InitProductCommands(ps *service.ProductService) {
	productService = ps
}
