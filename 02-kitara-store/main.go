package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kerti/evm/02-kitara-store/repository"
	"github.com/kerti/evm/02-kitara-store/service"

	"github.com/kerti/evm/02-kitara-store/config"
	"github.com/kerti/evm/02-kitara-store/database"
	"github.com/kerti/evm/02-kitara-store/handler"
	"github.com/kerti/evm/02-kitara-store/inject"
	"github.com/kerti/evm/02-kitara-store/server"
	"github.com/kerti/evm/02-kitara-store/util/logger"
)

func main() {
	// Register logger
	logger.SetupLoggerAuto("", "")

	// Initialize config
	config.Get()

	// Prepare containers
	container := inject.NewContainer()

	// Prepare containers - database
	var db database.MySQL
	container.RegisterService("mysql", &db)

	// Prepare containers - repositories
	container.RegisterService("inventoryRepository", new(repository.InventoryMySQLRepo))
	container.RegisterService("orderRepository", new(repository.OrderMySQLRepo))

	// Prepare containers - services
	container.RegisterService("orderService", new(service.OrderImpl))

	// Prepare containers - handlers
	container.RegisterService("healthHandler", new(handler.HealthImpl))
	container.RegisterService("orderHandler", new(handler.OrderImpl))

	// Prepare containers - HTTP server
	var s server.Server
	container.RegisterService("server", &s)

	// call this after all dependencies are registered
	if err := container.Ready(); err != nil {
		logger.Fatal("Failed to populate services -- %v", err)
	} else {
		logger.Info("Service registry started successfully.")
	}

	// Handle shutdown
	handleShutdown(container)

	// Run server
	s.Serve()
}

// handle graceful shutdown
func handleShutdown(container inject.ServiceContainer) {
	config := config.Get()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func(ch chan os.Signal) {
		<-ch
		defer os.Exit(0)
		duration := config.Server.ShutdownPeriod
		logger.Info("SIGTERM received. Waiting %v seconds to shutdown", duration.Seconds())
		container.PrepareShutdown()
		time.Sleep(duration)
		logger.Info("Cleaning up resources...")
		container.Shutdown()
		logger.Info("Bye!")
	}(ch)
}
