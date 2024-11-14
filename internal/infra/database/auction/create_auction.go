package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection      *mongo.Collection
	auctionDuration time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	repo := &AuctionRepository{
		Collection:      database.Collection("auctions"),
		auctionDuration: getAuctionDuration(),
	}

	go repo.startAuctionCloser()

	return repo
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func getAuctionDuration() time.Duration {
	auctionDurationStr := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDurationStr)
	if err != nil {
		// Default value of 5 minutes if the environment variable is not defined
		return time.Minute * 5
	}
	return duration
}

func (ar *AuctionRepository) startAuctionCloser() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ar.closeExpiredAuctions()
		}
	}
}

func (ar *AuctionRepository) closeExpiredAuctions() {
	ctx := context.Background()
	expiredTimestamp := time.Now().Unix() - int64(ar.auctionDuration.Seconds())

	filter := bson.M{
		"status": auction_entity.Active,
		"timestamp": bson.M{
			"$lte": expiredTimestamp,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status": auction_entity.Completed,
		},
	}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error updating expired auctions to 'Completed'", err)
		return
	}

	if result.ModifiedCount > 0 {
		logger.Info(fmt.Sprintf("Closed %d expired auctions", result.ModifiedCount))
	}
}
