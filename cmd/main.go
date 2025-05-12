package main

import (
	"log"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/awsgw"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	gwservice := awsgw.NewService()
	e.GET("/healthcheck", healthcheck.HealthCheckHandler)
	e.GET("/connect", gwservice.ConntectWebSocket)
	e.POST("/disconnect", gwservice.DisconnectWebSocket)
	e.POST("/subscribe", gwservice.SubscribeChannel)
	e.POST("/send-message", gwservice.SendMessage)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
