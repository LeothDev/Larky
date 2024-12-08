package bot

import (
	"fmt"
	"github.com/go-lark/lark"
	"sync"
)

// type CommandsHandler func(bot *lark.Bot, botHandler *Commands, userID, content string)
type CommandsHandler func(bot *lark.Bot, userID, content string, commands *Commands)

type Commands struct {
	Commands map[string]CommandsHandler
	Sessions map[string]string
	mu       sync.Mutex
}

func NewCommands() *Commands {
	return &Commands{
		Commands: map[string]CommandsHandler{
			"!hello":     CommandHelloFunc,
			"!cleanxcel": CommandCleanXcelFunc,
		},
		Sessions: map[string]string{},
	}
}

func (c *Commands) SetSession(userID, state string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Sessions[userID] = state
}

func (c *Commands) GetSession(userID string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	state, exists := c.Sessions[userID]
	return state, exists
}

func (c *Commands) ClearSession(userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.Sessions, userID)
}

func CommandHelloFunc(bot *lark.Bot, userID, content string, commands *Commands) {
	fmt.Println("I'm in CommandHelloFunc!")
	// email := os.Getenv("EMAIL")
	resp, _ := bot.PostText("Hey!", lark.WithUserID(userID))
	fmt.Println(resp)
}

func CommandCleanXcelFunc(bot *lark.Bot, userID, content string, commands *Commands) {
	// var botHandler *Commands
	// Check if botHandler or its Sessions are nil
	if commands == nil {
		fmt.Println("botHandler is nil!")
	}
	if commands.Sessions == nil {
		fmt.Println("The session is nil!")
	}

	fmt.Println("I'm in CommandCleanXcelFunc!")

	// Set the user state to expect a xlsx file
	commands.SetSession(userID, "awaiting_xlsx")
	session, _ := commands.GetSession(userID)
	fmt.Printf("Sessions for user %s: %s\n", userID, session)

	// Notify the user
	_, _ = bot.PostText("Please send the .xlsx file to clean.", lark.WithUserID(userID))
}
