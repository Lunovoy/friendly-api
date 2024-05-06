package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/lunovoy/friendly/docs"
	"github.com/sirupsen/logrus"

	"github.com/lunovoy/friendly/internal/handler"
	"github.com/lunovoy/friendly/internal/repository"
	"github.com/lunovoy/friendly/internal/server"
	"github.com/lunovoy/friendly/internal/service"
	"github.com/spf13/viper"
)

const uploadDir = "./images/"

// @title Friendly app API
// @version 1.0
// @description API Server for Friendly Application

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
	})

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error connecting postgres DB: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	serverPort := viper.GetString("port")
	if serverPort == "" {
		logrus.Fatalf("error reading server port from config")
	}

	server := server.NewAPIServer(serverPort, handler.InitRoutes())

	go func() {
		if err := server.Run(); err != nil {
			logrus.Fatalf("error starting server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Server shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error while shutting down server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error while closing database connection: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
