package bot

import (
	"fmt"
	"github.com/go-lark/lark"
	"log"
	"os"
)

func Init() (string, string) {
	appID := os.Getenv("APP_ID")
	appSecret := os.Getenv("APP_SECRET")
	return appID, appSecret
}

func MsgTest(bot *lark.Bot) error {
	fmt.Println("Sending Test Message to User!")
	email := os.Getenv("EMAIL")
	_, err := bot.PostText("Testing Larky!", lark.WithEmail(email))
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
		return err
	}
	return nil
}
