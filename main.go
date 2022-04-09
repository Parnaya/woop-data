package main

import (
	"github.com/labstack/echo/v4"
	"time"
	"woop-tags/configuration"
	"woop-tags/database"
	"woop-tags/integration/tag"
	"woop-tags/operations"
)

func main() {
	configuration.ShouldParseViperConfig()
	couchbaseConfig := configuration.ShouldParseCouchbaseConfig()
	cluster := database.ShouldGetCluster(couchbaseConfig)
	if err := cluster.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	var tg = tag.Handlers(cluster)

	e := echo.New()
	e.GET("/tags", operations.GetTagsHandler(&tg))
	e.POST("/tags", operations.CreateTagsHandler(&tg))
	e.Logger.Fatal(e.Start(":2434"))
}
