package main

import (
	"context"
	"vod/initializer"
	"vod/modules/middleware"
	"vod/modules/users"
	"vod/modules/video/videoMetadata"

	"github.com/gin-gonic/gin"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.ConnectToDB()

	ctx := context.Background()
    initializer.NewRedisClient(ctx)
}

func main() {
	r := gin.Default()

	users.Apis(r)
	protected := middleware.Apis(r)

	videoMetadata.Apis(protected)


	r.Run()
}
