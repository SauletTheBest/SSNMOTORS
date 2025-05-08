package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"order-service/internal/model"

	"github.com/nats-io/nats.go"
)

type Publisher interface {
	PublishOrderCreated(ctx context.Context, order *model.Order) error
}

type NATSPublisher struct {
	js nats.JetStreamContext
}

func NewNATSPublisher(nc *nats.Conn) (*NATSPublisher, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	return &NATSPublisher{js: js}, nil
}

func (p *NATSPublisher)PublishOrderCreated(ctx context.Context, order *model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	log.Printf("Publishing order.created event: order_id=%s, products=%+v at %v", order.ID, order.Products, time.Now())
	_, err = p.js.Publish("order.created", data)
	return err
}