package tag

import (
	"fmt"
	"github.com/couchbase/gocb/v2"
	"strings"
	"time"
	"woop-tags/log"
	"woop-tags/model"
)

type Tags struct {
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

func Handlers(cluster *gocb.Cluster) Tags {
	name := "data"

	bucket := cluster.Bucket(name)
	if err := bucket.WaitUntilReady(5*time.Second, nil); err != nil {
		panic(err)
	}

	collection := bucket.DefaultCollection()

	return Tags{
		Get:    couchbaseGet(cluster),
		Create: couchbaseCreate(collection),
		Update: couchbaseUpdate(collection),
		Delete: couchbaseDelete(collection),
	}
}
