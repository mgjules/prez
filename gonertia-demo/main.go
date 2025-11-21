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
	// Bootstrap Inertia
	i, err := inertia.NewFromFile("resources/views/app.html",
		inertia.WithSSR(),
	)
	if err != nil {
		slog.Error("inertia.NewFromFile failed", "err", err)
		os.Exit(1)
	}

	// Wrap with Vite and configure Vite-specific options
	app, err := inertia.NewWithVite(i,
		inertia.WithHotFile("public/hot"),                          // Hot reload file path
		inertia.WithBuildManifest("public/build/manifest.json"),    // Build manifest path
		inertia.WithFallbackManifest("public/.vite/manifest.json"), // Fallback manifest
		inertia.WithBuildDir("public/build"),                       // Build output directory
		inertia.WithHotReloadPort("//localhost:3000"),              // Hot reload server port
	)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(i.Middleware)
	r.Get("/", index(app))
	r.Handle("/build/*", http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	log.Println("listening on http://localhost:3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		slog.Error("failed to listen", "err", err)
		os.Exit(1)
	}
}

func index(i *inertia.ViteInstance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Index", inertia.Props{
			"text": "Inertia.js with Svelte and Go! 💙",
		})
		if err != nil {
			slog.Error("render error", "err", err)
			http.Error(w, "render error", http.StatusInternalServerError)
			return
		}
	}
}
