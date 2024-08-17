package notifier

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

const (
	// deliver without sound notifications
	Info = "info"
	// deliver with sound notifications
	Warning = "warning"
	Alert   = "alert"
)

type Notifier struct {
	pubsub *pubsub.Client
	topic  string
}

type Config struct {
	Topic          string
	ProjectId      string
	ServiceKeyFile string
}

type Message struct {
	Header  string `json:"header,omitempty"`
	Content string `json:"content"`
	Channel string `json:"channel,omitempty"`
	Urgency string `json:"urgency,omitempty"`
}

func New(ctx context.Context, cfg Config) (*Notifier, error) {
	opts := option.WithCredentialsFile(cfg.ServiceKeyFile)

	ps, err := pubsub.NewClient(ctx, cfg.ProjectId, opts)
	if err != nil {
		return nil, err
	}

	n := Notifier{
		pubsub: ps,
		topic:  cfg.Topic,
	}

	return &n, nil
}

func (n *Notifier) Send(ctx context.Context, message Message) {
	data, _ := json.Marshal(message)
	msg := pubsub.Message{Data: data}

	t := n.pubsub.Topic(n.topic)
	t.Publish(ctx, &msg)
}

func (n *Notifier) Release() {
	n.pubsub.Close()
}
