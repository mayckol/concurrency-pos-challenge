package auction_test

import (
	"context"
	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuctionAutoClose(t *testing.T) {
	// Set environment variables
	os.Setenv("AUCTION_DURATION", "2s")
	os.Setenv("MONGODB_URL", "mongodb://admin:admin@localhost:27017")
	os.Setenv("MONGODB_DB", "auction_db")

	// Initialize MongoDB connection
	client, err := mongodb.NewMongoDBConnection(context.Background())
	if err != nil {
		logger.Error("Error trying to connect to MongoDB", err)
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	database := client.Client().Database("auction_db")

	repo := auction.NewAuctionRepository(database)

	// Create an auction
	auctionEntity, _ := auction_entity.CreateAuction(
		"Produto Teste",
		"Categoria Teste",
		"Descrição do Produto Teste",
		auction_entity.New,
	)
	err = repo.CreateAuction(context.Background(), auctionEntity)
	assert.Nil(t, err)

	// Wait for the auction to expire
	time.Sleep(3 * time.Second)

	// Check if the auction status has been updated to Completed
	updatedAuction, err := repo.FindAuctionById(context.Background(), auctionEntity.Id)
	assert.Nil(t, err)
	assert.Equal(t, auction_entity.Completed, updatedAuction.Status)
}
