package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"

	broker "github.com/mrjvadi/go-broker/broker"
	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: getenv("REDIS_ADDR", "127.0.0.1:6379"),
	})

	streamUpdates := getenv("STREAM_UPDATES", "tg_updates")
	groupName := getenv("GROUP_NAME", "bot")
	app := broker.New(rdb, streamUpdates, groupName, 0) // فقط ارسال می‌کنیم، مصرف نداریم

	botToken := mustGetenv("BOT_TOKEN")
	bot, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 60 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	// فقط متن‌ها؛ می‌تونی هندلرهای دیگر هم اضافه کنی
	bot.Handle(tele.OnText, func(c tele.Context) error {
		msg := c.Message()
		if msg == nil || msg.Chat == nil {
			return nil
		}

		payload := map[string]any{
			"chat_id":   msg.Chat.ID,
			"user_id":   c.Sender().ID,
			"text":      c.Text(),
			"timestamp": time.Now().Unix(),
		}

		if err := app.Enqueue(ctx, "TG_INCOMING", payload); err != nil {
			log.Printf("enqueue error: %v", err)
		}
		// رویداد لحظه‌ای (اختیاری)
		_ = app.Publish(ctx, "tg_events", map[string]any{
			"type": "received", "chat_id": msg.Chat.ID,
		})
		return nil
	})

	log.Println("telegram-gateway started (polling)…")
	bot.Start()
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
