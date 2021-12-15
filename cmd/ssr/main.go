package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/quantonganh/ssr"
	"github.com/quantonganh/ssr/http"
	"github.com/quantonganh/ssr/postgresql"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("http.addr", ":8080")

	var config *ssr.Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	app, err := NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop
		cancel()
	}()

	if err := app.Run(ctx); err != nil {
		_ = app.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	<-ctx.Done()

	if err := app.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type app struct {
	config *ssr.Config
	httpServer *http.Server
}

func NewApp(config *ssr.Config) (*app, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

	db, err := gorm.Open(postgres.Open(psqlConn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&ssr.Repository{}, &ssr.Scan{}); err != nil {
		return nil, err
	}

	return &app{
		config: config,
		httpServer: http.NewServer(
			postgresql.NewRepositoryService(db),
			postgresql.NewScanService(db),
		),
	}, nil
}

func (a *app) Run(ctx context.Context) error {
	a.httpServer.Addr = a.config.HTTP.Addr
	if err := a.httpServer.Open(); err != nil {
		return err
	}
	return nil
}

func (a *app) Close() error {
	if a.httpServer != nil {
		if err := a.httpServer.Close(); err != nil {
			return err
		}
	}

	return nil
}