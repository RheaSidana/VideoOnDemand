package main

import (
	"context"
	"vod/initializer"
	"vod/modules/middleware"
	"vod/modules/users"
	videos "vod/modules/video"

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

	videos.Apis(protected)


	r.Run()
}
