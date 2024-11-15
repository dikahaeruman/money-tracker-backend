package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"money-tracker-backend/internal/auth"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"

	"money-tracker-backend/internal/controllers"
	"money-tracker-backend/internal/repositories"
	"money-tracker-backend/internal/services"
	"money-tracker-backend/pkg/database"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}(db)

	// if err := runMigrations(db); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	accountController := controllers.NewAccountController(accountService)

	// Set up Gin router
	gin.SetMode("debug")
	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	// Define routes
	api := r.Group("/api")

	api.POST("/auth/login", authController.Login)
	api.POST("/auth/refresh", authController.RefreshToken)
	api.POST("/auth/logout", authController.Logout)
	api.POST("/users", userController.CreateUser)

	api.Use(auth.JWTMiddleware())
	{
		api.GET("/verify", authController.VerifyToken)
		api.GET("/users", userController.GetUser)
		api.POST("/users/search", userController.SearchUser)
		api.POST("/accounts", accountController.CreateAccount)
		api.GET("/accounts/", accountController.GetAccounts)
		api.GET("/accounts/:account_id", accountController.GetAccountByID)
		api.PUT("/accounts/:account_id", accountController.UpdateAccount)
		api.DELETE("/accounts/:account_id", accountController.DeleteAccount)
		// Add other protected routes here
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run("localhost:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	absPath, _ := filepath.Abs("migrations")
	log.Printf("Using migrations folder at: %s", absPath)

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Make sure this path is correct
		"money-tracker-database", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations ran successfully or no changes were needed.")
	return nil
}
