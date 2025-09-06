// Package cli provides the command-line interface for the inventory management system.
// It uses Cobra for robust command parsing and user-friendly CLI interactions.
package cli

import (
	"fmt"
	"net/http"
	"os"

	"cli-inventory/internal/auth"
	"cli-inventory/internal/database"
	"cli-inventory/internal/db"
	"cli-inventory/internal/handlers"
	"cli-inventory/internal/repository"
	"cli-inventory/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

// serveCmd represents the command to start the HTTP server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long:  `Start the HTTP server to expose the inventory management API.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := initDatabase(); err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}

		// Ensure all services are initialized
		if productService == nil || stockService == nil {
			return fmt.Errorf("services not initialized")
		}

		// Initialize LocationService (it's not global yet, so we initialize it here)
		queries := db.New(database.DB)
		locationRepo := repository.NewLocationRepository(queries)
		locationService := service.NewLocationService(locationRepo)

		// Initialize Auth Handler
		authConfig, err := auth.LoadConfig()
		if err != nil {
			return fmt.Errorf("failed to load auth config: %w", err)
		}
		authHandler, err := auth.NewAuthHandler(authConfig)
		if err != nil {
			return fmt.Errorf("failed to initialize auth handler: %w", err)
		}

		// Initialize handlers
		productHandler := handlers.NewProductHandler(productService)
		locationHandler := handlers.NewLocationHandler(locationService)
		stockHandler := handlers.NewStockHandler(stockService)

		// Setup Chi router
		r := chi.NewRouter()

		// Middleware
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(auth.Authenticator(authHandler.SessionSecret()))

		// Auth Routes (no middleware)
		r.Get("/login", authHandler.LoginHandler)
		r.Get("/callback", authHandler.CallbackHandler)
		r.Get("/logout", authHandler.LogoutHandler)

		// API Routes (protected by AuthMiddleware)
		r.Route("/api/v1", func(r chi.Router) {
			// Product routes
			r.Route("/products", func(r chi.Router) {
				r.Post("/", productHandler.CreateProduct)
				r.Get("/", productHandler.ListProducts)
				r.Get("/{sku}", productHandler.GetProductBySKU)
			})

			// Location routes
			r.Route("/locations", func(r chi.Router) {
				r.Post("/", locationHandler.CreateLocation)
				r.Get("/", locationHandler.ListLocations)
				r.Get("/{name}", locationHandler.GetLocationByName)
			})

			// Stock routes
			r.Route("/stock", func(r chi.Router) {
				r.Post("/add", stockHandler.AddStock)
				r.Post("/move", stockHandler.MoveStock)
				r.Get("/low-stock", stockHandler.GetLowStockReport)
			})
		})

		fmt.Println("Starting server on :8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			return fmt.Errorf("failed to start server: %w", err)
		}
		return nil
	},
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
	rootCmd.AddCommand(serveCmd) // Add the new serve command
}
