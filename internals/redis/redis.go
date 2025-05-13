package redis

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-redis/redis/v8"
)

// errors
var (
	ErrUnsupportDataType = errors.New("unsupported data type")
	ErrInputNotSlice     = errors.New("input is not a slice")
	ErrInputEmpty        = errors.New("input is Nil")
)

type RedisHandler struct {
	RedisClient *redis.Client
}
type RedisLimitHandler struct {
	RedisClient *redis.Client
}

func NewRedisConnection() (*RedisHandler, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")
	index := os.Getenv("REDIS_INDEX")
	keyPrefix := os.Getenv("REDIS_KEY_PREFIX")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Username: username,
		Password: password,
	})
	ctx := redisClient.Context()

	//check and set index
	_, err := redisClient.Do(ctx, "FT.INFO", index).Result()
	if err != nil {
		userIdQuery := fmt.Sprintf("FT.CREATE %s ON JSON PREFIX 1 %s: SCHEMA $.user_id AS user_id TAG CASESENSITIVE", index, keyPrefix)
		argS := strings.Split(userIdQuery, " ")
		argI, _ := SliceToInterface(argS)
		redisClient.Do(ctx, argI...)
	}

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil || pong == "" {
		return &RedisHandler{}, err
	}

	return &RedisHandler{
		redisClient,
	}, nil
}

func SliceToInterface(inp interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(inp)
	if s.Kind() != reflect.Slice {
		return nil, ErrInputNotSlice
	}
	if s.IsNil() {
		return nil, ErrInputEmpty
	}

	res := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		res[i] = s.Index(i).Interface()
	}

	return res, nil
}

// set json data to redis
func (r *RedisHandler) SetJSONData(c context.Context, key, path string, data []byte) error {
	_, err := r.RedisClient.Do(c, "JSON.SET", key, ".", data).Result()
	return err
}

// get json data from redis
func (r *RedisHandler) GetJSONData(c context.Context, key, path string) (interface{}, error) {
	reply, err := r.RedisClient.Do(c, "JSON.GET", key, path).Result()
	if err != nil {
		return nil, err
	}
	return reply, nil
}
