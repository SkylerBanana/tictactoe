package redis

import (
	"context"
	"time"

	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type store struct {
	client *redis.Client
}

type Store interface {
	Subscribe(ctx context.Context, channel string, handler func(message string)) error
	Publish(ctx context.Context, channel string, message []byte) error
	Que(ctx context.Context, player []byte) error
	QuePop(ctx context.Context, timeout time.Duration) ([]string, error)
	Length(ctx context.Context, key string) (int64, error)
}

func NewRedisInstance(redis *redis.Client) Store {
	return &store{redis}
}

func (s *store) Subscribe(ctx context.Context, channel string, handler func(message string)) error {
	pubsub := s.client.Subscribe(ctx, channel)

	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		//fmt.Println(msg.Channel, msg.Payload)
		handler(msg.Payload)
		if err := ctx.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (s *store) Publish(ctx context.Context, channel string, message []byte) error {

	err := s.client.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) Que(ctx context.Context, player []byte) error {
	err := s.client.LPush(ctx, "MatchMaking", player).Err()
	if err != nil {
		return err
	}
	return nil

}

func (s *store) Length(ctx context.Context, key string) (int64, error) {
	return s.client.LLen(ctx, key).Result()
}

func (s *store) QuePop(ctx context.Context, timeout time.Duration) ([]string, error) {
	result, err := s.client.BRPop(ctx, timeout, "MatchMaking").Result()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	address := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0, // Use default DB
		Protocol: 2, // Connection protocol
	})
	return client
}
