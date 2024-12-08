package handlers

import (
	"github.com/larky/bot"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, c *bot.Commands) {
	mux.HandleFunc("/hello", HelloHandler)
	// mux.HandleFunc("/auth/webhook", WebhookHandler)
	mux.HandleFunc("/auth/webhook", func(w http.ResponseWriter, r *http.Request) {
		WebhookHandler(w, r, c)
	})
	// mux.HandleFunc("/test", TestHandler)
	// mux.HandleFunc("/send-test-msg", SendTestMessageHandler)
}
