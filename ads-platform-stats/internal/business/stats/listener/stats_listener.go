package listener

import (
	"context"
	"encoding/json"
	"log"

	"ads-platform-stats/internal/business/stats/model"
	"ads-platform-stats/internal/business/stats/service"
	"ads-platform-stats/internal/core/natsconn"

	"github.com/nats-io/nats.go"
)

type StatsListener struct {
	subscription *nats.Subscription
}

func NewStatsListener(natsConn *natsconn.Connection, subject, queue string, svc service.StatsService) (*StatsListener, error) {
	sub, err := natsConn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		var event model.Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("[stats] parse error: %v", err)
			return
		}
		if err := svc.Apply(context.Background(), event); err != nil {
			log.Printf("[stats] apply ad_id=%d event=%s: %v", event.AdID, event.Event, err)
		}
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[stats] listening on %s queue=%s", subject, queue)
	return &StatsListener{subscription: sub}, nil
}

func (l *StatsListener) Stop() {
	if l.subscription != nil {
		_ = l.subscription.Unsubscribe()
	}
}
