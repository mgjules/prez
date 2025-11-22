package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mgjules/prez/gonertia-demo/internal/user"
	inertia "github.com/romsar/gonertia/v2"
)

func main() {
	// Bootstrap Inertia
	i, err := inertia.NewFromFile("resources/views/app.html",
		inertia.WithSSR(),
	)
	if err != nil {
		panic(err)
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

	r.Route("/users", func(r chi.Router) {
		userRepo := user.NewRepository()
		userRepo.Seed(256)
		userHandler := user.NewHandler(app, userRepo)

		r.Get("/", userHandler.Index)
		r.Post("/", userHandler.Create)
		r.Patch("/", userHandler.Update)
		r.Delete("/", userHandler.Delete)
	})

	r.Handle("/build/*", http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	log.Println("listening on http://localhost:3000")

	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}

func index(app *inertia.ViteInstance) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.Render(w, r, "Index", inertia.Props{
			"text": "Inertia.js with Svelte and Go 💙 !",
		})
		if err != nil {
			app.ShareProp("error", err.Error())
		}
	}
}
