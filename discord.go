// Package discord is a Discord channel for togo notifications (webhook).
// Set DISCORD_WEBHOOK_URL; PushNotification title/body become a message + embed.
package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/togo-framework/notifications"
	"github.com/togo-framework/togo"
)

func init() {
	notifications.RegisterChannel("discord", func(k *togo.Kernel) notifications.Channel {
		return &channel{webhook: os.Getenv("DISCORD_WEBHOOK_URL"), client: &http.Client{Timeout: 15 * time.Second}}
	})
}

type channel struct {
	webhook string
	client  *http.Client
}

func (c *channel) Send(ctx context.Context, to notifications.Notifiable, n notifications.Notification) error {
	pn, ok := n.(notifications.PushNotification)
	if !ok || c.webhook == "" {
		return nil
	}
	m := pn.ToPush(to)
	body, _ := json.Marshal(map[string]any{
		"content": m.Title,
		"embeds":  []map[string]any{{"title": m.Title, "description": m.Body, "color": 1951964}},
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.webhook, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord: %s", resp.Status)
	}
	return nil
}
