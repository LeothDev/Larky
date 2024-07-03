package bot

import (
	"fmt"
	"github.com/go-lark/lark"
	"os"
)

type CommandsHandler func(bot *lark.Bot)
type Commands struct {
	Commands map[string]CommandsHandler
}

func NewCommands() *Commands {
	return &Commands{
		Commands: map[string]CommandsHandler{
			"!hello": CommandHelloFunc,
		},
	}
}

func CommandHelloFunc(bot *lark.Bot) {
	fmt.Println("I'm in CommandHelloFunc!")
	email := os.Getenv("EMAIL")
	resp, _ := bot.PostText("Hey!", lark.WithEmail(email))
	fmt.Println(resp)
}
