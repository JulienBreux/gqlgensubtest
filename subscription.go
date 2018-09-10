package main

import (
	"context"
	"fmt"
	"time"

	"github.com/julienbreux/gqlgensubtest/schema/graph"
	"github.com/julienbreux/gqlgensubtest/schema/model"
)

type subscriptionResolver struct{}

// Subscription returns a subscription resolver
func (rr *rootResolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{}
}

// Updated returns a fake object each 5 seconds
func (sr *subscriptionResolver) Updated(ctx context.Context) (<-chan model.ObjectReturned, error) {
	object := make(chan model.ObjectReturned, 1)

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for t := range ticker.C {
			object <- model.ObjectReturned{
				Name: fmt.Sprintf("%s", t),
			}
		}
	}()

	// Stop ticker when client disconnect
	go func() {
		<-ctx.Done()
		ticker.Stop()
	}()

	return object, nil
}
