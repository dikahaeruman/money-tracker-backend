package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	controller "money-tracker-backend/src/controller"
	database "money-tracker-backend/src/db"
	middleware "money-tracker-backend/src/middleware"
)

func main() {

	// Load environment variables from .env file
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error loading .env file:", err)
			os.Exit(1) // Exit the program with non-zero status code indicating failure
		}
	}

	// Continue with your application logic
	fmt.Println("Environment variables loaded successfully!")

	database.ConnectDatabase()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed : %v\n", err)
	}

	db := database.GetDB()

	app := gin.Default()

	// Middleware to set jwtKey in context
	jwtKey := middleware.GetJWTKey()

	app.Use(KeyMiddleware(jwtKey))

	mode := os.Getenv("APP_DEBUG")
	fmt.Println("Gin mode:", mode)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "OPIONS", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin", "content-type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app.Use(sentrygin.New(sentrygin.Options{}))

	app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	// Public routes
	app.POST("/api/login", controller.LoginHandler(db))
	app.POST("/api/refresh", controller.RefreshTokenHandler)
	app.POST("/api/logout", controller.LogoutHandler)
	app.POST("/api/create-user", controller.CreateUserHandler(db))
	app.GET("api/verify-token", controller.VerifyTokenHandler(jwtKey)) // New verification endpoint

	// Protected routes
	api := app.Group("/api", middleware.JWTMiddleware(jwtKey))
	{
		api.GET("/users", controller.AllUser(db))
		api.POST("/search", controller.SearchHandler(db))
	}
	app.Run(":8282")
}

func KeyMiddleware(jwtKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(jwtKey) == 0 {
			fmt.Println("JWT_KEY environment variable is not set or is empty in main")
		}
		// fmt.Printf("========= JWT_KEY is set: %v\n", jwtKey)
		c.Set("jwtKey", jwtKey)
		c.Next()
	}
}
