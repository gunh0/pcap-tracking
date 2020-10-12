package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "pcap-tracking/backend/docs"
	"pcap-tracking/backend/internal/api"
	"pcap-tracking/backend/internal/capture"
	"pcap-tracking/backend/internal/config"
	"pcap-tracking/backend/internal/websocket"

	"github.com/sirupsen/logrus"
)

func main() {
	// Parse command line flags
	var configPath = flag.String("config", "config.json", "path to config file")
	var port = flag.String("port", "8080", "server port")
	var interfaces = flag.String("interfaces", "", "network interfaces to monitor (comma-separated)")
	flag.Parse()

	// Setup logging
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Reduce log noise in production
	if os.Getenv("LOG_LEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logrus.WithError(err).Warn("Failed to load config file, using defaults")
		cfg = config.Default()
	}

	// Override config with command line flags
	if *port != "8080" {
		cfg.Server.Port = *port
	}
	if *interfaces != "" {
		cfg.Capture.Interfaces = parseInterfaces(*interfaces)
	} else {
		// Auto-detect active interface if no interface specified
		if len(cfg.Capture.Interfaces) == 0 || (len(cfg.Capture.Interfaces) == 1 && cfg.Capture.Interfaces[0] == "en0") {
			if activeInterface, err := config.AutoDetectActiveInterface(); err == nil {
				cfg.Capture.Interfaces = []string{activeInterface}
				logrus.WithField("detected_interface", activeInterface).Info("Auto-detected active network interface")
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"port":       cfg.Server.Port,
		"interfaces": cfg.Capture.Interfaces,
	}).Info("Starting PCAP tracking server")

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize packet capture manager
	captureManager := capture.NewManager(cfg.Capture, hub)

	// Start capturing on configured interfaces
	if err := captureManager.Start(); err != nil {
		logrus.WithError(err).Fatal("Failed to start packet capture")
	}

	// Initialize and start HTTP server
	server := api.NewServer(cfg.Server, captureManager, hub)
	go func() {
		if err := server.Start(); err != nil {
			logrus.WithError(err).Fatal("Failed to start HTTP server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	logrus.Info("Shutting down server...")

	// Graceful shutdown
	captureManager.Stop()
	hub.Stop()
	server.Stop()

	logrus.Info("Server stopped")
}

func parseInterfaces(interfaces string) []string {
	if interfaces == "" {
		return []string{}
	}

	var result []string
	for _, iface := range splitByComma(interfaces) {
		if iface != "" {
			result = append(result, iface)
		}
	}
	return result
}

func splitByComma(s string) []string {
	var result []string
	var current string

	for _, char := range s {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}

	return result
}
