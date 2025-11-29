package main

import (
	"flag"
	"log"
	"os"

	"market-data-insights/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Printf("Failed to load configuration: %v", err)
		log.Printf("\n" +
			"🔑 SETUP REQUIRED:\n" +
			"Please set your Massive.com API key in configs/config.yaml under api.massive_api_key\n" +
			"💡 Get a free API key at: https://www.massive.com/\n" +
			"   1. Sign up for a free account\n" +
			"   2. Go to Dashboard -> API Keys\n" +
			"   3. Copy your API key\n")
		os.Exit(1)
	}

	// Set up Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Load HTML templates
	router.LoadHTMLGlob("web/templates/*")

	// Static files
	router.Static("/static", "web/static")

	// Routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Market Data Insights",
		})
	})

	// API health check
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Market Data Insights API is running",
			})
		})
	}

	// Start server
	log.Printf("Starting server on %s", cfg.GetAddress())
	log.Printf("Dashboard available at: http://%s", cfg.GetAddress())
	log.Printf("API available at: http://%s/api", cfg.GetAddress())

	if err := router.Run(cfg.GetAddress()); err != nil {
		log.Printf("Failed to start server: %v", err)
		os.Exit(1)
	}
}