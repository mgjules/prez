package user

import (
	"encoding/json"
	"net/http"

	inertia "github.com/romsar/gonertia/v2"
)

type handler struct {
	app  *inertia.ViteInstance
	repo *repository
}

func NewHandler(
	app *inertia.ViteInstance,
	repo *repository,
) *handler {
	if app == nil {
		panic("app cannot be nil")
	}
	if repo == nil {
		panic("repo cannot be nil")
	}

	return &handler{
		app:  app,
		repo: repo,
	}
}

func (h *handler) Index(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List()
	if err != nil {
		h.app.ShareProp("error", err.Error())
		return
	}
	if err := h.app.Render(w, r, "User/Index", inertia.Props{
		"users": inertia.Defer(users),
	}); err != nil {
		h.app.ShareProp("error", err.Error())
	}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.app.ShareProp("error", err.Error())
		return
	}

	created, err := h.repo.Create(user)
	if err != nil {
		h.app.ShareProp("error", err.Error())
		return
	}

	if err := h.app.Render(w, r, "User/Index", inertia.Props{
		"user_created": created,
	}); err != nil {
		h.app.ShareProp("error", err.Error())
	}
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.app.ShareProp("error", err.Error())
		return
	}

	updated, err := h.repo.Update(user)
	if err != nil {
		h.app.ShareProp("error", err.Error())
		return
	}

	if err := h.app.Render(w, r, "User/Index", inertia.Props{
		"user_updated": updated,
	}); err != nil {
		h.app.ShareProp("error", err.Error())
	}
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.app.ShareProp("error", err.Error())
		return
	}

	if err := h.repo.Delete(user.ID); err != nil {
		h.app.ShareProp("error", err.Error())
	}

	if err := h.app.Render(w, r, "User/Index", inertia.Props{
		"user_id_deleted": user.ID,
	}); err != nil {
		h.app.ShareProp("error", err.Error())
	}
}
