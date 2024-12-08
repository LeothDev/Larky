package bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-lark/lark"
	"github.com/larky/utils"
	"log"
	"strings"
)

type Header struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
}

type SenderID struct {
	UserID string `json:"user_id"`
}

type Sender struct {
	SenderID SenderID `json:"sender_id"`
}

type Message struct {
	ChatID      string `json:"chat_id"`
	ChatType    string `json:"chat_type"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	// TODO: Deal with MessageType to differentiate between messages
	MessageID string `json:"message_id"`
}

type Event struct {
	Message Message `json:"message"`
	Sender  Sender  `json:"sender"`
}

type FullEvent struct {
	EventBody json.RawMessage
	Header    Header `json:"header"`
	Event     Event  `json:"event"`
	Schema    string `json:"schema"`
}

func (e *FullEvent) GetEventJSON(reqBody json.RawMessage) error {
	if err := json.Unmarshal(reqBody, e); err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}
	e.EventBody = reqBody
	return nil
}

// type EventsHandler func(event FullEvent, bot *lark.Bot, commands *Commands)
type EventsHandler func(event FullEvent, bot *lark.Bot, commands *Commands)
type Handler struct {
	Handlers map[string]EventsHandler
}

func NewHandler() *Handler {
	return &Handler{
		Handlers: map[string]EventsHandler{
			"im.message.receive_v1":                        HandleEventMessageReceived,
			"im.chat.access_event.bot_p2p_chat_entered_v1": HandleEventUserEnterChatWithBot,
			"p2p_chat_create":                              HandleEventUserAndBotFirstTimeChat,
			"im.chat.member.bot.added_v1":                  HandleEventBotAddedToGroup,
			"im.chat.member.user.added_v1":                 HandleEventUserAddedToGroup,
			"im.message.message_read_v1":                   HandleEventMessageRead,
			"application.bot.menu_v6":                      HandleEventBotMenu,
		},
	}
}

func HandleEventMessageReceived(event FullEvent, bot *lark.Bot, commands *Commands) {
	// botHandler := NewCommands()
	fmt.Printf("I'm handling %s\n", event.Header.EventType)
	content := event.Event.Message.Content
	userID := event.Event.Sender.SenderID.UserID

	// Check the user's session state
	state, exists := commands.GetSession(userID)
	if !exists {
		fmt.Println("No active session for user.")
	}

	extractedMsgContent := utils.ExtractContent(content)

	// Handle based on the state
	if state == "awaiting_xlsx" {
		if event.Event.Message.MessageType == "file" {
			fileKey, _ := utils.ExtractFileMsgContents(content)
			accessToken := bot.TenantAccessToken()
			messageID := event.Event.Message.MessageID
			bot.PostText("Your .xlsx file has been received! Processing...", lark.WithUserID(userID))

			RetrieveFile(bot, accessToken, messageID, fileKey)
			commands.ClearSession(userID)
		} else {
			commands.ClearSession(userID)
		}
		// TODO: Check if the file is an .xlsx
		// fmt.Printf("Exists | State : %t, %s\n", exists, state)
		// fmt.Printf("Content: %s\n", content)
		// fmt.Printf("ExtractedMsgContent: %s\n", extractedMsgContent)
		if len(extractedMsgContent) == 0 {
			// That means that the message received is not a text, but rather a file (supposedly)
		}
	} else {
		if !strings.HasPrefix(extractedMsgContent, "!") {
			return
		}
		if command, ok := commands.Commands[extractedMsgContent]; ok {
			command(bot, userID, content, commands)
		} else {
			fmt.Println("Unknown Command")
			return
		}
	}

}

func HandleEventUserEnterChatWithBot(event FullEvent, bot *lark.Bot, commands *Commands) {
}

func HandleEventUserAndBotFirstTimeChat(event FullEvent, bot *lark.Bot, commands *Commands) {
}

func HandleEventBotAddedToGroup(event FullEvent, bot *lark.Bot, commands *Commands) {
}

func HandleEventUserAddedToGroup(event FullEvent, bot *lark.Bot, commands *Commands) {
}

func HandleEventMessageRead(event FullEvent, bot *lark.Bot, commands *Commands) {
}

func HandleEventBotMenu(event FullEvent, bot *lark.Bot, commands *Commands) {
}
