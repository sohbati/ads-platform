package listener

import (
	"context"
	"encoding/json"
	"log"

	"ads-platform-stats/internal/business/stats/model"
	"ads-platform-stats/internal/business/stats/service"
	"ads-platform-stats/internal/core/broker"
)

type StatsListener struct {
	subscription broker.Subscription
}

func NewStatsListener(brokerConn *broker.Connection, subject, queue string, svc service.StatsService) (*StatsListener, error) {
	sub, err := brokerConn.QueueSubscribe(subject, queue, func(data []byte) {
		var event model.Event
		if err := json.Unmarshal(data, &event); err != nil {
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
