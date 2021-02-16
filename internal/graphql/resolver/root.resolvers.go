package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
)

func (r *mutationResolver) RootMutation(ctx context.Context) (*string, error) {
	log.Print("Root Mutation to allow for scema stitching.")
	return nil, nil
}

func (r *queryResolver) RootQuery(ctx context.Context) (*string, error) {
	log.Print("Root Query to allow for scema stitching.")
	return nil, nil
}

func (r *subscriptionResolver) RootSubscription(ctx context.Context) (<-chan *string, error) {
	log.Print("Root Subscription to allow for scema stitching.")
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
