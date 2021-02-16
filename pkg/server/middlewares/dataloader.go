package middlewares

import (
	"log"
	"time"

	"github.com/RagOfJoes/spoonfed-go/internal/database"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/dataloader"
	"github.com/RagOfJoes/spoonfed-go/internal/graphql/model"
	"github.com/RagOfJoes/spoonfed-go/pkg/util"
	"github.com/gin-gonic/gin"
)

// Dataloader sets Dataloaders into Context
func Dataloader() gin.HandlerFunc {
	client, err := database.Client()
	if err != nil {
		log.Panic(err)
	}
	return func(c *gin.Context) {
		// UserLoader
		userLoader := dataloader.NewUserLoader(dataloader.UserLoaderConfig{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(ids []string) ([]*model.User, []error) {
				return client.FindUsersByID(c.Request.Context(), ids)
			},
		})
		util.AddToContext(c, util.ProjectContextKeys.Dataloader, &dataloader.Loaders{
			UserByID: userLoader,
		})
		c.Next()
	}
}
