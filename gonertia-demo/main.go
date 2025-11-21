package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	i := initInertia()

	r := chi.NewRouter()
	r.Use(i.Middleware)
	r.Get("/", index(i))
	r.Handle("/build/*", http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	log.Println("listening on http://localhost:3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Error("http.ListenAndServe failed", "err", err)
		os.Exit(1)
	}
}

func index(i *inertia.Inertia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Index", inertia.Props{
			"text": "Inertia.js with Svelte and Go! 💙",
		})
		if err != nil {
			handleServerErr(w, err)
			return
		}
	}
}

func handleServerErr(w http.ResponseWriter, err error) {
	log.Printf("http error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("server error"))
}
