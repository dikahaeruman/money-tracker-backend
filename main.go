package main

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	controller "money-tracker-backend/src/controller"
	database "money-tracker-backend/src/db"
	middleware "money-tracker-backend/src/middleware"
)

var jwtKey = []byte("my_secret_key")

func main() {

	database.ConnectDatabase()
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://aeb1b28f74b00ab978941319e3a675ee@o4507263852740608.ingest.us.sentry.io/4507263980994560",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed : %v\n", err)
	}

	db := database.GetDB()
	app := gin.Default()
	gin.SetMode(gin.DebugMode)
	app.Use(sentrygin.New(sentrygin.Options{}))

	app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello world!")
	})

	app.POST("/login", controller.LoginHandler(db, jwtKey))
	app.POST("/logout", controller.LogoutHandler)
	app.GET("/users", middleware.JWTMiddleware(jwtKey), controller.AllUser(db))
	app.POST("/search", middleware.JWTMiddleware(jwtKey), controller.SearchHandler(db))

	app.Run(":8181")
}
