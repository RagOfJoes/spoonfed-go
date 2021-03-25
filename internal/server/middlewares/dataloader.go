package middlewares

import (
	"time"

	"github.com/RagOfJoes/spoonfed-go/internal/graphql/dataloader"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/internal/orm"
	"github.com/RagOfJoes/spoonfed-go/internal/orm/services"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// Dataloader sets Dataloaders into Context
func Dataloader(o *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {
		// UserLoader
		userLoader := dataloader.NewUserLoader(dataloader.UserLoaderConfig{
			MaxBatch: 100,
			Wait:     2 * time.Millisecond,
			Fetch: func(keys []string) ([]*model.User, []error) {
				return services.UserDataloader(c.Request.Context(), o.DB, keys)
			},
		})
		// RecipeLikeLoader
		recipeLikeLoader := dataloader.NewRecipeLikeLoader(dataloader.RecipeLikeLoaderConfig{
			MaxBatch: 100,
			Wait:     2 * time.Millisecond,
			Fetch: func(ids []string) ([]*int64, []error) {
				return services.RecipeLikeDataloader(c.Request.Context(), o.DB, ids)
			},
		})

		loaders := dataloader.Loaders{
			UserByID:       userLoader,
			RecipeLikeByID: recipeLikeLoader,
		}
		util.AddToContext(c, util.ProjectContextKeys.Dataloader, &loaders)
		c.Next()
	}
}
