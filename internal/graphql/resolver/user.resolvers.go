package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/internal/auth"
	"github.com/RagOfJoes/spoonfed-go/internal/database"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	client, err := database.Client()
	if err != nil {
		return nil, err
	}
	user, err := auth.GetUserFromContext(ctx, util.ProjectContextKeys.User)
	if user == nil || err != nil {
		return nil, auth.ErrUnauthenticated
	}
	userID, err := primitive.ObjectIDFromHex(user.Sub)
	if err != nil {
		return nil, err
	}
	return client.FindUserByID(ctx, userID)
}

func (r *userResolver) Name(ctx context.Context, obj *model.User) (*model.UserName, error) {
	var sb strings.Builder
	sb.WriteString(*obj.GivenName)
	sb.WriteString(" ")
	sb.WriteString(*obj.FamilyName)
	fullName := sb.String()
	return &model.UserName{
		First: obj.GivenName,
		Last:  obj.FamilyName,
		Full:  &fullName,
	}, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
