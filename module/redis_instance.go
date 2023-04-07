package module

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	RedisDeployTypeCluster    = "CLUSTER"
	RedisDeployTypeStandalone = "STANDALONE"

	EXPIRATION = 10 * time.Minute
)

type RedisInstance struct {
	Name       string   `json:"name"`
	Url        string   `json:"url"`
	ClusterUrl []string `json:"cluster_url"`
	RedisType  string   `json:"redis_type"`
	Password   string   `json:"password"`
}

type SetItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (r *RedisInstance) NewClient(ctx context.Context) (redis.Cmdable, error) {
	var rdb redis.Cmdable
	if r.RedisType == RedisDeployTypeCluster {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    r.ClusterUrl,
			Password: r.Password,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     r.Url,
			Password: r.Password,
			DB:       0,
		})
	}
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return rdb, nil
}
