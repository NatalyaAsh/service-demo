package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"service-demo/internal/config"
	"strconv"
	"time"

	modeldb "service-demo/internal/models"

	redis "github.com/redis/go-redis/v9"
)

var (
	db  *redis.Client
	ttl int
)

func Init(cfg *config.Config) error {
	ttl = cfg.RDS.TTL
	db = redis.NewClient(&redis.Options{
		Addr: cfg.RDS.Addr,
		// Password:     cfg.RDS.Password,
		DB: cfg.RDS.DB,
		// Username:     cfg.RDS.User,
		MaxRetries:   cfg.RDS.MaxRetries,
		DialTimeout:  time.Duration(cfg.RDS.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.RDS.Timeout) * time.Second,
		WriteTimeout: time.Duration(cfg.RDS.Timeout) * time.Second,
	})

	if err := db.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis server: %s", err.Error())
	}
	slog.Info("Connect to Redis")
	return nil
}

func CloseDB() {
	db.Close()
}

func Set(good *modeldb.Goods) error {
	// Переводим структуру в json
	slog.Info("Redis: Set", "good", good.ID)
	value, err := json.Marshal(&good)
	if err != nil {
		return err
	}

	key := strconv.Itoa(good.ID)
	if err := db.Set(context.Background(), key, value, time.Duration(ttl)*time.Second).Err(); err != nil {
		slog.Error("Redis Set: failed to set data, error:", "err", err.Error())
		return fmt.Errorf("redis: failed to set data, error: %v", err)
	}
	return nil
}

func Get(id string) (modeldb.Goods, error) {
	val, err := db.Get(context.Background(), id).Result()

	if err == redis.Nil {
		slog.Info("Redis Get: value not found")
		return modeldb.Goods{}, fmt.Errorf("redis: value not found")
	} else if err != nil {
		slog.Error("Redis Get: failed to get value, error:", "err", err.Error())
		return modeldb.Goods{}, fmt.Errorf("redis: failed to get value, error: %v", err)
	}
	// Переводим из json в структуру modeldb.Goods
	var good modeldb.Goods
	if err = json.Unmarshal([]byte(val), &good); err != nil {
		slog.Error(err.Error())
		return modeldb.Goods{}, err
	}

	slog.Info("Redis Get", "id", id)
	return good, nil
}
