package main

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {
	// Start embedded NATS server
	ns, err := server.NewServer(&server.Options{
		DontListen: true,
	})
	if err != nil {
		slog.Error("failed to start embedded nats server", "err", err)
		return
	}

	go ns.Start()
	defer ns.Shutdown()

	if !ns.ReadyForConnections(5 * time.Second) {
		slog.Error("nats server not ready for connections")
		return
	}

	slog.Info("nats server started")

	// In-process connection (no TCP)
	nc, err := nats.Connect(
		ns.ClientURL(),
		nats.InProcessServer(ns),
	)
	if err != nil {
		slog.Error("failed to connect to nats server", "err", err)
		return
	}
	defer nc.Close()

	slog.Info("connected to nats server")

	// Start modules
	userModule := NewUserModule(nc)
	userModule.Start()

	productModule := NewProductModule(nc)
	productModule.Start()

	orderModule := NewOrderModule(nc)
	orderModule.Start()

	paymentModule := NewPaymentModule(nc)
	paymentModule.Start()

	// Add a user
	if _, err := nc.Subscribe("events.user.created", func(m *nats.Msg) {
		var user struct {
			ID uuid.UUID `json:"id"`
		}
		if err := json.Unmarshal(m.Data, &user); err != nil {
			slog.Error("failed to unmarshal user created event", "err", err)
			return
		}
		slog.Info("user created", "user", user)
	}); err != nil {
		slog.Error("failed to subscribe to user created event", "err", err)
		return
	}
	resp, err := nc.Request("user.create", []byte(`{"name": "John Doe"}`), time.Second)
	if err != nil {
		slog.Error("failed to create user", "err", err)
		return
	}
	var user struct {
		ID uuid.UUID `json:"id"`
	}
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		slog.Error("failed to unmarshal user response", "err", err)
		return
	}

	// Add a product
	if _, err := nc.Subscribe("events.product.created", func(m *nats.Msg) {
		var product struct {
			ID    uuid.UUID `json:"id"`
			Name  string    `json:"name"`
			Stock uint16    `json:"stock"`
		}
		if err := json.Unmarshal(m.Data, &product); err != nil {
			slog.Error("failed to unmarshal user created event", "err", err)
			return
		}
		slog.Info("product added", "product", product)
	}); err != nil {
		slog.Error("failed to subscribe to product created event", "err", err)
		return
	}
	resp, err = nc.Request("product.create", []byte(`{"name": "Widget", "stock": 100}`), time.Second)
	if err != nil {
		slog.Error("failed to create product", "err", err)
		return
	}
	var product struct {
		ID uuid.UUID `json:"id"`
	}
	if err := json.Unmarshal(resp.Data, &product); err != nil {
		slog.Error("failed to unmarshal product response", "err", err)
		return
	}

	// Place an order
	if _, err := nc.Subscribe("events.order.created", func(m *nats.Msg) {
		var order struct {
			ID        uuid.UUID `json:"id"`
			UserID    uuid.UUID `json:"user_id"`
			ProductID uuid.UUID `json:"product_id"`
			Quantity  uint16    `json:"quantity"`
		}
		if err := json.Unmarshal(m.Data, &order); err != nil {
			slog.Error("failed to unmarshal order created event", "err", err)
			return
		}
		slog.Info("order placed", "order", order)
	}); err != nil {
		slog.Error("failed to subscribe to order created event", "err", err)
		return
	}
	if err = nc.Publish(
		"order.create",
		[]byte(`{"user_id": "`+user.ID.String()+`", "product_id": "`+product.ID.String()+`", "quantity": 10}`),
	); err != nil {
		slog.Error("failed to place order", "err", err)
		return
	}

	sub, err := nc.SubscribeSync("events.payment.success")
	if err != nil {
		slog.Error("failed to subscribe to payment success event", "err", err)
		return
	}

	msg, err := sub.NextMsg(3 * time.Second)
	if err != nil {
		slog.Error("failed to receive payment success event", "err", err)
		return
	}
	var payment struct {
		ID      uuid.UUID `json:"id"`
		OrderID uuid.UUID `json:"order_id"`
	}
	if err := json.Unmarshal(msg.Data, &payment); err != nil {
		slog.Error("failed to unmarshal payment success event", "err", err)
		return
	}

	slog.Info("payment successful", "payment", payment)
}
