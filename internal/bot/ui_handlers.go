package bot

import (
	"net/http"
)

func (b *Bot) handleDashboard(w http.ResponseWriter, r *http.Request) {
	b.renderTemplate(w, "index.gohtml", nil)
}

func (b *Bot) handleFishingView(w http.ResponseWriter, r *http.Request) {
	// 1. If it's a direct browser access (F5), give them the whole dashboard
	if r.Header.Get("HX-Request") != "true" {
		b.handleDashboard(w, r)
		return
	}

	// 2. If it's an HTMX click, return ONLY the fishing fragment
	data := struct {
		FishingData interface{}
	}{
		FishingData: struct {
			Endpoint string
			Enabled  bool
		}{
			Endpoint: "/api/toggle-fishing",
			Enabled:  b.fishingEnabled,
		},
	}

	// IMPORTANT: Render the "fishing-view" block, NOT "index.gohtml"
	b.renderTemplate(w, "fishing-view", data)
}

func (b *Bot) handleToggleFishing(w http.ResponseWriter, r *http.Request) {
	// 1. Toggle the logic
	b.fishingEnabled = !b.fishingEnabled

	// 2. Prepare the FULL data the generic toggle needs
	// If you miss 'Endpoint', the next click won't work!
	data := struct {
		Endpoint string
		Enabled  bool
	}{
		Endpoint: "/api/toggle-fishing",
		Enabled:  b.fishingEnabled,
	}

	// 3. Use the correct template name "toggle"
	b.renderTemplate(w, "toggle", data)
}

func (b *Bot) handleStatsView(w http.ResponseWriter, r *http.Request) {
	b.renderTemplate(w, "stats-view", nil)

	b.renderNavItem(w, "fishing", "🎣", "Fishing", false)

	b.renderNavItem(w, "stats", "🛡️", "Stats", true)
}

func (b *Bot) renderNavItem(w http.ResponseWriter, id, icon, label string, active bool) {
	templates.ExecuteTemplate(w, "nav-item", map[string]interface{}{
		"ID":     id,
		"Icon":   icon,
		"Label":  label,
		"Active": active,
		"IsOOB":  true, // Tells HTMX to find the element by ID and swap it
	})
}
