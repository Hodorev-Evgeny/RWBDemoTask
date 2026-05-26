package feature_ingest

import (
	core_domain "RWBDwmoTask/internal/core/domain"
	core_storage "RWBDwmoTask/internal/core/storage"
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type NatsConsumer struct {
	js      jetstream.JetStream
	storage *core_storage.Storage
}

func NewNatsConsumer(js jetstream.JetStream, storage *core_storage.Storage) *NatsConsumer {
	return &NatsConsumer{
		js:      js,
		storage: storage,
	}
}

func (n *NatsConsumer) ReadEvents(
	ctx context.Context,
) error {
	stream, err := n.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "SEARCH",
		Subjects: []string{"search.events"},
		MaxAge:   10 * time.Minute,
	})
	if err != nil {
		return err
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "trend-service",
		FilterSubject: "search.events",
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	consumerCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		var event core_domain.SearchEvent

		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			_ = msg.Ack()
			return
		}

		n.storage.Add(
			event.UserID,
			event.SessionID,
			event.Query,
			1,
		)

		_ = msg.Ack()
	})
	if err != nil {
		return err
	}

	<-ctx.Done()
	consumerCtx.Stop()

	return nil
}
