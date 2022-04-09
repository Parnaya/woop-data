package model

type JsonObject map[string]interface{}

type Tag struct {
	Id   string     `json:"id"`
	Data JsonObject `json:"data"`
}
