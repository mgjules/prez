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

type order struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  uint16    `json:"quantity"`
}

type OrderModule struct {
	nc *nats.Conn

	mu     sync.RWMutex
	orders map[uuid.UUID]order
}

func NewOrderModule(nc *nats.Conn) *OrderModule {
	return &OrderModule{
		nc:     nc,
		orders: make(map[uuid.UUID]order),
	}
}

func (o *OrderModule) Start() error {
	// Handle order creation requests
	o.nc.Subscribe("order.create", o.handleCreateOrder())

	// Handle order queries
	o.nc.Subscribe("order.get", o.handleGetOrder())

	return nil
}

func (o *OrderModule) handleCreateOrder() nats.MsgHandler {
	type createOrderRequest struct {
		UserID    uuid.UUID `json:"user_id"`
		ProductID uuid.UUID `json:"product_id"`
		Quantity  uint16    `json:"quantity"`
	}

	type getUserRequest struct {
		ID uuid.UUID `json:"id"`
	}
	type getUserResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	type getProductRequest struct {
		ID uuid.UUID `json:"id"`
	}
	type getProductResponse struct {
		ID    uuid.UUID `json:"id"`
		Stock uint16    `json:"stock"`
	}

	return func(m *nats.Msg) {
		var req createOrderRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal create order request", "err", err)
			return
		}

		// Validate request
		if req.UserID == uuid.Nil {
			slog.Error("user id is required")
			return
		}
		if req.ProductID == uuid.Nil {
			slog.Error("product id is required")
			return
		}
		if req.Quantity == 0 {
			slog.Error("quantity is required")
			return
		}

		// Get User (make sure they exist)
		reqData, err := json.Marshal(getUserRequest{ID: req.UserID})
		if err != nil {
			slog.Error("failed to marshal get user request", "err", err, "user_id", req.UserID)
			return
		}
		resp, err := o.nc.Request("user.get", reqData, time.Second)
		if err != nil {
			slog.Error("failed to get user", "err", err, "user_id", req.UserID)
			return
		}
		var user getUserResponse
		if err := json.Unmarshal(resp.Data, &user); err != nil {
			slog.Error("failed to unmarshal get user response", "err", err, "user_id", req.UserID)
			return
		}

		// Get Product (make sure there's enough stock)
		reqData, err = json.Marshal(getProductRequest{ID: req.ProductID})
		if err != nil {
			slog.Error("failed to marshal get product request", "err", err, "product_id", req.ProductID)
			return
		}
		resp, err = o.nc.Request("product.get", reqData, time.Second)
		if err != nil {
			slog.Error("failed to get product", "err", err, "product_id", req.ProductID)
			return
		}
		var product getProductResponse
		if err := json.Unmarshal(resp.Data, &product); err != nil {
			slog.Error("failed to unmarshal get product response", "err", err, "product_id", req.ProductID)
			return
		}
		if product.Stock < req.Quantity {
			slog.Error("not enough stock", "product", product, "quantity", req.Quantity)
			return
		}

		// Create order logic here
		o.mu.Lock()
		defer o.mu.Unlock()
		ordr := order{
			ID:        uuid.New(),
			UserID:    user.ID,
			ProductID: product.ID,
			Quantity:  req.Quantity,
		}
		o.orders[ordr.ID] = ordr

		slog.Info("placed order", "order", ordr)

		// Publish order created event
		orderData, err := json.Marshal(ordr)
		if err != nil {
			slog.Error("failed to marshal order", "err", err, "order", ordr)
			return
		}
		if err := o.nc.Publish("events.order.created", orderData); err != nil {
			slog.Error("failed to publish order created event", "err", err, "order", ordr)
			return
		}

		// Reply with created order
		if err := m.Respond(orderData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to order creation request", "err", err, "order", ordr)
			return
		}
	}
}

func (o *OrderModule) handleGetOrder() nats.MsgHandler {
	type getOrderRequest struct {
		ID uuid.UUID `json:"id"`
	}

	return func(m *nats.Msg) {
		var req getOrderRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal get order request", "err", err)
			return
		}

		// Validate request
		if req.ID == uuid.Nil {
			slog.Error("order id is required")
			return
		}

		// Get order logic here
		o.mu.RLock()
		defer o.mu.RUnlock()
		ordr, ok := o.orders[req.ID]
		if !ok {
			slog.Error("order not found", "id", req.ID)
			return
		}

		// Reply with order
		orderData, err := json.Marshal(ordr)
		if err != nil {
			slog.Error("failed to marshal order", "err", err, "order", ordr)
			return
		}
		if err := m.Respond(orderData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to order get request", "err", err, "order", ordr)
			return
		}
	}
}
