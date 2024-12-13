package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

func TestAuctionCloserWithDatabase(t *testing.T) {
	err := os.Setenv("AUCTION_DURATION", "10s")
	require.NoError(t, err)
	defer os.Unsetenv("AUCTION_DURATION")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/auctions?authSource=admin"))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	collection := client.Database("auctions").Collection("auctions")
	repo := &AuctionRepository{Collection: collection}

	now := time.Now().Unix()
	_, err = collection.InsertMany(ctx, []interface{}{
		bson.M{"_id": uuid.New().String(), "timestamp": now - 3600, "status": auction_entity.Active},
	})
	if err != nil {
		t.Fatalf("Failed to insert documents: %v", err)
	}
	require.NoError(t, err)

	go StartAuctionCloser(ctx, repo)

	time.Sleep(15 * time.Second)

	var result bson.M
	err = collection.FindOne(ctx, bson.M{"timestamp": now - 3600}).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, int32(auction_entity.Completed), result["status"])
}
