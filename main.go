package main

import (
	"context"
	"vod/initializer"
	"vod/modules/users"

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

	r.Run()
}
