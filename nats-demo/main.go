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
	if err = nc.Publish(
		"order.create",
		[]byte(`{"user_id": "`+user.ID.String()+`", "product_id": "`+product.ID.String()+`", "quantity": 1}`),
	); err != nil {
		slog.Error("failed to place order", "err", err)
		return
	}

	select {}
}
