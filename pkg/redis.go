package pkg

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init(rdb *redis.Client) {
	Rdb = rdb
}

func SetBlacklistToken(ctx context.Context, token string, exp time.Duration) error {
	return Rdb.Set(ctx, "blacklist:"+token, true, exp).Err()
}

func IsBlacklisted(ctx context.Context, token string) (bool, error) {
	res, err := Rdb.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
