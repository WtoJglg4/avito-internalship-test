package main

import (
	"context"
	"github/avito/entities"
	"github/avito/pkg/handler"
	"github/avito/pkg/repository"
	"github/avito/pkg/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s\n", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s\n", err.Error())
	}

	dbConfig := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		logrus.Fatalf("error initializing db: %s\n", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	router := handler.NewHandler(services)
	srv := new(entities.Server)

	if _, err := services.Authorization.CreateUser(entities.User{
		Login:    os.Getenv("ADMIN_LOGIN"),
		Password: os.Getenv("ADMIN_PASSWORD"),
		Role:     "admin",
	}); err != nil {
		logrus.Error(err.Error())
	}

	go func() {
		if err := srv.Run(viper.GetString("port"), router.InitRoutes()); err != nil {
			logrus.Fatalf("error while running http server: %s\n", err.Error())
		}
	}()

	logrus.Printf("API started on :%s\n", viper.GetString("port"))
	// logrus.Printf("Admin`s login: %s\n Admin`s password: %s\n", "admin", "admin")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Println("API Shutting Down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
