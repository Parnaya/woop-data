package database

import (
	"github.com/couchbase/gocb/v2"
	"woop-tags/configuration"
)

func GetCluster(couchbaseConfig *configuration.CouchbaseConfig) (*gocb.Cluster, error) {
	return gocb.Connect(
		couchbaseConfig.Connect,
		gocb.ClusterOptions{
			Username: couchbaseConfig.Auth.Username,
			Password: couchbaseConfig.Auth.Password,
		})
}

func ShouldGetCluster(couchbaseConfig *configuration.CouchbaseConfig) *gocb.Cluster {
	cluster, err := GetCluster(couchbaseConfig)
	if err != nil {
		panic(err)
	}

	return cluster
}
