package notifier


import (
"bytes"
"context"
"encoding/json"
"net/http"


"scansched/internal/config"
)


type webhook struct { cfg config.WebhookConfig }


func NewWebhook(c config.WebhookConfig) Notifier { return &webhook{cfg: c} }


func (w *webhook) Notify(ctx context.Context, ev Event) error {
	b, _ := json.Marshal(ev)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, w.cfg.URL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if w.cfg.Secret != "" { req.Header.Set("X-Scansched-Signature", w.cfg.Secret) }
		_, err := http.DefaultClient.Do(req)
	return err
}