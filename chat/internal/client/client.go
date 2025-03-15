package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	Clientname  string
	redisClient *redis.Client
}

func New(ctx context.Context, cfg Config) (Client, error) {
	rClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Addr,
		Username: cfg.RedisConfig.Username,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})

	cmd := rClient.Ping(ctx)
	if cmd.Err() != nil {
		return Client{}, cmd.Err()
	}

	err := checkExistClientname(ctx, cfg.ClientConfig.Username, rClient)
	if err != nil {
		return Client{}, err
	}

	client := Client{
		Clientname:  cfg.ClientConfig.Username,
		redisClient: rClient,
	}

	cmd2 := rClient.HSet(ctx, "usernames", cfg.ClientConfig.Username, "")
	if cmd2.Err() != nil {
		return Client{}, cmd2.Err()
	}

	return client, nil
}

func (c Client) Connect(ctx context.Context, chanNames ...string) {
	for i := range chanNames {
		sub := c.redisClient.Subscribe(ctx, chanNames[i])

		fmt.Printf("start listen: %s\n", sub.String())

		go func() {
			ch := sub.Channel()

			for message := range ch {
				fmt.Printf("Receive message: %s", message.Payload)
			}
		}()
	}
}

func (c Client) StartPublish(ctx context.Context) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		fmt.Println("Print message:")
		scanner.Scan()
		message := c.Clientname + ": " + scanner.Text() + "\n"

		fmt.Println("Print channel:")
		scanner.Scan()
		channel := scanner.Text()

		pub := c.redisClient.Publish(ctx, channel, message)
		if pub.Err() != nil {
			log.Println(pub.Err())
			continue
		}
	}
}

func (c Client) Close(ctx context.Context) error {
	select {
	case <- ctx.Done():
		return ctx.Err()
	default:
	}

	cmd := c.redisClient.HDel(ctx, "usernames", c.Clientname)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func checkExistClientname(ctx context.Context, name string, redisClient *redis.Client) error {
	cmd := redisClient.HGet(ctx, "clientNames:", name)

	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return nil
		}

		return cmd.Err()
	}

	return errors.New("clientname exists")
}
