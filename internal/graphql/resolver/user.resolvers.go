package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/generated"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
)

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
