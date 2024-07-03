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

type Message struct {
	ChatID      string `json:"chat_id"`
	ChatType    string `json:"chat_type"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
}

type Event struct {
	Message Message `json:"message"`
}

type FullEvent struct {
	Header Header `json:"header"`
	// EventRequest map[string]interface{} `json:"event"`
	EventBody json.RawMessage
	Event     Event `json:"event"`
}

func (e *FullEvent) GetEventJSON(reqBody json.RawMessage) error {
	if err := json.Unmarshal(reqBody, e); err != nil {
		log.Fatalf("Unable to marshal JSON due to %s", err)
	}
	e.EventBody = reqBody
	return nil
}

type EventsHandler func(event FullEvent, bot *lark.Bot)
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

func HandleEventMessageReceived(event FullEvent, bot *lark.Bot) {
	commands := NewCommands()
	fmt.Printf("I'm handling %s\n", event.Header.EventType)
	content := event.Event.Message.Content
	extractedMsgContent := utils.ExtractContent(content)
	if len(extractedMsgContent) == 0 {
		return
	}

	if !strings.HasPrefix(extractedMsgContent, "!") {
		return
	}
	if command, ok := commands.Commands[extractedMsgContent]; ok {
		command(bot)
	} else {
		fmt.Println("Unknown Command")
		return
	}
}

func HandleEventUserEnterChatWithBot(event FullEvent, bot *lark.Bot) {
}

func HandleEventUserAndBotFirstTimeChat(event FullEvent, bot *lark.Bot) {
}

func HandleEventBotAddedToGroup(event FullEvent, bot *lark.Bot) {
}

func HandleEventUserAddedToGroup(event FullEvent, bot *lark.Bot) {
}

func HandleEventMessageRead(event FullEvent, bot *lark.Bot) {
}

func HandleEventBotMenu(event FullEvent, bot *lark.Bot) {
}
