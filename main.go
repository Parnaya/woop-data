package main

import (
	"encoding/json"
	"fmt"
	"github.com/couchbase/gocb/v2"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
	"woop-tags/configuration"
	"woop-tags/database"
	"woop-tags/log"
	"woop-tags/model"
)

func main() {
	configuration.ShouldParseViperConfig()
	couchbaseConfig := configuration.ShouldParseCouchbaseConfig()
	cluster := database.ShouldGetCluster(couchbaseConfig)
	if err := cluster.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	e := echo.New()

	var tg = TagsHandler(cluster)

	e.GET("/tags", func(c echo.Context) error {
		var q = c.Request().URL.Query()
		tagsIds := q["id"]

		tags := tg.Get(&tagsIds)
		datas := make(map[string]model.JsonObject)

		for _, tag := range tags {
			datas[tag.Id] = tag.Data
		}

		jsonString, err := json.Marshal(datas)
		fmt.Println(err)

		return c.String(http.StatusOK, string(jsonString))
	})

	e.POST("/tags", func(c echo.Context) error {
		tags := new(map[string]model.JsonObject)
		json.NewDecoder(c.Request().Body).Decode(&tags)
		for id, data := range *tags {
			tg.Create(&model.Tag{Id: id, Data: data})
		}
		return c.String(http.StatusOK, "")
	})

	e.Logger.Fatal(e.Start(":2434"))
}

type DBTags struct {
	Create func(entity *model.Tag)
	Update func(entity *model.Tag)
	Delete func(entity *model.Tag)
	Get    func(params *[]string) []model.Tag
}

func couchbaseGet(cluster *gocb.Cluster) func(idArray *[]string) []model.Tag {
	return func(params *[]string) []model.Tag {

		var sqlBuilder strings.Builder
		sqlBuilder.WriteString("SELECT * FROM data USE KEYS [")

		var paramsLen = len(*params) - 1
		for ind, id := range *params {
			sqlBuilder.WriteString("'")
			sqlBuilder.WriteString(id)
			if ind != paramsLen {
				sqlBuilder.WriteString("', ")
			} else {
				sqlBuilder.WriteString("']")
			}
		}

		sql := sqlBuilder.String()

		rows := log.Proxy(
			cluster.Query(sql, nil),
		).(*gocb.QueryResult)

		var item map[string]model.Tag
		var items []model.Tag

		for rows.Next() {
			rows.Row(&item)
			items = append(items, item["data"])
		}

		return items
	}
}

func couchbaseCreate(collection *gocb.Collection) func(tag *model.Tag) {
	return func(tag *model.Tag) {
		if _, err := collection.Insert(tag.Id, tag, nil); err != nil {
			fmt.Errorf("[Couchbase] Ошибка во время вставки tag: %s", err)
		}
	}
}

func couchbaseUpdate(collection *gocb.Collection) func(entity *model.Tag) {
	return func(entity *model.Tag) {

	}
}

func couchbaseDelete(collection *gocb.Collection) func(entity *model.Tag) {
	return func(entity *model.Tag) {

	}
}

func TagsHandler(cluster *gocb.Cluster) DBTags {
	name := "data"

	bucket := cluster.Bucket(name)
	if err := bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	collection := bucket.DefaultCollection()

	return DBTags{
		Get:    couchbaseGet(cluster),
		Create: couchbaseCreate(collection),
		Update: couchbaseUpdate(collection),
		Delete: couchbaseDelete(collection),
	}
}
