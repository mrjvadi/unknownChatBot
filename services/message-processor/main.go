package main

import (
	"context"
	"log"
	"os"
	"strings"

	broker "github.com/mrjvadi/go-broker/broker"
	"github.com/mrjvadi/unknownChatBot/pkg/contracts"
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
		if err := ctx.Bind(&msg); err != nil {
			log.Printf("failed to decode message: %v", err)
			return nil // ignore decoding errors
		}

		reply := "متوجه نشدم. /help رو بزن."
		switch {
		case strings.HasPrefix(msg.Text, "/start"):
			reply = "سلام! 👋 به ربات خوش اومدی."
		case strings.HasPrefix(msg.Text, "/help"):
			reply = "دستورات: /start /help"
		}

		ctxs := context.Background()

		out := broker.New(rdb, streamOut, group, broker.WithMaxJobs(32))
		out.Enqueue(ctxs, "TG_SEND", contracts.TGSend{
			ChatID: msg.ChatID,
			Text:   reply,
		})
		return nil
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
