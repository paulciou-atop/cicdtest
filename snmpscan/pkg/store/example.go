package store

import (
	"context"
	"fmt"
	"nms/snmpscan/api/v1/serviceswatcher"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ExampleClient() (err error) {

	redishost, err := serviceswatcher.GetServiceHost("redis")
	if err != nil {
		return
	}
	fmt.Println("redishost ", redishost)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redishost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	return
	// Output: key value
	// key2 does not exist
}

func ExampleRedisClient() error {
	m := map[string]interface{}{
		"word":   "string",
		"number": 1234,
		"float":  12.334,
		"array":  []string{"1", "2"},
		"map": map[string]interface{}{
			"map.string": "abc",
		},
	}

	rdb := RedisClient()
	r1, err := rdb.Store(CreateID("ss"), "/testing", m)
	if err != nil {
		fmt.Println("Error store data", err)
		return err
	}
	fmt.Println("Store result ", r1)

	r2, err := rdb.Read(r1.Path, r1.Id)
	if err != nil {
		fmt.Println("Error read data", err)
		return err
	}
	fmt.Println("result ", r2.Payload)
	return nil
}
