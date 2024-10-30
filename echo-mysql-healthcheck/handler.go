package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (app *App) HealthCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err := app.DB.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		c.Logger().Errorf("Health check query failed", err)
		return c.String(http.StatusInternalServerError, "Health check query failed")
	}

	if result != 1 {
		c.Logger().Errorf("Health check query returned unexpected result", err)
		return c.String(http.StatusInternalServerError, "Health check query returned unexpected result")
	}

	return c.String(http.StatusOK, "healthy")
}
