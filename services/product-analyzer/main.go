package main

import (
	"context"
	"log"
	"net"
	"trendyol-tracker/pkg/config"
	"trendyol-tracker/pkg/database"
	pb "trendyol-tracker/proto"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type ProductAnalyzerService struct {
	pb.UnimplementedProductAnalyzerServer
	config   *config.Config
	db       *gorm.DB
	producer sarama.SyncProducer
}

func NewProductAnalyzerService(cfg *config.Config) (*ProductAnalyzerService, error) {
	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Configure Kafka producer
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaBootstrapServers}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &ProductAnalyzerService{
		config:   cfg,
		db:       db,
		producer: producer,
	}, nil
}

func (s *ProductAnalyzerService) AnalyzeProduct(ctx context.Context, req *pb.ProductData) (*pb.AnalysisResult, error) {
	// TODO: Implement product analysis logic
	return &pb.AnalysisResult{
		Status:  "success",
		Message: "Product analyzed",
	}, nil
}

func (s *ProductAnalyzerService) UpdateProduct(ctx context.Context, req *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	// TODO: Implement product update logic
	return &pb.ProductUpdateResponse{
		Success: true,
		Message: "Product updated",
	}, nil
}

func main() {
	cfg := config.LoadConfig()

	service, err := NewProductAnalyzerService(cfg)
	if err != nil {
		log.Fatalf("Failed to create product analyzer service: %v", err)
	}

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductAnalyzerServer(grpcServer, service)

	log.Printf("Starting gRPC server on port %d", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
