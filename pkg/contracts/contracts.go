package contracts

const (
	StreamUpdates = "tg_updates"
	StreamOutbox  = "tg_outbox"
	GroupBot      = "bot"

	TaskIncoming  = "TG_INCOMING"
	TaskSend      = "TG_SEND"
	ChannelEvents = "tg_events"
)

type TGIncoming struct {
	ChatID     int64  `json:"chat_id"`
	UserID     int64  `json:"user_id"`
	Text       string `json:"text"`
	Timestamp  int64  `json:"timestamp"`   // optional: زمان دریافت از تلگرام (ثانیه)
	StartedAt  int64  `json:"started_at"`  // REQUIRED: زمان شروع پروسه (UnixMilli)
	TraceID    string `json:"trace_id"`    // برای ردیابی
}

type TGSend struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	StartedAt int64  `json:"started_at"` // از TGIncoming کپی می‌شود
	TraceID   string `json:"trace_id"`   // از TGIncoming کپی می‌شود
}
