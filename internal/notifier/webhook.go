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
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.cfg.URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if w.cfg.Secret != "" {
		// Use HMAC-SHA256 to sign the payload
		import "crypto/hmac"
		import "crypto/sha256"
		import "encoding/hex"
		mac := hmac.New(sha256.New, []byte(w.cfg.Secret))
		mac.Write(b)
		signature := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Scansched-Signature", signature)
	}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}