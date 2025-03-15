package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			DB:       0,
			Password: "games7890",
			Username: "default",
		},
	)

	for {
		time.Sleep(time.Millisecond * 200)

		pub := client.Publish(context.TODO(), "messages:", "test message")

		if pub.Err() != nil {
			fmt.Println(pub.Err().Error())
			continue
		}

		fmt.Println("sent message")
	}
}
