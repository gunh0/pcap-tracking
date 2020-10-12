package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"pcap-tracking/backend/internal/capture"
	"pcap-tracking/backend/internal/config"
	"pcap-tracking/backend/internal/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server represents the HTTP server
type Server struct {
	config         config.ServerConfig
	captureManager *capture.Manager
	hub            *websocket.Hub
	server         *http.Server
	router         *gin.Engine
}

// NewServer creates a new HTTP server
func NewServer(config config.ServerConfig, captureManager *capture.Manager, hub *websocket.Hub) *Server {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(corsConfig))

	server := &Server{
		config:         config,
		captureManager: captureManager,
		hub:            hub,
		router:         router,
	}

	server.setupRoutes()

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logrus.WithField("addr", addr).Info("Starting HTTP server")
	return s.server.ListenAndServe()
}

// Stop stops the HTTP server
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Swagger documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", s.healthCheck)

		// Network interfaces
		api.GET("/interfaces", s.getInterfaces)

		// Statistics
		api.GET("/stats", s.getStats)

		// WebSocket endpoint
		api.GET("/ws", func(c *gin.Context) {
			s.hub.HandleWebSocket(c)
		})
	}

	// Serve static files (for development)
	s.router.Static("/static", "./static")

	// Default route
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "PCAP Tracking API",
			"version": "1.0.0",
			"status":  "running",
		})
	})
}

// healthCheck returns the health status of the service
// @Summary Health check endpoint
// @Description Check if the API service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Health status"
// @Router /api/v1/health [get]
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(time.Now()).Seconds(), // This should be calculated properly
	})
}

// getInterfaces returns available network interfaces
// @Summary Get network interfaces
// @Description Retrieve list of available network interfaces for packet capture
// @Tags interfaces
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of network interfaces"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/interfaces [get]
func (s *Server) getInterfaces(c *gin.Context) {
	interfaces, err := s.captureManager.GetInterfaces()
	if err != nil {
		logrus.WithError(err).Error("Failed to get network interfaces")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get network interfaces",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"interfaces": interfaces,
	})
}

// getStats returns current packet capture statistics
// @Summary Get packet capture statistics
// @Description Retrieve current packet capture statistics for all monitored interfaces
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} models.SystemStats "Packet capture statistics"
// @Router /api/v1/stats [get]
func (s *Server) getStats(c *gin.Context) {
	stats := s.captureManager.GetStats()
	c.JSON(http.StatusOK, stats)
}
