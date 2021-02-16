package middlewares

import (
	"time"

	"github.com/RagOfJoes/spoonfed-go/internal/database"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/dataloader"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// Dataloader sets Dataloaders into Context
func Dataloader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// UserLoader
		userLoader := dataloader.NewUserLoader(dataloader.UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []string) ([]*model.User, []error) {
				client, err := database.Client()
				if err != nil {
					c.Next()
					return nil, []error{database.ErrClientNotInitialized}
				}
				return client.FindUsersByID(ids)
			},
		})
		util.AddToContext(c, util.ProjectContextKeys.Dataloader, &dataloader.Loaders{
			UserByID: userLoader,
		})
		c.Next()
	}
}
