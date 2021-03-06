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

func (rs *redisService) GetRandomUsersInSquare(uid int, scenario string) (*[]response.RandomUser, error) {
	redis := global.App.Redis
	key := redis2.GenKey(uid, scenario)
	val, err := redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res, err := redis2.UnSerialize(val)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (rs *redisService) GetRecommendedUsers(uid int, scenario string) (*[]response.RecommendedUser, error) {
	key := redis2.GenKey(uid, scenario)
	val, err := global.App.Redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res, err := redis2.UnSerializeRecommendedUsers(val)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (rs *redisService) SetRecommendedUsers(uid int, scenario string, users *[]response.RecommendedUser) error {
	redis := global.App.Redis
	serialzedUserData, err := redis2.SerializeRecommendedUsers(users)
	if err != nil {
		return err
	}
	err = redis.Set(ctx, redis2.GenKey(uid, scenario), serialzedUserData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
