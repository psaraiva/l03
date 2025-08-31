package auction_entity

import (
	"context"
	"l03/internal/internal_error"
	"time"

	"github.com/google/uuid"
)

type Auction struct {
	ID          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota + 1
	Running
	Completed
)

const (
	New ProductCondition = iota + 1
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	ChangeStatus(ctx context.Context, id string, status AuctionStatus) *internal_error.InternalError
	Create(ctx context.Context, auction Auction) *internal_error.InternalError
	FindAuctions(ctx context.Context, status AuctionStatus, category, productName string) (*[]Auction, *internal_error.InternalError)
	FindById(ctx context.Context, id string) (*Auction, *internal_error.InternalError)
}

func Create(productname, category, description string, condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		ID:          uuid.New().String(),
		ProductName: productname,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validade(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (a *Auction) Validade() *internal_error.InternalError {
	if len(a.ProductName) < 5 {
		return internal_error.NewBadRequestError("Invalid product name")
	}

	if len(a.Category) < 5 {
		return internal_error.NewBadRequestError("Invalid category")
	}

	if len(a.Description) < 5 {
		return internal_error.NewBadRequestError("Invalid description")
	}

	if a.Condition != New && a.Condition != Used && a.Condition != Refurbished {
		return internal_error.NewBadRequestError("Invalid condition")
	}

	return nil
}
