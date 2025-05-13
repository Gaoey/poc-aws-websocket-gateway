package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/auth"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/aws"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/internals/rabbitmq"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/awsgw"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/channels"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/example"
	"github.com/Gaoey/poc-aws-websocket-gateway.git/services/healthcheck"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	rabbitmqClient, err := rabbitmq.NewClient(rabbitmq.Config{
		URL:          os.Getenv("RABBITMQ_URL"),
		ExchangeName: "ws_events",
		ExchangeType: "topic",
	})
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ client: %v", err)
	}

	authApi := auth.NewAuthAPI(os.Getenv("AUTH_API_URL"))
	awsApp := aws.NewAWSApplication()
	gwservice := awsgw.NewService(awsApp, authApi)
	exampleHandler := example.NewExampleHandler(rabbitmqClient)

	e.GET("/healthcheck", healthcheck.HealthCheckHandler)
	e.GET("/connect", gwservice.ConntectWebSocket)
	e.GET("/auth", gwservice.ConntectWebSocket)
	e.POST("/disconnect", gwservice.DisconnectWebSocket)
	e.POST("/subscribe", gwservice.SubscribeChannel)
	e.POST("/send-message", gwservice.SendMessage)

	e.POST("/publish", exampleHandler.PublishMessage)
	e.POST("/sign", exampleHandler.GetSignature)

	// Consumer
	ou := channels.NewChannel(
		rabbitmqClient,
		awsApp,
		"order_update",
		"ws.order.update",
		[]string{"ws.order.update"},
	)
	if err := ou.StartOrderUpdateChannel(); err != nil {
		log.Fatalf("Failed to start order_update consumer: %v", err)
	}

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// When shutting down, stop the channel properly
	log.Println("Stopping WebSocket channels...")
	ou.Stop()
	// Close RabbitMQ connections
	log.Println("Closing RabbitMQ connections...")
	rabbitmqClient.Close()

	log.Println("Server exited properly")
}
