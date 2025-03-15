package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cwby333/pubSubRedisStudy/chat/internal/client"
)

const (
	defaultShutdownContext = time.Second * 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-ch
		slog.Info("signal", slog.String("signal", sig.String()))
		cancel()
	}()

	cfg := client.MustLoadConfig()

	client, err := client.New(ctx, cfg)

	if err != nil {
		log.Println(err)
		return
	}

	client.Connect(ctx, "testChannel")


	go func() {
		client.StartPublish(ctx)
	}()

	graceful := NewGraceful()
	
	graceful.Add(client.Close)

	select {
	case <- ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownContext)
		defer cancel()

		errors := graceful.StartGraceful(ctx)
		
		for i := range errors {
			log.Println(errors[i])
		}
	}
}