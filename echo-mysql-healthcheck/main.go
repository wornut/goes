package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	ServerPort string
}

func main() {
	e := echo.New()
	e.GET("/up", up)

	e.Logger.Fatal(e.Start(":8080"))
}

func loadConfig() Config {
	return Config{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func up(c echo.Context) error {
	config := loadConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		c.Logger().Errorf("Error opening db: %v", err)
		return c.String(http.StatusInternalServerError, "Cannot opening db connection")
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Second * 5)

	if err := db.Ping(); err != nil {
		c.Logger().Errorf("Error pinging database: %v", err)
		return c.String(http.StatusInternalServerError, "Database ping failed")
	}

	var result int
	err = db.QueryRow("SELECT 1").Scan(&result)
	if err != nil || result != 1 {
		c.Logger().Errorf("Health check query failed: %v", err)
		return c.String(http.StatusInternalServerError, "Health check query failed")
	}

	return c.String(http.StatusOK, "healthy")
}
