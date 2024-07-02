package bot

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/go-lark/lark"
	"log"
	"os"
	_ "strings"
)

func Init() (string, string) {
	appID := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")
	return appID, appSecret
}

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
