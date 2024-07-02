package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/larky/utils"
	"log"
	"net/http"
	"os"
)

type WebhookEncrypted struct {
	Encrypt string `json:"encrypt"`
}

type WebhookValidation struct {
	Challenge string `json:"challenge"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Webhook Called!")
	if r.Method != "POST" {
		http.Error(w, "Method Not Supported!", http.StatusMethodNotAllowed)
		return
	}

	// Get encrypted Event Content
	decoder := json.NewDecoder(r.Body)
	var we WebhookEncrypted
	err := decoder.Decode(&we)
	if err != nil {
		panic(err)
	}
	eventContent := we.Encrypt

	encryptKey := os.Getenv("ENCRYPT_KEY")
	content := utils.Decrypt(encryptKey, eventContent)

	// Store the "challenge" string
	var wv WebhookValidation
	err = json.Unmarshal([]byte(content), &wv)
	if err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}
	challenge := wv.Challenge

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{"challenge": challenge})
	if err != nil {
		return
	}

}
