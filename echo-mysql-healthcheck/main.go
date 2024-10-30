package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	DB *sql.DB
}

func main() {
	cfg, err := LoadConFig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := InitDB(cfg)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer db.Close()

	app := &App{
		DB: db,
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Wait too long..",
		Timeout:      30 * time.Second,
	}))
	e.Use(middleware.Recover())

	e.GET("/up", app.HealthCheck)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ServerPort)))
}
