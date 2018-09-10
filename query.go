package main

import (
	context "context"

	"github.com/julienbreux/gqlgensubtest/schema/graph"
	"github.com/julienbreux/gqlgensubtest/schema/model"
)

type queryResolver struct{}

// Query returns a query resolver
func (rr *rootResolver) Query() graph.QueryResolver {
	return &queryResolver{}
}

// Get returns a fake object
func (qr *queryResolver) Get(ctx context.Context, name string) (*model.ObjectReturned, error) {
	return &model.ObjectReturned{
		Name: name,
	}, nil
}
