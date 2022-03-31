package services

import (
	"context"
	"webapp_gin/app/common/response"
	"webapp_gin/global"
	redis2 "webapp_gin/utils/redis"
)

type redisService struct {
}

var RedisService = new(redisService)
var ctx = context.Background()

func (rs *redisService) SetRandomUsersInSquare(uid int, scenario string, users *[]response.RandomUser) error {
	redis := global.App.Redis
	serialzedUserData, err := redis2.Serialize(users)
	if err != nil {
		return err
	}
	err = redis.Set(ctx, redis2.GenKey(uid, scenario), serialzedUserData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
