package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/labstack/echo/v4"
	"trendyol-tracker/pkg/config"
	"trendyol-tracker/services/crawler/crawler"
)

type CrawlerService struct {
	config     *config.Config
	producer   sarama.SyncProducer
	crawler    *crawler.Crawler
	mu         sync.Mutex
	isRunning  bool
}

func NewCrawlerService(cfg *config.Config) (*CrawlerService, error) {
	// Configure Kafka producer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaBootstrapServers}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &CrawlerService{
		config:     cfg,
		producer:   producer,
		crawler:    crawler.NewCrawler(),
		isRunning:  false,
	}, nil
}

func main() {
	cfg := config.LoadConfig()

	service, err := NewCrawlerService(cfg)
	if err != nil {
		log.Fatalf("Failed to create crawler service: %v", err)
	}

	e := echo.New()
	
	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// Crawler endpoints
	e.POST("/crawl", service.handleCrawl)
	e.GET("/status", service.handleStatus)

	// Start the server
	log.Printf("Starting crawler service on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func (s *CrawlerService) handleCrawl(c echo.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return c.JSON(http.StatusConflict, map[string]string{
			"status": "error",
			"message": "Crawler is already running",
		})
	}
	s.isRunning = true
	s.mu.Unlock()

	// Start crawling in a goroutine
	go func() {
		defer func() {
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
		}()

		categories, err := s.crawler.GetCategories()
		if err != nil {
			log.Printf("Error getting categories: %v", err)
			return
		}

		for _, category := range categories {
			// Crawl first 5 pages of each category
			for page := 1; page <= 5; page++ {
				products, err := s.crawler.GetProductsFromCategory(category, page)
				if err != nil {
					log.Printf("Error getting products from category %s page %d: %v", category, page, err)
					continue
				}

				// Send each product to Kafka for analysis
				for _, product := range products {
					productJSON, err := json.Marshal(product)
					if err != nil {
						log.Printf("Error marshaling product: %v", err)
						continue
					}

					_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
						Topic: "raw-products",
						Key:   sarama.StringEncoder(product.ID),
						Value: sarama.ByteEncoder(productJSON),
					})
					if err != nil {
						log.Printf("Error sending product to Kafka: %v", err)
					}
				}

				// Add a small delay between pages
				time.Sleep(time.Second)
			}
		}
	}()

	return c.JSON(http.StatusOK, map[string]string{
		"status": "crawling started",
	})
}

func (s *CrawlerService) handleStatus(c echo.Context) error {
	s.mu.Lock()
	isRunning := s.isRunning
	s.mu.Unlock()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": map[string]interface{}{
			"is_running": isRunning,
		},
	})
}
