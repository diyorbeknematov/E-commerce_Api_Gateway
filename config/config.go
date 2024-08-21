package config

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"

	_ "github.com/lib/pq"
)

type Config struct {
	HTTP_PORT         string
	GRPC_USER_PORT    string
	GRPC_PRODUCT_PORT string
	DB_HOST           string
	DB_PORT           string
	DB_USER           string
	DB_PASSWORD       string
	DB_NAME           string
	ACCESS_TOKEN      string
}

func Load() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	config := Config{}

	config.HTTP_PORT = cast.ToString(coalesce("HTTP_PORT", ":8080"))
	config.GRPC_USER_PORT = cast.ToString(coalesce("GRPC_USER_PORT", 50050))
	config.GRPC_PRODUCT_PORT = cast.ToString(coalesce("GRPC_PRODUCT_PORT", 50051))
	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToString(coalesce("DB_PORT", "5432"))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "123321"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "postgres"))
	config.ACCESS_TOKEN = cast.ToString(coalesce("ACCESS_TOKEN", "key_is_really_easy"))

	return config
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if !exists {
		return defaultValue
	}
	return value
}

func CasbinEnforcer(logger *slog.Logger) (*casbin.Enforcer, error) {
	config := Load()
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
	config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASSWORD))
	if err != nil {
		log.Println("Eror connecting to database", "error", err.Error())
		logger.Error("Error connecting to database", "error", err.Error())
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec("DROP DATABASE IF EXISTS casbin;")
	if err != nil {
		log.Println("Error dropping table", "error", err.Error())
        logger.Error("Error dropping table", "error", err.Error())
        return nil, err
	}

	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_NAME, config.DB_PASSWORD)
	fmt.Println(conn)
	adapter, err := xormadapter.NewAdapter("postgres", conn)
	if err != nil {
		log.Println("error creating Casbin adapter", "error", err.Error())
		logger.Error("Error creating Casbin adapter", "error", err.Error())
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("config/model.conf", adapter)
	if err != nil {
		logger.Error("Error creating Casbin enforcer", "error", err.Error())
		log.Println("error creating Casbin enforcer", "error", err.Error())
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		log.Println("error loading Casbin policy", "error", err.Error())
		logger.Error("Error loading Casbin policy", "error", err.Error())
		return nil, err
	}

	policies := [][]string{
		{"user", "/api/users", "GET"},
		{"user", "/api/users", "PUT"},
		{"user", "/api/users", "DELETE"},
		{"user", "/api/users/recommendation", "GET"},
		{"user", "/api/users/products", "GET"},
		{"admin", "/api/users/:id", "GET"},
		{"admin", "/api/users/:id", "PUT"},
		{"admin", "/api/users/:id", "DELETE"},
		{"admin", "/api/users", "POST"},
		{"admin", "/api/users/products/:id", "GET"},
		{"admin", "/api/users/list", "GET"},
		{"user", "/api/products/media", "POST"},
		{"user", "/api/orders/:product_id", "POST"},
		{"user", "/api/basket/:product_id", "POST"},
		{"user", "/api/basket", "GET"},
		{"user", "/api/basket/:product_id", "DELETE"},
		{"user", "/api/products/list", "GET"},
		{"user", "/api/categories", "GET"},
		{"user", "/api/reviews/:product_id", "GET"},
		{"user", "/api/reviews/:product_id", "POST"},
		{"user", "/api/reviews/:id", "PUT"},
		{"user", "/api/reviews/:id", "DELETE"},
		{"admin", "/api/products/list", "GET"},
		{"admin", "/api/products/:id", "GET"},
		{"admin", "/api/products", "POST"},
		{"admin", "/api/products/:id", "PUT"},
		{"admin", "/api/products/:id", "DELETE"},
		{"admin", "/api/categories", "GET"},
		{"admin", "/api/categories", "POST"},
		{"admin", "/api/categories/:id", "PUT"},
		{"admin", "/api/categories/:id", "DELETE"},
		{"admin", "/api/reviews", "GET"},
		{"admin", "/api/reviews", "POST"},
		{"admin", "/api/reviews/admin/:id", "PUT"},
		{"admin", "/api/reviews/admin/:id", "DELETE"},
		{"admin", "/api/order/:product_id", "GET"},
	}

	_, err = enforcer.AddPolicies(policies)
	if err != nil {
		log.Println("error adding Casbin policy", "error", err.Error())
		logger.Error("Error adding Casbin policy", "error", err.Error())
		return nil, err
	}

	err = enforcer.SavePolicy()
	if err != nil {
		log.Println("Error saving Casbin policy", "error", err.Error())
		logger.Error("Error saving Casbin policy", "error", err.Error())
		return nil, err
	}
	return enforcer, nil
}
