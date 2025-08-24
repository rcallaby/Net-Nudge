package notifier


import (
"bytes"
"context"
"encoding/json"
"net/http"


"scansched/internal/config"
)


type slack struct { cfg config.SlackConfig }


func NewSlack(c config.SlackConfig) Notifier { return &slack{cfg: c} }


func (s *slack) Notify(ctx context.Context, ev Event) error {
	payload := map[string]any{"text": "*"+ev.Title+"*\n"+ev.Text}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.Webhook, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)
	return err
}