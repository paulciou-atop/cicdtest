package store

import (
	"encoding/json"
	"fmt"
	"nms/snmpscan/api/v1/serviceswatcher"

	"github.com/go-redis/redis/v8"
)

type redisclient struct {
}

func RedisClient() *redisclient {
	return &redisclient{}
}

func key(path, id string) string {
	return path + "/" + id
}

func redisClient() (*redis.Client, error) {
	redishost, err := serviceswatcher.GetServiceHost("redis")
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redishost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb, nil
}

//Store store data into redis
func (r *redisclient) Store(path string, id string, data JsonObj) (Result, error) {

	rdb, err := redisClient()
	if err != nil {
		return Result{}, fmt.Errorf("Connect redis error %+v", err)
	}
	k := key(path, id)
	v, err := json.Marshal(data)
	if err != nil {
		return Result{}, fmt.Errorf("marshal data error %+v ", err)
	}
	err = rdb.Set(ctx, k, v, 0).Err()
	if err != nil {
		return Result{}, err
	}

	return Result{
		Path: path,
		Id:   id,
		Payload: JsonObj{
			"key":   k,
			"value": data,
		},
	}, nil
}

//Read read data from redis
func (r *redisclient) Read(path, id string) (Result, error) {
	rdb, err := redisClient()
	if err != nil {
		return Result{}, fmt.Errorf("Connect redis error %+v", err)
	}
	k := key(path, id)
	val, err := rdb.Get(ctx, k).Result()
	if err != nil {
		return Result{}, err
	}
	var ret map[string]interface{}
	json.Unmarshal([]byte(val), &ret)
	return Result{
		Path:    path,
		Id:      id,
		Payload: ret,
	}, nil
}
