package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	inertia "github.com/romsar/gonertia"
)

func main() {
	i := initInertia()

	mux := http.NewServeMux()
	mux.Handle("GET /", i.Middleware(indexHandler(i)))
	mux.Handle("GET /build/", http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	log.Println("listening on http://localhost:3000")

	if err := http.ListenAndServe(":3000", mux); err != nil {
		slog.Error("http.ListenAndServe failed", "err", err)
		os.Exit(1)
	}
}

func initInertia() *inertia.Inertia {
	const (
		viteHotFile  = "./public/hot"
		rootViewFile = "resources/views/app.html"
	)

	// check if laravel-vite-plugin is running in dev mode (it puts a "hot" file in the public folder)
	_, err := os.Stat(viteHotFile)
	if err == nil {
		i, err := inertia.NewFromFile(
			rootViewFile,
			inertia.WithSSR(),
		)
		if err != nil {
			slog.Error("inertia.NewFromFile failed", "err", err)
			os.Exit(1)
		}

		if err = i.ShareTemplateFunc("vite", func(entry string) (string, error) {
			content, err := os.ReadFile(viteHotFile)
			if err != nil {
				return "", err
			}
			url := strings.TrimSpace(string(content))
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				url = url[strings.Index(url, ":")+1:]
			} else {
				url = "//localhost:8080"
			}
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}
			return url + entry, nil
		}); err != nil {
			slog.Warn("shareTemplateFunc failed", "err", err)
		}

		return i
	}

	// laravel-vite-plugin not running in dev mode, use build manifest file
	const manifestPath = "./public/build/manifest.json"

	// check if the manifest file exists, if not, rename it
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		// move the manifest from ./public/build/.vite/manifest.json to ./public/build/manifest.json
		// so that the vite function can find it
		err := os.Rename("./public/build/.vite/manifest.json", "./public/build/manifest.json")
		if err != nil {
			return nil
		}
	}

	i, err := inertia.NewFromFile(
		rootViewFile,
		inertia.WithVersionFromFile(manifestPath),
		inertia.WithSSR(),
	)
	if err != nil {
		slog.Error("inertia.NewFromFile failed", "err", err)
		os.Exit(1)
	}

	if err := i.ShareTemplateFunc("vite", vite(manifestPath, "/build/")); err != nil {
		slog.Warn("shareTemplateFunc failed", "err", err)
		os.Exit(1)
	}

	return i
}

func vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		slog.Error("os.Open for manifest failed", "err", err)
		os.Exit(1)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)
	// print content of viteAssets
	for k, v := range viteAssets {
		log.Printf("%s: %s\n", k, v.File)
	}

	if err != nil {
		log.Fatalf("cannot unmarshal vite manifest file to json: %s", err)
	}

	return func(p string) (string, error) {
		if val, ok := viteAssets[p]; ok {
			return path.Join("/", buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q not found", p)
	}
}

func indexHandler(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := i.Render(w, r, "Index", inertia.Props{
			"text": "Inertia.js with Svelte and Go! 💙",
		})
		if err != nil {
			handleServerErr(w, err)
			return
		}
	}

	return http.HandlerFunc(fn)
}

func handleServerErr(w http.ResponseWriter, err error) {
	log.Printf("http error: %s\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("server error"))
}
