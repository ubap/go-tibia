package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (b *Bot) loopWebUI() {
	b.wg.Add(1)
	defer b.wg.Done()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", b.HandleWS)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("[UI] Dashboard live at http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[UI] HTTP server error: %v", err)
		}
	}()

	<-b.stopChan

	log.Println("[UI] Shutting down server...")

	// Create a context with a timeout so it doesn't hang forever
	// if a browser tab stays connected
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[UI] Shutdown error: %v", err)
	}
	log.Println("[UI] Server stopped.")
}
