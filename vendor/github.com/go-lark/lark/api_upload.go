package lark

// Import from Lark API Go demo
// with adaption to go-lark frame

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	uploadImageURL        = "/open-apis/im/v1/images"
	uploadFileURL         = "/open-apis/im/v1/files"
	uploadFilesURL        = "/open-apis/drive/v1/files/upload_all" // CUSTOM
	importSpreadsheetURL  = "/open-apis/sheets/v2/import"          // CUSTOM
	queryImportResultsURL = "/open-apis/sheets/v2/import/result"   // CUSTOM
)

// UploadImageResponse .
type UploadImageResponse struct {
	BaseResponse
	Data struct {
		ImageKey string `json:"image_key"`
	} `json:"data"`
}

// UploadFileRequest .
type UploadFileRequest struct {
	FileType string    `json:"-"`
	FileName string    `json:"-"`
	Duration int       `json:"-"`
	Path     string    `json:"-"`
	Reader   io.Reader `json:"-"`
}

// UploadFileResponse .
type UploadFileResponse struct {
	BaseResponse
	Data struct {
		FileKey string `json:"file_key"`
	} `json:"data"`
}

// UploadFilesRequest . CUSTOM
type UploadFilesRequest struct {
	FileName string `json:"-"` // The name of the file
	// FilePath   string    `json:"-"` // If Reader is nil
	File       []byte `json:"-"` // File content as a []byte
	ParentType string `json:"-"` // File type, e.g., "explorer" (My Space)
	ParentNode string `json:"-"` // Parent node ID (Folder Token)
}

// UploadFilesResponse . CUSTOM
type UploadFilesResponse struct {
	BaseResponse
	Data struct {
		FileToken string `json:"file_token"`
	} `json:"data"`
}

// ImportSpreadsheetRequest . CUSTOM
type ImportSpreadsheetRequest struct {
	File        *bytes.Buffer `json:"-"` // File as a buffer
	FileName    string        `json:"-"` // Name of the file (e.g., "test.xlsx")
	FolderToken string        `json:"-"` // Folder where the file will be stored
}

// ImportSpreadsheetResponse . CUSTOM
type ImportSpreadsheetResponse struct {
	BaseResponse
	Data struct {
		Ticket string `json:"ticket"` // Import ticket
	} `json:"data"`
}

// UploadImage uploads image to Lark server
func (bot Bot) UploadImage(path string) (*UploadImageResponse, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("image_type", "message")
	part, err := writer.CreateFormFile("image", path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	var respData UploadImageResponse
	header := make(http.Header)
	header.Set("Content-Type", writer.FormDataContentType())
	err = bot.DoAPIRequest("POST", "UploadImage", uploadImageURL, header, true, body, &respData)
	if err != nil {
		return nil, err
	}
	return &respData, err
}

// UploadImageObject uploads image to Lark server
func (bot Bot) UploadImageObject(img image.Image) (*UploadImageResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("image_type", "message")
	part, err := writer.CreateFormFile("image", "temp_image")
	if err != nil {
		return nil, err
	}
	err = jpeg.Encode(part, img, nil)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	var respData UploadImageResponse
	header := make(http.Header)
	header.Set("Content-Type", writer.FormDataContentType())
	err = bot.DoAPIRequest("POST", "UploadImage", uploadImageURL, header, true, body, &respData)
	if err != nil {
		return nil, err
	}
	return &respData, err
}

// UploadFile uploads file to Lark server
func (bot Bot) UploadFile(req UploadFileRequest) (*UploadFileResponse, error) {
	var content io.Reader
	if req.Reader == nil {
		file, err := os.Open(req.Path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		content = file
	} else {
		content = req.Reader
		req.Path = req.FileName
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("file_type", req.FileType)
	writer.WriteField("file_name", req.FileName)
	if req.FileType == "mp4" && req.Duration > 0 {
		writer.WriteField("duration", fmt.Sprintf("%d", req.Duration))
	}
	part, err := writer.CreateFormFile("file", req.Path)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, content)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	var respData UploadFileResponse
	header := make(http.Header)
	header.Set("Content-Type", writer.FormDataContentType())
	err = bot.DoAPIRequest("POST", "UploadFile", uploadFileURL, header, true, body, &respData)
	if err != nil {
		return nil, err
	}
	return &respData, err
}

func (bot Bot) UploadFiles(req UploadFilesRequest) (*UploadFilesResponse, error) {
	content := bytes.NewReader(req.File)
	fileSize := fmt.Sprintf("%d", len(req.File))

	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	if err := writer.WriteField("file_name", req.FileName); err != nil {
		return nil, fmt.Errorf("failed to write file_name field: %w", err)
	}
	if err := writer.WriteField("parent_type", req.ParentType); err != nil {
		return nil, fmt.Errorf("failed to write parent_type field: %w", err)
	}
	if err := writer.WriteField("parent_node", req.ParentNode); err != nil {
		return nil, fmt.Errorf("failed to write parent_node field: %w", err)
	}
	if err := writer.WriteField("size", fileSize); err != nil {
		return nil, fmt.Errorf("failed to write size field: %w", err)
	}

	// Add file content
	part, err := writer.CreateFormFile("file", req.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, content); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}

	// Close the multipart writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Prepare request headers
	// Directly call DoAPIRequest to handle raw binary data
	header := make(http.Header)
	header.Set("Content-Type", "multipart/form-data")
	if bot.TenantAccessToken() != "" {
		header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.TenantAccessToken()))
	}

	// Perform the API request
	var respData UploadFilesResponse
	err = bot.DoAPIRequest("POST", "UploadFiles", uploadFilesURL, header, true, body, &respData)
	if err != nil {
		return nil, err
	}

	return &respData, nil
}

func (bot Bot) ImportSpreadsheet(req ImportSpreadsheetRequest) (*ImportSpreadsheetResponse, error) {
	// Validate input
	if req.File == nil {
		return nil, fmt.Errorf("file buffer cannot be nil")
	}
	if req.FileName == "" {
		return nil, fmt.Errorf("file name cannot be empty")
	}
	if req.FolderToken == "" {
		return nil, fmt.Errorf("folder token cannot be empty")
	}

	// Convert file buffer to byte array
	fileBytes := req.File.Bytes()

	// Prepare the payload
	payload := map[string]interface{}{
		"file":        fileBytes,
		"name":        req.FileName,
		"folderToken": req.FolderToken,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Set headers
	headers := make(http.Header)
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", bot.TenantAccessToken()))
	headers.Set("Content-Type", "application/json; charset=utf-8")

	/*// Perform the API Request
	var respData ImportSpreadsheetResponse
	err = bot.DoAPIRequest("POST", "ImportSpreadSheet", importSpreadsheetURL, headers, true, bytes.NewReader(payloadBytes), &respData)
	if err != nil {
		return nil, err
	}*/

	// Perform the request manually
	newReq, err := http.NewRequest(http.MethodPost, bot.ExpandURL(importSpreadsheetURL), bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	newReq.Header = headers

	resp, err := bot.client.Do(newReq)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	content, _ := io.ReadAll(resp.Body)
	fmt.Println("RESPONSE BODY: ", string(content))

	var respData ImportSpreadsheetResponse
	err = json.Unmarshal(content, &respData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}
	return &respData, nil
}
