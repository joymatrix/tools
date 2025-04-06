package utils

import (
	"configs"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedisClient() {
	redisInfo := configs.Config.Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port),
		Password: redisInfo.Passwd, // no password set
		DB:       0,                // use default DB
	})
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("连接Redis出错，错误信息：", err)
		return
	}
	fmt.Println("pong ", pong)
}

func SetToken(ctx context.Context, userId, token string, expireTime int64) error {
	res := redisClient.Set(ctx, userId, token, time.Duration(expireTime)*time.Second)
	return res.Err()
}

func GetToken(ctx context.Context, userId string) (string, error) {
	res, err := redisClient.Get(ctx, userId).Result()
	return res, err
}

func DeleteToken(ctx context.Context, userId string) error {
	res := redisClient.Del(ctx, userId)
	return res.Err()
}

func SetSms(ctx context.Context, phoneNum, randNum string, expireTime int64) error {
	res := redisClient.Set(ctx, phoneNum, randNum, time.Duration(expireTime)*time.Second)
	return res.Err()
}

func GetSms(ctx context.Context, phoneNum string) (string, error) {
	res, err := redisClient.Get(ctx, phoneNum).Result()
	return res, err
}
