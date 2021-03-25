package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/models"
	"github.com/RagOfJoes/spoonfed-go/pkg/auth"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
)

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user, err := auth.GetUserFromContext(ctx, util.ProjectContextKeys.User)
	if err != nil {
		return nil, err
	}
	u, ok := user.(*models.User)
	if uOk, err := u.IsValid(false); !ok || !uOk || err != nil {
		return nil, auth.ErrUnauthenticated
	}
	built, err := model.BuildUser(u)
	if err != nil {
		return nil, err
	}
	return built, nil
}
