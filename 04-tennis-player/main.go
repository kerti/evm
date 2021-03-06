package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kerti/evm/04-tennis-player/config"
	"github.com/kerti/evm/04-tennis-player/database"
	_ "github.com/kerti/evm/04-tennis-player/docs"
	"github.com/kerti/evm/04-tennis-player/handler"
	"github.com/kerti/evm/04-tennis-player/inject"
	"github.com/kerti/evm/04-tennis-player/repository"
	"github.com/kerti/evm/04-tennis-player/server"
	"github.com/kerti/evm/04-tennis-player/service"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

// @title Tennis Player API
// @version 1.0
// @description Submitted as part of Evermos Backend Engineer Assessment
// @license.name MIT

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
	container.RegisterService("containerRepository", new(repository.ContainerMySQLRepo))
	container.RegisterService("playerRepository", new(repository.PlayerMySQLRepo))

	// Prepare containers - services
	container.RegisterService("containerService", new(service.ContainerImpl))
	container.RegisterService("playerService", new(service.PlayerImpl))

	// Prepare containers - handlers
	container.RegisterService("containerHandler", new(handler.ContainerImpl))
	container.RegisterService("healthHandler", new(handler.HealthImpl))
	container.RegisterService("playerHandler", new(handler.PlayerImpl))

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
