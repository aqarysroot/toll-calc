package main

import (
	"context"
	"toll-calculator/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

// transport Layer
// JSON -> types.Distance
// GRPC -> types.AggregateRequest -> type.Distance
// Webpack => types.Webpack -> types.Distance

// business layer -> business Layer type (main type everyone needs to convert)

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}

	return &types.None{}, s.svc.AggregateDistance(distance)
}
