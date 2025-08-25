package notifier


import (
	"context"
)


type Event struct {
	Title string `json:"title"`
	Text string `json:"text"`
	Meta map[string]any `json:"meta"`
}


type Notifier interface { Notify(ctx context.Context, e Event) error }


