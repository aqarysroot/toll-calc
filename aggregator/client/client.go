package client

import (
	"context"
	"toll-calculator/types"
)

type Client interface {
	Aggregate(ctx context.Context, request *types.AggregateRequest) error
	GetInvoice(ctx context.Context, int2 int) (*types.Invoice, error)
}
