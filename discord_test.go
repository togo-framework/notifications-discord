package discord

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/togo-framework/notifications"
)

type nfy struct{}

func (nfy) RouteID() string           { return "1" }
func (nfy) RouteEmail() string        { return "" }
func (nfy) RoutePushTokens() []string { return nil }

type note struct{ t, b string }

func (note) Via(notifications.Notifiable) []string { return []string{"discord"} }
func (n note) ToPush(notifications.Notifiable) notifications.PushMessage {
	return notifications.PushMessage{Title: n.t, Body: n.b}
}

func TestSendPostsWebhook(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		w.WriteHeader(204)
	}))
	defer srv.Close()
	c := &channel{webhook: srv.URL, client: srv.Client()}
	if err := c.Send(context.Background(), nfy{}, note{t: "Deploy ok", b: "v1 live"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	if got["content"] != "Deploy ok" {
		t.Fatalf("content = %v", got["content"])
	}
}

func TestNoWebhookIsNoop(t *testing.T) {
	c := &channel{client: http.DefaultClient}
	if err := c.Send(context.Background(), nfy{}, note{t: "x", b: "y"}); err != nil {
		t.Fatalf("got %v", err)
	}
}
