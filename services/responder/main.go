package main

import (
	"context"
	"log"
	"os"

	broker "github.com/mrjvadi/go-broker/broker"
	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: getenv("REDIS_ADDR", "127.0.0.1:6379"),
	})

	streamOut := getenv("STREAM_OUTBOX", "tg_outbox")
	group := getenv("GROUP_NAME", "bot")

	app := broker.New(rdb, streamOut, group, 16)

	botToken := mustGetenv("BOT_TOKEN")
	// نیازی به Start کردن Poller نداریم؛ فقط برای Send استفاده می‌کنیم.
	bot, err := tele.NewBot(tele.Settings{Token: botToken})
	if err != nil {
		log.Fatal(err)
	}

	app.OnTask("TG_SEND", func(ctx context.Context, msg map[string]any) error {
		chatID := int64Val(msg["chat_id"])
		text := stringVal(msg["text"])
		if chatID == 0 || text == "" {
			return nil
		}
		_, err := bot.Send(tele.ChatID(chatID), text)
		return err
	})

	log.Println("responder started (consuming & sending)…")
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env %s", k)
	}
	return v
}
func stringVal(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
func int64Val(v any) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	case float64:
		return int64(t)
	default:
		return 0
	}
}
