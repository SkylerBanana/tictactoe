package redis

import (
	"context"

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
