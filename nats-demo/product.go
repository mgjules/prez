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

type product struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Stock uint16    `json:"stock"`
}

type ProductModule struct {
	nc *nats.Conn

	mu       sync.RWMutex
	products map[uuid.UUID]product
}

func NewProductModule(nc *nats.Conn) *ProductModule {
	return &ProductModule{
		nc:       nc,
		products: make(map[uuid.UUID]product),
	}
}

func (p *ProductModule) Start() error {
	// Handle product creation requests
	p.nc.Subscribe("product.create", p.handleCreateProduct())

	// Handle product queries
	p.nc.Subscribe("product.get", p.handleGetProduct())

	// Handle payment success
	p.nc.Subscribe("events.payment.success", p.handlePaymentSuccess())

	return nil
}

func (p *ProductModule) handleCreateProduct() nats.MsgHandler {
	type createProductRequest struct {
		Name  string `json:"name"`
		Stock uint16 `json:"stock"`
	}

	return func(m *nats.Msg) {
		var req createProductRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal create product request", "err", err)
			return
		}

		// Validate request
		if req.Name == "" {
			slog.Error("product name is required")
			return
		}
		if req.Stock == 0 {
			slog.Error("product stock is required")
			return
		}

		// Create product logic here
		prod := product{
			ID:    uuid.New(),
			Name:  req.Name,
			Stock: req.Stock,
		}
		p.mu.Lock()
		p.products[prod.ID] = prod
		p.mu.Unlock()

		// Publish product created event
		prodData, err := json.Marshal(prod)
		if err != nil {
			slog.Error("failed to marshal product", "err", err, "product", prod)
			return
		}
		if err := p.nc.Publish("events.product.created", prodData); err != nil {
			slog.Error("failed to publish product created event", "err", err, "product", prod)
			return
		}

		// Reply with created product
		if err := m.Respond(prodData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to product creation request", "err", err, "product", prod)
			return
		}
	}
}

func (p *ProductModule) handleGetProduct() nats.MsgHandler {
	type getProductRequest struct {
		ID uuid.UUID `json:"id"`
	}

	return func(m *nats.Msg) {
		var req getProductRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal get product request", "err", err)
			return
		}

		// Validate request
		if req.ID == uuid.Nil {
			slog.Error("product id is required")
			return
		}

		// Get product logic here
		p.mu.RLock()
		defer p.mu.RUnlock()
		prod, ok := p.products[req.ID]
		if !ok {
			slog.Error("product not found", "id", req.ID)
			return
		}

		// Reply with product
		prodData, err := json.Marshal(prod)
		if err != nil {
			slog.Error("failed to marshal product", "err", err, "prod", prod)
			return
		}
		if err := m.Respond(prodData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to product get request", "err", err, "prod", prod)
			return
		}
	}
}

func (p *ProductModule) handlePaymentSuccess() nats.MsgHandler {
	type paymentSuccessRequest struct {
		OrderID uuid.UUID `json:"order_id"`
	}

	type getOrderRequest struct {
		ID uuid.UUID `json:"id"`
	}

	type getOrderResponse struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  uint16    `json:"quantity"`
	}

	return func(m *nats.Msg) {
		var req paymentSuccessRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal payment success request", "err", err)
			return
		}

		// Validate request
		if req.OrderID == uuid.Nil {
			slog.Error("order id is required")
			return
		}

		// Get Order
		reqData, err := json.Marshal(getOrderRequest{ID: req.OrderID})
		if err != nil {
			slog.Error("failed to marshal get order request", "err", err, "order_id", req.OrderID)
			return
		}
		resp, err := p.nc.Request("order.get", reqData, time.Second)
		if err != nil {
			slog.Error("failed to get order", "err", err, "order_id", req.OrderID)
			return
		}
		var order getOrderResponse
		if err := json.Unmarshal(resp.Data, &order); err != nil {
			slog.Error("failed to unmarshal get order response", "err", err, "order_id", req.OrderID)
			return
		}

		// Update product stock
		p.mu.Lock()
		defer p.mu.Unlock()
		prod, ok := p.products[order.ProductID]
		if !ok {
			slog.Error("product not found", "id", order.ProductID)
			return
		}
		if prod.Stock < order.Quantity {
			// We should rarely be here
			slog.Error("not enough stock", "stock", prod.Stock, "count", order.Quantity)
			return
		}
		prod.Stock -= order.Quantity
		p.products[prod.ID] = prod

		// Inform updated stock
		prodData, err := json.Marshal(prod)
		if err != nil {
			slog.Error("failed to marshal product", "err", err, "prod", prod)
			return
		}
		if err := p.nc.Publish("events.product.stock.updated", prodData); err != nil {
			slog.Error("failed to publish product stock updated event", "err", err, "product", prod)
			return
		}
	}
}
