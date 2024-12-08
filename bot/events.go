package bot

import (
	"fmt"
	"github.com/go-lark/lark"
)

func DummyProcessXLSX() {
	fmt.Println("I'M IN DUMMY XLSX METHOD!")
}

func ProcessXcelFile(bot *lark.Bot, userID, content string) {
	// fmt.Println("User triggered !cleanxcel command")
}

// (*bytes.Buffer, error)
func RetrieveFile(bot *lark.Bot, accessToken, messageID, fileKey string) {
	resp, _ := bot.GetMessage(messageID)
	fmt.Printf("Get Message RESPONSE: %s", resp)
	// TODO: Check UploadFileRequest struct
	// TODO: Add getMessageResource method in api_message.go by using bot.GetAPIRequest
	// var respData lark.GetMessageResponse
	// bot.GetAPIRequest("GetMessageResource", fmt.Sprintf(getMessageResourceURL, messageID, fileKey), true, nil, &respData)
}
