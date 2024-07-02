package handlers

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/auth/webhook", WebhookHandler)
	// mux.HandleFunc("/test", TestHandler)
	// mux.HandleFunc("/send-test-msg", SendTestMessageHandler)
}
