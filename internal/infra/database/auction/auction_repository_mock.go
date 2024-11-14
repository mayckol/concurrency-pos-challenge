package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"github.com/stretchr/testify/mock"
)

type AuctionRepositoryMock struct {
	mock.Mock
}

func (m *AuctionRepositoryMock) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	args := m.Called(ctx, auction)
	if args.Get(0) != nil {
		return args.Get(0).(*internal_error.InternalError)
	}
	return nil
}

func (m *AuctionRepositoryMock) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	args := m.Called(ctx, status, category, productName)
	if args.Get(0) != nil {
		return args.Get(0).([]auction_entity.Auction), nil
	}
	return nil, args.Get(1).(*internal_error.InternalError)
}

func (m *AuctionRepositoryMock) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*auction_entity.Auction), nil
	}
	return nil, args.Get(1).(*internal_error.InternalError)
}
