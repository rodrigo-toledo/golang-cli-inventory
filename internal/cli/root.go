// Package cli provides the command-line interface for the inventory management system.
// It uses Cobra for robust command parsing and user-friendly CLI interactions.
package cli

import (
	"fmt"
	"os"

	"cli-inventory/internal/database"
	"cli-inventory/internal/db"
	"cli-inventory/internal/repository"
	"cli-inventory/internal/service"

	"github.com/spf13/cobra"
)

// initDatabase initializes the database connection when needed
func initDatabase() error {
	if database.IsInitialized() {
		return nil
	}

	if err := database.InitDB(); err != nil {
		return err
	}

	// Initialize services after database is connected
	queries := db.New(database.DB)
	InitializeServices(queries)

	return nil
}

// Global service variables
var productService *service.ProductService
var stockService *service.StockService

// InitializeServices initializes all services after database connection
func InitializeServices(queries *db.Queries) {
	// Initialize repositories
	productRepo := repository.NewProductRepository(queries)
	locationRepo := repository.NewLocationRepository(queries)
	stockRepo := repository.NewStockRepository(queries)
	movementRepo := repository.NewStockMovementRepository(queries)

	// Initialize services
	productService = service.NewProductService(productRepo)
	stockService = service.NewStockService(productRepo, locationRepo, stockRepo, movementRepo, database.DB)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "inventory",
	Short: "CLI Inventory Management System",
	Long: `A command-line interface for managing inventory, products, and stock levels.
This application allows you to add products, manage stock, move inventory between locations,
and generate reports.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your command '%s'", err)
		os.Exit(1)
	}
}

// init initializes the root command and adds all subcommands
func init() {
	// Add subcommands
	rootCmd.AddCommand(addProductCmd)
	rootCmd.AddCommand(addStockCmd)
	rootCmd.AddCommand(findProductCmd)
	rootCmd.AddCommand(moveStockCmd)
	rootCmd.AddCommand(generateReportCmd)
	rootCmd.AddCommand(listProductsCmd)
}
