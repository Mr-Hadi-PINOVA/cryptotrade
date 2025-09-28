package config

import (
    "fmt"
    "os"
)

// Config contains runtime configuration for the API server.
type Config struct {
    Environment string
    ServerPort string
}

// Load reads configuration values from the environment and applies sensible defaults.
func Load() Config {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    return Config{
        Environment: env,
        ServerPort:  fmt.Sprintf(":%s", port),
    }
}
