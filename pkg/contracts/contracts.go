package contracts

const (
	StreamUpdates = "tg_updates"
	StreamOutbox  = "tg_outbox"
	GroupBot      = "bot"

	TaskIncoming = "TG_INCOMING"
	TaskSend     = "TG_SEND"

	ChannelEvents = "tg_events" // اختیاری (Pub/Sub)
)

type TGIncoming struct {
	ChatID    int64  `json:"chat_id"`
	UserID    int64  `json:"user_id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type TGSend struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
