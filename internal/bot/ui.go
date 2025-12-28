package bot

import (
	"fmt"
	"net/http"
)

func (b *Bot) loopWebUI() {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", b.HandleWS)

	fmt.Println("[UI] Dashboard live at http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
