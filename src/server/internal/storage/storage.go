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
var placesKey = "places"
var argsKeys = []string{"mode", "color", "brightness"}


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

// ------------------------------------------------------------ //

func (s *Storage) AddPlace(place string) error {
	_, err := s.client.RPush(ctx, placesKey, place).Result()
	if err != nil {	return err }
	err = s.UpdatePlace(place, []string{"", "ffffff", "1"})
	return err
}

func (s *Storage) GetPlaces() []string {
	val, err := s.client.LRange(ctx, placesKey, 0, -1).Result()
	if err != nil {
		log.Println("[storage]", err)
		return []string{}
	}
	return val
}

func (s *Storage) RemovePlace(place string) error {
	_, err := s.client.LRem(ctx, placesKey, 0, place).Result()
	return err
}

// ----------------------------------------------------------- //

func (s *Storage) UpdatePlace(place string, args []string) error {
	keyValue := make([]string, len(args)*2)
	for i, v := range args {
		keyValue[i*2] = argsKeys[i]
		keyValue[i*2+1] = v
	}
	_, err := s.client.HSet(ctx, placesKey + ":" + place, keyValue).Result()
	return err
}

func (s *Storage) GetPlace(place string) []string {
	res, err := s.client.HGetAll(ctx, placesKey + ":" + place).Result()
	if err != nil {
		log.Println("[storage]", err)
		return []string{}
	}
	vals := make([]string, len(argsKeys))
	for i := 0; i < len(argsKeys); i++ {
		vals[i] = res[argsKeys[i]]
	}
	return vals
}
