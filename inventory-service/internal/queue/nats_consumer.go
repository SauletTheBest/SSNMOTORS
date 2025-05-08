package queue

import (
    "context"
    "encoding/json"
    "inventory-service/internal/usecase"
    "log"
    "github.com/nats-io/nats.go"
)

type Consumer struct {
    conn     *nats.Conn
    subject  string
    usecase  *usecase.ProductUsecase
}

func NewConsumer(conn *nats.Conn, subject string, uc *usecase.ProductUsecase) *Consumer {
    return &Consumer{conn: conn, subject: subject, usecase: uc}
}

type OrderCreatedMessage struct {
	ID       string `json:"ID"`
	UserID   string `json:"UserID"`
	Total    float64 `json:"Total"`
	Status   string `json:"Status"`
	Products []struct {
		ProductID string `json:"ProductID"`
		Quantity  int    `json:"Quantity"`
	} `json:"Products"`
}

func (c *Consumer) Subscribe(ctx context.Context) error {
    _, err := c.conn.Subscribe(c.subject, func(msg *nats.Msg) {
        var order OrderCreatedMessage
        if err := json.Unmarshal(msg.Data, &order); err != nil {
            log.Printf("❌ Failed to parse order: %v\n", err)
            return
        }

        log.Printf("📨 Received message on %s: %+v\n", c.subject, order)

        for _, item := range order.Products {
            prod, err := c.usecase.GetProduct(ctx, item.ProductID)
            if err != nil {
                log.Printf("❌ Product not found: %s\n", item.ProductID)
                continue
            }

            if prod.Stock < int32(item.Quantity) {
                log.Printf("❌ Not enough stock for product %s: have %d, need %d\n", item.ProductID, prod.Stock, item.Quantity)
                continue
            }

            prod.Stock -= int32(item.Quantity)
            if err := c.usecase.UpdateProduct(ctx, prod); err != nil {
                log.Printf("❌ Failed to update stock for product %s: %v\n", item.ProductID, err)
                continue
            }

            log.Printf("✅ Stock updated for product %s: new stock = %d\n", item.ProductID, prod.Stock)
        }
    })

    return err
}