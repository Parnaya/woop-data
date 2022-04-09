package operations

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
	"sync"
	"woop-tags/integration/tag"
	"woop-tags/model"
)

func CreateTagsHandler(
	databaseTag *tag.Tags,
) echo.HandlerFunc {
	return func(config echo.Context) error {
		var wg sync.WaitGroup
		wg.Add(1)
		handleCreateTags(config.Request().Body, databaseTag, &wg)
		return config.String(http.StatusOK, "")
	}
}

func handleCreateTags(
	body io.ReadCloser,
	databaseTag *tag.Tags,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	tags := new(map[string]model.JsonObject)
	json.NewDecoder(body).Decode(&tags)
	for id, data := range *tags {
		databaseTag.Create(&model.Tag{Id: id, Data: data})
	}
}

func GetTagsHandler(
	databaseTag *tag.Tags,
) echo.HandlerFunc {
	return func(config echo.Context) error {
		var wg sync.WaitGroup
		wg.Add(1)
		return config.String(http.StatusOK, handleGetTags(config.Request().URL.Query(), databaseTag, &wg))
	}
}

func handleGetTags(
	params url.Values,
	databaseTag *tag.Tags,
	wg *sync.WaitGroup,
) string {
	defer wg.Done()
	tagsIds := params["id"]

	tags := databaseTag.Get(&tagsIds)
	datas := make(map[string]model.JsonObject)

	for _, tag := range tags {
		datas[tag.Id] = tag.Data
	}

	jsonString, err := json.Marshal(datas)
	fmt.Println(err)

	return string(jsonString)
}
