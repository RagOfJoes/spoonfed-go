package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
)

func (r *mutationResolver) RootMutation(ctx context.Context) (*string, error) {
	logger.Info("[GraphQL] Root Mutation to allow for scema stitching.")
	return nil, nil
}

func (r *queryResolver) RootQuery(ctx context.Context) (*string, error) {
	logger.Info("[GraphQL] Root Query to allow for scema stitching.")
	return nil, nil
}

func (r *subscriptionResolver) RootSubscription(ctx context.Context) (<-chan *string, error) {
	logger.Info("[GraphQL] Root Subscription to allow for scema stitching.")
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
