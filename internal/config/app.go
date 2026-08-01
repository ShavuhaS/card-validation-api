package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	Port int
}

func Load() (*AppConfig, error) {
	portStr, ok := os.LookupEnv("PORT")
	if !ok {
		portStr = "5000"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse the application port: %v", err)
	}
	return &AppConfig{Port: port}, nil
}
