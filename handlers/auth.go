package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type WebhookRequest struct {
	Challenge string `json:"challenge"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Webhook Called!")
	if r.Method != "POST" {
		http.Error(w, "Method Not Supported!", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req WebhookRequest
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}

	challenge := req.Challenge
	log.Println("Challenge", challenge)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
	if err != nil {
		return
	}
}
