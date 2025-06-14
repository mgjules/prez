package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type status string

const (
	statusPending status = "pending"
	statusPaid    status = "paid"
	statusFailed  status = "failed"
)

type payment struct {
	ID      uuid.UUID `json:"id"`
	OrderID uuid.UUID `json:"order_id"`
	Status  status    `json:"status"`
}

type PaymentModule struct {
	nc *nats.Conn

	mu       sync.RWMutex
	payments map[uuid.UUID]payment
}

func NewPaymentModule(nc *nats.Conn) *PaymentModule {
	return &PaymentModule{
		nc:       nc,
		payments: make(map[uuid.UUID]payment),
	}
}

func (p *PaymentModule) Start() error {
	p.nc.Subscribe("events.order.created", p.handleOrderCreated())

	return nil
}

func (p *PaymentModule) handleOrderCreated() nats.MsgHandler {
	type orderCreatedRequest struct {
		ID uuid.UUID `json:"id"`
	}
	type getOrderRequest struct {
		ID uuid.UUID `json:"id"`
	}
	type getOrderResponse struct {
		ID uuid.UUID `json:"id"`
	}

	return func(m *nats.Msg) {
		var req orderCreatedRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal create user request", "err", err)
			return
		}

		// Validate request
		if req.ID == uuid.Nil {
			slog.Error("order id is required")
			return
		}

		// Get Order (make sure it exists)
		reqData, err := json.Marshal(getOrderRequest{ID: req.ID})
		if err != nil {
			slog.Error("failed to marshal get order request", "err", err, "order_id", req.ID)
			return
		}
		resp, err := p.nc.Request("order.get", reqData, time.Second)
		if err != nil {
			slog.Error("failed to get order", "err", err, "order_id", req.ID)
			return
		}
		var order getOrderResponse
		if err := json.Unmarshal(resp.Data, &order); err != nil {
			slog.Error("failed to unmarshal get order response", "err", err, "order_id", req.ID)
			return
		}

		paym := payment{
			ID:      uuid.New(),
			OrderID: order.ID,
			Status:  statusPending,
		}

		// Simulate payment process
		time.Sleep(time.Second)

		// We assume the user has moneeeeyyy
		paym.Status = statusPaid

		// Create payment logic here
		p.mu.Lock()
		p.payments[paym.ID] = paym
		p.mu.Unlock()

		slog.Info("created payment", "payment", paym)

		// Publish payment success event
		paymentData, err := json.Marshal(paym)
		if err != nil {
			slog.Error("failed to marshal payment", "err", err, "payment", paym)
			return
		}
		if err := p.nc.Publish("events.payment.success", paymentData); err != nil {
			slog.Error("failed to publish payment success event", "err", err, "payment", paym)
			return
		}

		// Reply with created payment
		if err := m.Respond(paymentData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to payment success request", "err", err, "payment", paym)
			return
		}
	}
}
