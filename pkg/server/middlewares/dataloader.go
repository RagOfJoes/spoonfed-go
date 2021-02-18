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
			Wait:     2 * time.Millisecond,
			Fetch: func(ids []string) ([]*model.User, []error) {
				users, errs := client.FindUsersByID(c.Request.Context(), ids)
				if len(errs) > 0 {
					return nil, errs
				}
				usersArr := []*model.User{}
				usersMap := make(map[string]*model.User)
				for _, user := range users {
					usersMap[user.Sub] = user
				}
				for _, id := range ids {
					usersArr = append(usersArr, usersMap[id])
				}
				return usersArr, nil
			},
		})
		util.AddToContext(c, util.ProjectContextKeys.Dataloader, &dataloader.Loaders{
			UserByID: userLoader,
		})
		c.Next()
	}
}
