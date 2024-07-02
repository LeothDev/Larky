package handlers

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/larky/bot"
	"github.com/larky/utils"
	"io"
	"log"
	_ "log"
	"net/http"
	"os"
	"time"
)

// WebhookValidation is a struct holding important request information to parse
type WebhookValidation struct {
	EncryptKey   string
	Timestamp    string
	Nonce        string
	Signature    string
	RequestType  string
	BodyBytes    []byte
	Encrypt      string `json:"encrypt"`
	Challenge    string `json:"challenge"`
	EventContent string `json:"event"`
}

// newWebhookValidation is a constructor function that initialises RequestType to "Unknown"
func newWebhookValidation() *WebhookValidation {
	return &WebhookValidation{
		EncryptKey:   "",
		Timestamp:    "",
		Nonce:        "",
		Signature:    "",
		RequestType:  "Unknown",
		BodyBytes:    nil,
		Encrypt:      "",
		Challenge:    "",
		EventContent: "",
	}
}

// routeWebhook sets WebhookValidation RequestType to either "Event" or "Challenge"
func (wv *WebhookValidation) routeWebhook(headers http.Header) {
	if headers.Get("X-Lark-Signature") != "" {
		wv.Timestamp = headers.Get("X-Lark-Request-Timestamp")
		wv.Nonce = headers.Get("X-Lark-Request-Nonce")
		wv.Signature = headers.Get("X-Lark-Signature")
	} else {
		wv.RequestType = "Verification"
	}
}

// verificationStep simply unmarshals "encrypt" into the struct
func (wv *WebhookValidation) verificationStep(w http.ResponseWriter) {
	if err := json.Unmarshal(wv.BodyBytes, wv); err != nil {
		http.Error(w, "Unknown Request", http.StatusBadRequest)
		return
	}
	wv.RequestType = "Verification"

}

// verificationChallenge decrypts the message, retrieves the challenge code and returns it
// to the Larksuite server
func (wv *WebhookValidation) verificationChallenge(w http.ResponseWriter) {
	encryptedContent := wv.Encrypt
	decryptedContent := utils.Decrypt(wv.EncryptKey, encryptedContent)
	if err := json.Unmarshal([]byte(decryptedContent), wv); err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"challenge": wv.Challenge})
	if err != nil {
		return
	}
}

// eventStep takes care of handling events coming from user interactions
func (wv *WebhookValidation) eventStep(w http.ResponseWriter) {
	wv.RequestType = "Event"
	w.WriteHeader(http.StatusOK)
}

// WebhookHandler handles all the requests to the endpoint 'auth/webhook'
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received request at %v from %v\n", time.Now(), r.RemoteAddr)
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	wv := newWebhookValidation()

	headers := r.Header
	wv.routeWebhook(headers)
	wv.EncryptKey = os.Getenv("ENCRYPT_KEY")

	wv.BodyBytes, _ = io.ReadAll(r.Body)
	if len(wv.BodyBytes) == 0 {
		http.Error(w, "Bad Request, Who are you?", http.StatusBadRequest)
		return
	}

	fmt.Printf("BodyBytes: %s\n", wv.BodyBytes)
	if wv.Signature != "" {
		isValidRequest := bot.SignatureValidation(wv.Timestamp, wv.Nonce, wv.EncryptKey, wv.Signature, wv.BodyBytes)
		if isValidRequest {
			wv.eventStep(w)
		} else {
			http.Error(w, "Bad Request, Failed Signature Validation", http.StatusBadRequest)
		}
	} else {
		wv.verificationStep(w)
		wv.verificationChallenge(w)
	}
}
