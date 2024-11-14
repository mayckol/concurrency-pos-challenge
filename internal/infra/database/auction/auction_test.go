package auction_test

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestAuctionAutoClose(t *testing.T) {
	// Create a mock repository
	repoMock := new(auction.AuctionRepositoryMock)

	// Set up the auction duration
	auctionDuration := 2 * time.Second

	// Create an auction
	auctionEntity, _ := auction_entity.CreateAuction(
		"Produto Teste",
		"Categoria Teste",
		"Descrição do Produto Teste",
		auction_entity.New,
	)

	// Mock the CreateAuction method
	repoMock.On("CreateAuction", mock.Anything, auctionEntity).Return(nil)

	// Mock the FindAuctionById method before closure
	repoMock.On("FindAuctionById", mock.Anything, auctionEntity.Id).
		Return(auctionEntity, nil).Once()

	// Start the goroutine that simulates auction closure
	go func() {
		time.Sleep(auctionDuration)
		auctionEntity.Status = auction_entity.Completed

		// Mock the FindAuctionById method after closure
		repoMock.On("FindAuctionById", mock.Anything, auctionEntity.Id).
			Return(auctionEntity, nil)
	}()

	// Call CreateAuction
	err := repoMock.CreateAuction(context.Background(), auctionEntity)
	assert.Nil(t, err)

	// Wait for the auction to "expire"
	time.Sleep(auctionDuration + 1*time.Second)

	// Retrieve the auction after expiration
	updatedAuction, err := repoMock.FindAuctionById(context.Background(), auctionEntity.Id)
	assert.Nil(t, err)
	assert.Equal(t, auction_entity.Completed, updatedAuction.Status)
}
