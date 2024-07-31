package config

import (
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
	DB_CASBIN_DRIVER  string
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
	config.DB_CASBIN_DRIVER = cast.ToString(coalesce("DB_CASBIN_DRIVER", "postgres"))
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
	adapter, err := xormadapter.NewAdapter("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_CASBIN_DRIVER, config.DB_PASSWORD))
	if err != nil {
		logger.Error("Error creating Casbin adapter", "error", err.Error())
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("config/model.conf", adapter)
	if err != nil {
		logger.Error("Error creating Casbin enforcer", "error", err.Error())
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		logger.Error("Error loading Casbin policy", "error", err.Error())
		return nil, err
	}

	policies := [][]string{
		// user user service
		{"user", "/api/users", "GET"}, //ok
		{"user", "/api/users", "PUT"}, //ok
		{"user", "/api/users", "DELETE"}, //ok
		{"user", "/api/users/recommendation", "GET"}, //ok
		{"user", "/api/users/products", "GET"}, //ok
	  
		// admin user service
		{"admin", "/api/users/:id", "GET"}, //ok
		{"admin", "/api/users/:id", "PUT"}, //ok
		{"admin", "/api/users/:id", "DELETE"}, //ok
		{"admin", "/api/users", "POST"}, //ok
		{"admin", "/api/users/products/:id", "GET"}, //ok
		{"admin", "/api/users/list", "GET"}, //ok
	  
		// user product service
		{"user", "/api/media", "POST"}, //ok
		{"user", "/api/orders/:product_id", "POST"}, //ok
		{"user", "/api/basket/:product_id", "POST"},//ok
		{"user", "/api/basket", "GET"}, //ok
		{"user", "/api/basket/:product_id", "DELETE"},//k
		{"user", "/api/products/list", "GET"}, //ok
		{"user", "/api/categories", "GET"}, //ok
		{"user", "/api/reviews/:product_id", "GET"}, //ok
		{"user", "/api/reviews/:product_id", "POST"},//ok
		{"user", "/api/reviews/:id", "PUT"}, //ok
		{"user", "/api/reviews/:id", "DELETE"},//k
	  
		// admin product service
		{"admin", "/api/products/list", "GET"}, //ok
		{"admin", "/api/products/:id", "GET"}, //ok
		{"admin", "/api/products", "POST"}, //ok
		{"admin", "/api/products/:id", "PUT"}, //ok
		{"admin", "/api/products/:id", "DELETE"}, //ok
		{"admin", "/api/categories", "GET"}, //ok
		{"admin", "/api/categories", "POST"}, //ok
		{"admin", "/api/categories/:id", "PUT"}, //ok
		{"admin", "/api/categories/:id", "DELETE"}, //ok
		{"admin", "/api/reviews", "GET"}, //ok
		{"admin", "/api/reviews", "POST"},// ok
		{"admin", "/api/reviews/admin/:id", "PUT"}, //ok
		{"admin", "/api/reviews/admin/:id", "DELETE"},//o
		{"admin", "/api/order/:product_id", "GET"}, //ok
	   }

	_, err = enforcer.AddPolicies(policies)
	if err != nil {
		logger.Error("Error adding Casbin policy", "error", err.Error())
		return nil, err
	}

	err = enforcer.SavePolicy()
	if err != nil {
		logger.Error("Error saving Casbin policy", "error", err.Error())
		return nil, err
	}
	return enforcer, nil
}
