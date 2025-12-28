package bot

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

func (b *Bot) loopWebUI() {
	mux := http.NewServeMux()

	// Main Page
	mux.HandleFunc("/", b.handleDashboard)

	// 2. The Detail Views (Master-Detail Pattern)
	// This matches <button hx-get="/views/fishing" ...> in index.html
	mux.HandleFunc("/views/fishing", b.handleFishingView)
	mux.HandleFunc("/views/stats", b.handleStatsView)

	// 3. API Endpoints (Fragments)
	// This matches <button hx-post="/api/toggle-fishing" ...>
	mux.HandleFunc("/api/toggle-fishing", b.handleToggleFishing)

	// API Endpoints
	//mux.HandleFunc("/api/stats", b.handleGetStats) // Simple text return

	fmt.Println("[UI] Dashboard live at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

//go:embed web/*
var webAssets embed.FS
var templates = template.Must(
	template.New(""). // Create a base template container
				ParseFS(webAssets, // Now it is safe to parse
			"web/index.gohtml",
			"web/**/*.gohtml",
		),
)

const DevMode = true

func (b *Bot) renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	var t *template.Template
	var err error

	if DevMode {
		// CRITICAL: Re-parse files on EVERY request so changes show instantly
		// We use os.DirFS to ensure we are looking at the actual disk
		t, err = template.New("").
			ParseFS(os.DirFS("internal/bot"),
				"web/index.gohtml",
				"web/components/*.gohtml",
				"web/views/*.gohtml")
	} else {
		// In production, use the pre-parsed global 'templates' variable
		// (or parse once at startup from embeddedAssets)
		t = templates
	}

	if err != nil {
		http.Error(w, "Template Error: "+err.Error(), 500)
		return
	}

	t.ExecuteTemplate(w, templateName, data)
}
