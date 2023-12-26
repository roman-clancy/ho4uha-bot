package main

import (
	"context"
	"github.com/roman-clancy/ho4uha-bot/internal/client"
	"github.com/roman-clancy/ho4uha-bot/internal/config"
	"github.com/roman-clancy/ho4uha-bot/internal/model/messages"
	"github.com/roman-clancy/ho4uha-bot/internal/storage/inmemory"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	tgClient, err := client.New(cfg.Token, client.ProcessingMessage)
	if err != nil {
		return
	}
	storage, err := inmemory.New()
	if err != nil {
		return
	}
	botModel := messages.New(storage, tgClient)
	tgClient.ListenUpdates(botModel)
	//opts := []bot.Option{
	//	bot.WithDefaultHandler(handler),
	//	bot.
	//}
	//
	//
	//b.Start(ctx)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
