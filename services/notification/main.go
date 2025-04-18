package main

import (
	"context"
	"encoding/json"
	"log"
	"trendyol-tracker/pkg/config"
	"trendyol-tracker/pkg/database"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type NotificationService struct {
	config   *config.Config
	db       *gorm.DB
	consumer sarama.Consumer
}

type PriceDropNotification struct {
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	OldPrice  float64 `json:"old_price"`
	NewPrice  float64 `json:"new_price"`
}

func NewNotificationService(cfg *config.Config) (*NotificationService, error) {
	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Configure Kafka consumer
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer([]string{cfg.KafkaBootstrapServers}, config)
	if err != nil {
		return nil, err
	}

	return &NotificationService{
		config:   cfg,
		db:       db,
		consumer: consumer,
	}, nil
}

func (s *NotificationService) startConsumer(ctx context.Context) error {
	// Subscribe to price drop notifications topic
	partitionConsumer, err := s.consumer.ConsumePartition("price-drops", 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var notification PriceDropNotification
			if err := json.Unmarshal(msg.Value, &notification); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			
			// TODO: Implement notification sending logic
			log.Printf("Price drop notification: Product %d price changed from %.2f to %.2f",
				notification.ProductID, notification.OldPrice, notification.NewPrice)

		case <-ctx.Done():
			return nil
		}
	}
}

func main() {
	cfg := config.LoadConfig()

	service, err := NewNotificationService(cfg)
	if err != nil {
		log.Fatalf("Failed to create notification service: %v", err)
	}

	ctx := context.Background()
	if err := service.startConsumer(ctx); err != nil {
		log.Fatalf("Error starting consumer: %v", err)
	}
}
