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
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
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

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) *internal_error.InternalError {
	now := time.Now()
	auctionDuration := getAuctionDuration()
	expirationTime := now.Add(-auctionDuration).Unix()

	filter := bson.M{"timestamp": bson.M{"$lt": expirationTime}, "status": auction_entity.Active}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	result, err := ar.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to close expired auctions", err)
		return internal_error.NewInternalServerError("Error trying to close expired auctions")
	}

	logger.Info(fmt.Sprintf("Closed %d expired auctions", result.ModifiedCount))
	return nil
}

func StartAuctionCloser(ctx context.Context, repo *AuctionRepository) {
	auctionDuration := getAuctionDuration()
	ticker := time.NewTicker(auctionDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Auction closer stopped")
			return
		case <-ticker.C:
			logger.Info("Checking for expired auctions to close")
			err := repo.CloseExpiredAuctions(ctx)
			if err != nil {
				logger.Error("Error closing expired auctions", err)
			}
		}
	}
}

func getAuctionDuration() time.Duration {
	auctionDuration := os.Getenv("AUCTION_DURATION")
	duration, err := time.ParseDuration(auctionDuration)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}
