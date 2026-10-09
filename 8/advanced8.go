package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var defaultConfig = Config{
	Port: 8080,
}

type Config struct {
	Port int
}

func LoadConfig(path string) (*Config, error) {
	cfg := defaultConfig

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &cfg, nil
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("incorrect config line: %q", line)
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "port":
			port, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("incorrect port value: %w", err)
			}
			cfg.Port = port
		}
	}

	return &cfg, nil
}

func main() {
	config, err := LoadConfig("config.txt")
	if err != nil {
		log.Fatalf("load config: %s", err.Error())
	}

	fmt.Printf("Значение порта: %d\n", config.Port)
}
