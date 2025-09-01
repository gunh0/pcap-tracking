package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig  `json:"server"`
	Capture CaptureConfig `json:"capture"`
}

// ServerConfig represents the HTTP server configuration
type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

// CaptureConfig represents the packet capture configuration
type CaptureConfig struct {
	Interfaces    []string `json:"interfaces"`
	SnapLength    int32    `json:"snapLength"`
	Promiscuous   bool     `json:"promiscuous"`
	TimeoutMs     int      `json:"timeoutMs"`
	BufferSizeMB  int      `json:"bufferSizeMB"`
	StatsInterval int      `json:"statsInterval"` // seconds
}

// Load loads configuration from a JSON file
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Default returns default configuration
func Default() *Config {
	// Try to auto-detect active interface
	defaultInterface := "en0"
	if activeInterface, err := AutoDetectActiveInterface(); err == nil && activeInterface != "" {
		defaultInterface = activeInterface
	}

	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		Capture: CaptureConfig{
			Interfaces:    []string{defaultInterface},
			SnapLength:    1024,
			Promiscuous:   true,
			TimeoutMs:     5000,
			BufferSizeMB:  100,
			StatsInterval: 5,
		},
	}
}

// Save saves configuration to a JSON file
func (c *Config) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c)
}
