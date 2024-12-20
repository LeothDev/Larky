package lark

import (
	"fmt"
	"io"
	"net/http"
)

const (
	messageURL                = "/open-apis/im/v1/messages?receive_id_type=%s"
	replyMessageURL           = "/open-apis/im/v1/messages/%s/reply"
	reactionsMessageURL       = "/open-apis/im/v1/messages/%s/reactions"
	deleteReactionsMessageURL = "/open-apis/im/v1/messages/%s/reactions/%s"
	getMessageURL             = "/open-apis/im/v1/messages/%s"
	updateMessageURL          = "/open-apis/im/v1/messages/%s"
	recallMessageURL          = "/open-apis/im/v1/messages/%s"
	getMessageResourceURL     = "/open-apis/im/v1/messages/%s/resources/%s?type=%s" // CUSTOM
	messageReceiptURL         = "/open-apis/message/v4/read_info/"
	ephemeralMessageURL       = "/open-apis/ephemeral/v1/send"
	deleteEphemeralMessageURL = "/open-apis/ephemeral/v1/delete"
	pinMessageURL             = "/open-apis/im/v1/pins"
	unpinMessageURL           = "/open-apis/im/v1/pins/%s"
	forwardMessageURL         = "/open-apis/im/v1/messages/%s/forward?receive_id_type=%s"
	downloadFileURL           = "open-apis/im/v1/files/%s" // CUSTOM
)

// PostMessageResponse .
type PostMessageResponse struct {
	BaseResponse

	Data IMMessage `json:"data"`
}

// IMMessageRequest .
type IMMessageRequest struct {
	Content       string `json:"content"`
	MsgType       string `json:"msg_type,omitempty"`
	ReceiveID     string `json:"receive_id,omitempty"`
	UUID          string `json:"uuid,omitempty"`
	ReplyInThread bool   `json:"reply_in_thread,omitempty"`
}

// IMSender .
type IMSender struct {
	ID         string `json:"id"`
	IDType     string `json:"id_type"`
	SenderType string `json:"sender_type"`
	TenantKey  string `json:"tenant_key"`
}

// IMMention .
type IMMention struct {
	ID     string `json:"id"`
	IDType string `json:"id_type"`
	Key    string `json:"key"`
	Name   string `json:"name"`
}

// IMBody .
type IMBody struct {
	Content string `json:"content"`
}

// IMMessage .
type IMMessage struct {
	MessageID      string      `json:"message_id"`
	UpperMessageID string      `json:"upper_message_id"`
	RootID         string      `json:"root_id"`
	ParentID       string      `json:"parent_id"`
	ThreadID       string      `json:"thread_id"`
	ChatID         string      `json:"chat_id"`
	MsgType        string      `json:"msg_type"`
	CreateTime     string      `json:"create_time"`
	UpdateTime     string      `json:"update_time"`
	Deleted        bool        `json:"deleted"`
	Updated        bool        `json:"updated"`
	Sender         IMSender    `json:"sender"`
	Mentions       []IMMention `json:"mentions"`
	Body           IMBody      `json:"body"`
}

// ReactionResponse .
type ReactionResponse struct {
	BaseResponse
	Data struct {
		ReactionID string `json:"reaction_id"`
		Operator   struct {
			OperatorID   string `json:"operator_id"`
			OperatorType string `json:"operator_type"`
			ActionTime   string `json:"action_time"`
		} `json:"operator"`
		ReactionType struct {
			EmojiType EmojiType `json:"emoji_type"`
		} `json:"reaction_type"`
	} `json:"data"`
}

// GetMessageResponse .
type GetMessageResponse struct {
	BaseResponse

	Data struct {
		Items []IMMessage `json:"items"`
	} `json:"data"`
}

// PostEphemeralMessageResponse .
type PostEphemeralMessageResponse struct {
	BaseResponse
	Data struct {
		MessageID string `json:"message_id"`
	} `json:"data"`
}

// DeleteEphemeralMessageResponse .
type DeleteEphemeralMessageResponse = BaseResponse

// RecallMessageResponse .
type RecallMessageResponse = BaseResponse

// GetMessageResourceResponse .
type GetMessageResourceResponse struct {
	Data []byte // Raw binary data
}

// UpdateMessageResponse .
type UpdateMessageResponse = BaseResponse

// ForwardMessageResponse .
type ForwardMessageResponse = PostMessageResponse

// MessageReceiptResponse .
type MessageReceiptResponse struct {
	BaseResponse
	Data struct {
		ReadUsers []struct {
			OpenID    string `json:"open_id"`
			UserID    string `json:"user_id"`
			Timestamp string `json:"timestamp"`
		} `json:"read_users"`
	} `json:"data"`
}

// PinMessageResponse .
type PinMessageResponse struct {
	BaseResponse
	Data struct {
		Pin struct {
			MessageID      string `json:"message_id"`
			ChatID         string `json:"chat_id"`
			OperatorID     string `json:"operator_id"`
			OperatorIDType string `json:"operator_id_type"`
			CreateTime     string `json:"create_time"`
		} `json:"pin"`
	} `json:"data"`
}

// UnpinMessageResponse .
type UnpinMessageResponse = BaseResponse

// DownloadFileResponse .
type DownloadFileResponse struct {
	Data []byte // Raw binary data
}

func newMsgBufWithOptionalUserID(msgType string, userID *OptionalUserID) *MsgBuffer {
	mb := NewMsgBuffer(msgType)
	realID := userID.RealID
	switch userID.UIDType {
	case "email":
		mb.BindEmail(realID)
	case "open_id":
		mb.BindOpenID(realID)
	case "chat_id":
		mb.BindChatID(realID)
	case "user_id":
		mb.BindUserID(realID)
	case "union_id":
		mb.BindUnionID(realID)
	default:
		return nil
	}
	return mb
}

// PostText is a simple way to send text messages
func (bot Bot) PostText(text string, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgText, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	om := mb.Text(text).Build()
	return bot.PostMessage(om)
}

// PostRichText is a simple way to send rich text messages
func (bot Bot) PostRichText(postContent *PostContent, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgPost, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	om := mb.Post(postContent).Build()
	return bot.PostMessage(om)
}

// PostTextMention is a simple way to send text messages with @user
func (bot Bot) PostTextMention(text string, atUserID string, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgText, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	tb := NewTextBuilder()
	om := mb.Text(tb.Text(text).Mention(atUserID).Render()).Build()
	return bot.PostMessage(om)
}

// PostTextMentionAll is a simple way to send text messages with @all
func (bot Bot) PostTextMentionAll(text string, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgText, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	tb := NewTextBuilder()
	om := mb.Text(tb.Text(text).MentionAll().Render()).Build()
	return bot.PostMessage(om)
}

// PostTextMentionAndReply is a simple way to send text messages with @user and reply a message
func (bot Bot) PostTextMentionAndReply(text string, atUserID string, userID *OptionalUserID, replyID string) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgText, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	tb := NewTextBuilder()
	om := mb.Text(tb.Text(text).Mention(atUserID).Render()).BindReply(replyID).Build()
	return bot.PostMessage(om)
}

// PostImage is a simple way to send image
func (bot Bot) PostImage(imageKey string, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgImage, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	om := mb.Image(imageKey).Build()
	return bot.PostMessage(om)
}

// PostShareChat is a simple way to share chat
func (bot Bot) PostShareChat(chatID string, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgShareCard, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	om := mb.ShareChat(chatID).Build()
	return bot.PostMessage(om)
}

// PostShareUser is a simple way to share user
func (bot Bot) PostShareUser(openID string, userID *OptionalUserID) (*PostMessageResponse, error) {
	mb := newMsgBufWithOptionalUserID(MsgShareUser, userID)
	if mb == nil {
		return nil, ErrParamUserID
	}
	om := mb.ShareUser(openID).Build()
	return bot.PostMessage(om)
}

// PostMessage posts a message
func (bot Bot) PostMessage(om OutcomingMessage) (*PostMessageResponse, error) {
	req, err := BuildMessage(om)
	if err != nil {
		return nil, err
	}
	var respData PostMessageResponse
	if om.RootID == "" {
		err = bot.PostAPIRequest("PostMessage", fmt.Sprintf(messageURL, om.UIDType), true, req, &respData)
	} else {
		resp, err := bot.ReplyMessage(om)
		return resp, err
	}
	return &respData, err
}

// ReplyMessage replies a message
func (bot Bot) ReplyMessage(om OutcomingMessage) (*PostMessageResponse, error) {
	req, err := buildReplyMessage(om)
	if err != nil {
		return nil, err
	}
	if om.RootID == "" {
		return nil, ErrParamMessageID
	}
	var respData PostMessageResponse
	err = bot.PostAPIRequest("ReplyMessage", fmt.Sprintf(replyMessageURL, om.RootID), true, req, &respData)
	return &respData, err
}

// AddReaction adds reaction to a message
func (bot Bot) AddReaction(messageID string, emojiType EmojiType) (*ReactionResponse, error) {
	req := map[string]interface{}{
		"reaction_type": map[string]interface{}{
			"emoji_type": emojiType,
		},
	}
	var respData ReactionResponse
	err := bot.PostAPIRequest("AddReaction", fmt.Sprintf(reactionsMessageURL, messageID), true, req, &respData)
	return &respData, err
}

// DeleteReaction deletes reaction of a message
func (bot Bot) DeleteReaction(messageID string, reactionID string) (*ReactionResponse, error) {
	var respData ReactionResponse
	err := bot.DeleteAPIRequest("DeleteReaction", fmt.Sprintf(deleteReactionsMessageURL, messageID, reactionID), true, nil, &respData)
	return &respData, err
}

// UpdateMessage updates a message
func (bot Bot) UpdateMessage(messageID string, om OutcomingMessage) (*UpdateMessageResponse, error) {
	if om.MsgType != MsgInteractive &&
		om.MsgType != MsgText &&
		om.MsgType != MsgPost {
		return nil, ErrMessageType
	}
	req, err := buildUpdateMessage(om)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(updateMessageURL, messageID)
	var respData UpdateMessageResponse
	if om.MsgType == MsgInteractive {
		err = bot.PatchAPIRequest("UpdateMessage", url, true, req, &respData)
	} else {
		err = bot.PutAPIRequest("UpdateMessage", url, true, req, &respData)
	}
	return &respData, err
}

// GetMessage gets a message with im/v1
func (bot Bot) GetMessage(messageID string) (*GetMessageResponse, error) {
	var respData GetMessageResponse
	err := bot.GetAPIRequest("GetMessage", fmt.Sprintf(getMessageURL, messageID), true, nil, &respData)
	return &respData, err
}

// RecallMessage recalls a message with ID
func (bot Bot) RecallMessage(messageID string) (*RecallMessageResponse, error) {
	url := fmt.Sprintf(recallMessageURL, messageID)
	var respData RecallMessageResponse
	err := bot.DeleteAPIRequest("RecallMessage", url, true, nil, &respData)
	return &respData, err
}

// GetMessageResource gets a resource given a message ID
func (bot Bot) GetMessageResource(messageID, fileKey string) (*GetMessageResourceResponse, error) {
	url := fmt.Sprintf(getMessageResourceURL, messageID, fileKey, "file")
	/*params := map[string]interface{}{
		"type": "file",
	}*/

	// Directly call DoAPIRequest to handle raw binary data
	header := make(http.Header)
	header.Set("Content-Type", "application/json; charset=utf-8")
	if bot.TenantAccessToken() != "" {
		header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.TenantAccessToken()))
	}

	// Perform the request and read raw response data
	req, err := http.NewRequest(http.MethodGet, bot.ExpandURL(url), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = header

	resp, err := bot.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Return raw data in the response struct
	return &GetMessageResourceResponse{Data: data}, nil

	/*
		var respData GetMessageResourceResponse
		err := bot.GetAPIRequest("GetMessageResource", url, true, params, &respData)
		return &respData, err
		// https://github.com/larksuite/oapi-sdk-go/blob/3ab77a695d7494f3b7864c294783d269b1d687d8/service/im/v1/resource.go
	*/
}

// MessageReadReceipt queries message read receipt
func (bot Bot) MessageReadReceipt(messageID string) (*MessageReceiptResponse, error) {
	params := map[string]string{
		"message_id": messageID,
	}
	var respData MessageReceiptResponse
	err := bot.PostAPIRequest("MessageReadReceipt", messageReceiptURL, true, params, &respData)
	return &respData, err
}

// PostEphemeralMessage posts an ephemeral message
func (bot Bot) PostEphemeralMessage(om OutcomingMessage) (*PostEphemeralMessageResponse, error) {
	if om.UIDType == UIDUnionID {
		return nil, ErrUnsupportedUIDType
	}
	params := BuildOutcomingMessageReq(om)
	var respData PostEphemeralMessageResponse
	err := bot.PostAPIRequest("PostEphemeralMessage", ephemeralMessageURL, true, params, &respData)
	return &respData, err
}

// DeleteEphemeralMessage deletes an ephemeral message
func (bot Bot) DeleteEphemeralMessage(messageID string) (*DeleteEphemeralMessageResponse, error) {
	params := map[string]interface{}{
		"message_id": messageID,
	}
	var respData DeleteEphemeralMessageResponse
	err := bot.PostAPIRequest("DeleteEphemeralMessage", deleteEphemeralMessageURL, true, params, &respData)
	return &respData, err
}

// PinMessage pins a message
func (bot Bot) PinMessage(messageID string) (*PinMessageResponse, error) {
	params := map[string]interface{}{
		"message_id": messageID,
	}
	var respData PinMessageResponse
	err := bot.PostAPIRequest("PinMessage", pinMessageURL, true, params, &respData)
	return &respData, err
}

// UnpinMessage unpins a message
func (bot Bot) UnpinMessage(messageID string) (*UnpinMessageResponse, error) {
	url := fmt.Sprintf(unpinMessageURL, messageID)
	var respData UnpinMessageResponse
	err := bot.DeleteAPIRequest("PinMessage", url, true, nil, &respData)
	return &respData, err
}

// ForwardMessage forwards a message
func (bot Bot) ForwardMessage(messageID string, receiveID *OptionalUserID) (*ForwardMessageResponse, error) {
	url := fmt.Sprintf(forwardMessageURL, messageID, receiveID.UIDType)
	params := map[string]interface{}{
		"receive_id": receiveID.RealID,
	}
	var respData ForwardMessageResponse
	err := bot.PostAPIRequest("ForwardMessage", url, true, params, &respData)
	return &respData, err
}

// TODO: FIX THIS
func (bot Bot) DownloadFile(fileKey string) (*DownloadFileResponse, error) {
	url := fmt.Sprintf(downloadFileURL, fileKey)

	// Directly call DoAPIRequest to handle raw binary data
	header := make(http.Header)
	header.Set("Content-Type", "application/json; charset=utf-8")
	if bot.TenantAccessToken() != "" {
		header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.TenantAccessToken()))
	}

	// Perform the request and read raw response data
	req, err := http.NewRequest(http.MethodGet, bot.ExpandURL(url), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header = header

	resp, err := bot.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Return raw data in the response struct
	return &DownloadFileResponse{Data: data}, nil
}
