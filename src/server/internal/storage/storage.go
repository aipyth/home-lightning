package storage

import (
	"log"
	"os"
	"context"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	client *redis.Client
}

// var ctx, cancel = context.WithTimeout(context.Background(), 300 * time.Millisecond)
var ctx = context.Background()

var modesKey = "modes"


func connectToRedis() (*redis.Client, error) {
	redisUrl := os.Getenv("REDIS_URL")

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func NewStorage() (*Storage, error) {
	redisClient, err := connectToRedis()
	return &Storage{redisClient}, err
}

func (s *Storage) GetModes() []string {
	val, err := s.client.LRange(ctx, modesKey, 0, -1).Result()
	if err != nil {
		log.Printf("[storage] %v", err)
		return []string{}
	}
	return val
}

func (s *Storage) AddMode(mode string) error {
	_, err := s.client.RPush(ctx, modesKey, mode).Result()
	return err
}

func (s *Storage) RemoveMode(mode string) error {
	_, err := s.client.LRem(ctx, modesKey, 0, mode).Result()
	return err
}


