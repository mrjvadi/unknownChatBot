package main

import (
	"context"
	"log"
	"os"
	"strings"

	broker "github.com/mrjvadi/go-broker/broker"
	contracts "github.com/mrjvadi/pkg/contracts"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: getenv("REDIS_ADDR", "127.0.0.1:6379"),
	})

	streamIn := getenv("STREAM_UPDATES", "tg_updates")
	streamOut := getenv("STREAM_OUTBOX", "tg_outbox")
	group := getenv("GROUP_NAME", "bot")

	app := broker.New(rdb, streamIn, group, broker.WithStreamLength(32))

	app.OnTask("TG_INCOMING", func(ctx *broker.Context) error {
		msg := contracts.TGIncoming{}
		text := stringVal(msg["text"])
		chatID := int64Val(msg["chat_id"])

		reply := "متوجه نشدم. /help رو بزن."
		switch {
		case strings.HasPrefix(text, "/start"):
			reply = "سلام! 👋 به ربات خوش اومدی."
		case strings.HasPrefix(text, "/help"):
			reply = "دستورات: /start /help"
		}

		out := broker.New(rdb, streamOut, group, broker.WithMaxJobs(32))
		out.Enqueue(ctx, "TG_SEND", map[string]any{
			"chat_id": chatID,
			"text":    reply,
		})
	})

	log.Println("message-processor started (consuming)…")
	app.Run(ctx)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
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
