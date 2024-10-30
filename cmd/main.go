package main

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"money-tracker-backend/internal/auth"
	"os"

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

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	accountRepo := repositories.NewAccountRepository(db)

	// Initialize services
	authService := services.NewService(userRepo)
	userService := services.NewUserService(userRepo)
	accountService := services.NewAccountService(accountRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	accountController := controllers.NewAccountController(accountService)

	// Set up Gin router
	gin.SetMode("debug")
	r := gin.Default()

	// Define routes
	r.POST("/api/auth/login", authController.Login)
	r.POST("/api/auth/refresh", authController.RefreshToken)
	r.POST("/api/auth/logout", authController.Logout)
	r.POST("/api/users", userController.CreateUser)
	r.GET("/api/users", userController.GetAllUsers)

	api := r.Group("/api")
	api.Use(auth.JWTMiddleware())
	{
		api.GET("/verify", authController.VerifyToken)
		api.POST("/users/search", userController.SearchUser)
		api.POST("/accounts", accountController.CreateAccount)
		api.GET("/accounts/:user_id", accountController.GetAllAccounts)
		api.GET("/accounts/:user_id/:account_id", accountController.GetAccount)
		// Add other protected routes here
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = r.Run(":" + port)
	if err != nil {
		return
	}
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"money-tracker-database", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
