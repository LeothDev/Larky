package bot

// TODO: Update HandleEvent
import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-lark/lark"
	"log"
	"net/http"
	"os"
)

// SignatureValidation to verify Lark requests
func SignatureValidation(timestamp, nonce, encryptKey, signature string, bodyBytes []byte) bool {
	b1 := []byte(timestamp + nonce + encryptKey)
	toConcat := [][]byte{
		b1,
		bodyBytes,
	}
	sep := []byte("")
	b := bytes.Join(toConcat, sep)

	h := sha256.New()
	h.Write(b)
	bs := h.Sum(nil)
	s := fmt.Sprintf("%x", bs)

	if s == signature {
		return true
	}
	return false
}

func MsgTest(bot *lark.Bot) error {
	email := os.Getenv("EMAIL")
	_, err := bot.PostText("Testing Larky!", lark.WithEmail(email))
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}
	return nil
}

// getIDandSecret returns APP_ID and APP_SECRET
func getIDandSecret() (string, string) {
	appID := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")
	return appID, appSecret
}

// HandleEvent takes care of the user input and handles the response
func LogicEvent(reqBody json.RawMessage, w http.ResponseWriter, commands *Commands) error {
	handler := NewHandler()
	var e FullEvent
	// TODO: INITIALIZE NEWCOMMANDS IN THE MAIN SERVER FILE
	// commands := NewCommands()
	_ = e.GetEventJSON(reqBody)

	if eventHandler, ok := handler.Handlers[e.Header.EventType]; ok {
		bot := lark.NewChatBot(getIDandSecret())
		bot.SetDomain(lark.DomainLark)
		fmt.Printf("BOT DOMAIN: %s\n\n", bot.Domain())
		_ = bot.StartHeartbeat()
		eventHandler(e, bot, commands)
	} else {
		fmt.Printf("No handler found for event type %s\n", e.Header.EventType)
	}
	fmt.Printf("Header | EventID and EventType: %s, %s\n", e.Header.EventID, e.Header.EventType)
	fmt.Printf("EventBody: %s\n", e.EventBody)
	fmt.Printf("Message Params: %s, %s\n", e.Event.Message.MessageType, e.Event.Message.Content)
	return nil
	// fmt.Printf("Raw JSON: %s\n\n", reqBody)
	// fmt.Printf("Header EventID and EventType: %s, %s\n", e.Header.EventID, e.Header.EventType)
	// fmt.Printf("EventBody: %s\n", e.EventBody)
	// bot := lark.NewChatBot(getIDandSecret())
	// _ = bot.StartHeartbeat()
}
