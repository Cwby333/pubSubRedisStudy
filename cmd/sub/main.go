package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "games7890",
		Username: "default",
	})

	pubSub := client.Subscribe(context.Background(), "messages:")

	m, err := pubSub.Receive(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	switch m.(type) {
	case *redis.Subscription:
		sub := m.(*redis.Subscription)
		fmt.Printf("success subscribe:%s %d", sub.Channel, sub.Count)
	}

	ch := pubSub.Channel()
	ch2 := pubSub.Channel()

	go func() {
		for message := range ch {
			fmt.Println("from one", message.Payload)
		}
	}()

	go func() {
		for message := range ch2 {
			fmt.Println("from two", message.Payload)
		}
	}()

	time.Sleep(time.Second * 30)
	pubSub.Close()
}
