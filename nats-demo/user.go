package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type user struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UserModule struct {
	nc *nats.Conn

	mu    sync.RWMutex
	users map[uuid.UUID]user
}

func NewUserModule(nc *nats.Conn) *UserModule {
	return &UserModule{
		nc:    nc,
		users: make(map[uuid.UUID]user),
	}
}

func (u *UserModule) Start() error {
	// Handle user creation requests
	u.nc.Subscribe("user.create", u.handleCreateUser())

	// Handle user queries
	u.nc.Subscribe("user.get", u.handleGetUser())

	return nil
}

func (u *UserModule) handleCreateUser() nats.MsgHandler {
	type createUserRequest struct {
		Name string `json:"name"`
	}

	return func(m *nats.Msg) {
		var req createUserRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal create user request", "err", err)
			return
		}

		// Validate request
		if req.Name == "" {
			slog.Error("user name is required")
			return
		}

		// Create user logic here
		usr := user{
			ID:   uuid.New(),
			Name: req.Name,
		}
		u.mu.Lock()
		u.users[usr.ID] = usr
		u.mu.Unlock()

		// Publish user created event
		userData, err := json.Marshal(usr)
		if err != nil {
			slog.Error("failed to marshal user", "err", err, "user", usr)
			return
		}
		if err := u.nc.Publish("events.user.created", userData); err != nil {
			slog.Error("failed to publish user created event", "err", err, "user", usr)
			return
		}

		// Reply with created user
		if err := m.Respond(userData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to user creation request", "err", err, "user", usr)
			return
		}
	}
}

func (u *UserModule) handleGetUser() nats.MsgHandler {
	type getUserRequest struct {
		ID uuid.UUID `json:"id"`
	}

	return func(m *nats.Msg) {
		var req getUserRequest
		if err := json.Unmarshal(m.Data, &req); err != nil {
			slog.Error("failed to unmarshal get user request", "err", err)
			return
		}

		// Validate request
		if req.ID == uuid.Nil {
			slog.Error("user id is required")
			return
		}

		// Get user logic here
		u.mu.RLock()
		defer u.mu.RUnlock()
		usr, ok := u.users[req.ID]
		if !ok {
			slog.Error("user not found", "id", req.ID)
			return
		}

		// Reply with user
		userData, err := json.Marshal(usr)
		if err != nil {
			slog.Error("failed to marshal user", "err", err, "user", usr)
			return
		}
		if err := m.Respond(userData); err != nil && !errors.Is(err, nats.ErrMsgNoReply) {
			slog.Error("failed to respond to user get request", "err", err, "user", usr)
			return
		}
	}
}
