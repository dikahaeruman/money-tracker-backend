package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	controller "money-tracker-backend/src/controller"
	database "money-tracker-backend/src/db"
	middleware "money-tracker-backend/src/middleware"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func main() {

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
	mode := os.Getenv("APP_DEBUG")
	fmt.Println("Gin mode:", mode)

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	app.Use(sentrygin.New(sentrygin.Options{}))

	app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	api := app.Group("/api")
	{
		api.POST("/login", controller.LoginHandler(db, jwtKey))
		api.POST("/refresh", controller.RefreshTokenHandler(jwtKey))
		api.POST("/logout", controller.LogoutHandler)
		api.POST("/create-user", controller.CreateUserHandler(db))
		api.GET("/users", middleware.JWTMiddleware(jwtKey), controller.AllUser(db))
		api.POST("/search", middleware.JWTMiddleware(jwtKey), controller.SearchHandler(db))
	}

	app.Run(":8181")
}
