package bot

import (
	"bytes"
	"fmt"
	"github.com/go-lark/lark"
	"github.com/xuri/excelize/v2"
	"io"
	"strings"
)

func DummyProcessXLSX() {
	fmt.Println("I'M IN DUMMY XLSX METHOD!")
}

func ProcessXcelFile(file io.Reader, fileName string) (*bytes.Buffer, string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}

	// Read all rows from Sheet1
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, "", fmt.Errorf("failed to read rows: %w", err)
	}

	// Locate required columns: strId, EN, Italian
	var strIDIdx, enIdx, italianIdx int
	if len(rows) > 0 {
		headers := rows[0]
		for i, col := range headers {
			switch col {
			case "strId":
				strIDIdx = i
			case "EN":
				enIdx = i
			case "Italian":
				italianIdx = i
			}
		}
	}

	// Filter rows to keep only strId, EN, and Italian columns
	var filteredRows [][]string
	filteredRows = append(filteredRows, []string{"strId", "EN", "Italian"}) // Add header row
	for _, row := range rows[1:] {                                          // Skip the header row
		newRow := []string{"", "", ""}
		if strIDIdx < len(row) {
			newRow[0] = row[strIDIdx]
		}
		if enIdx < len(row) {
			newRow[1] = row[enIdx]
		}
		if italianIdx < len(row) {
			newRow[2] = row[italianIdx]
		}
		filteredRows = append(filteredRows, newRow)
	}

	// Map to track unique rows and their occurrences
	type rowKey struct {
		EN      string
		Italian string
	}
	uniqueRows := map[rowKey][]string{}
	occurrences := map[rowKey]int{}

	// Process filtered rows
	for _, row := range filteredRows[1:] { // Skip header
		key := rowKey{
			EN:      strings.ReplaceAll(row[1], "&lt;", "<"),
			Italian: strings.ReplaceAll(row[2], "&gt;", ">"),
		}
		key.EN = strings.ReplaceAll(key.EN, "&gt;", ">")
		key.Italian = strings.ReplaceAll(key.Italian, "&lt;", "<")
		occurrences[key]++
		if _, exists := uniqueRows[key]; !exists {
			uniqueRows[key] = row
		}
	}

	// Create a new Excel file
	output := excelize.NewFile()
	outputSheet := "Sheet1"
	output.NewSheet(outputSheet)

	// Write headers
	headers := []string{"strId", "EN", "Italian", "Occurrences"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		output.SetCellValue(outputSheet, cell, header)
	}

	// Write unique rows with occurrences
	rowIndex := 2
	for key, originalRow := range uniqueRows {
		// Ensure the `strId` value comes directly from the original row.
		output.SetCellValue(outputSheet, fmt.Sprintf("A%d", rowIndex), originalRow[0])   // strId
		output.SetCellValue(outputSheet, fmt.Sprintf("B%d", rowIndex), key.EN)           // EN
		output.SetCellValue(outputSheet, fmt.Sprintf("C%d", rowIndex), key.Italian)      // Italian
		output.SetCellValue(outputSheet, fmt.Sprintf("D%d", rowIndex), occurrences[key]) // Occurrences
		rowIndex++
	}

	// Set column widths
	output.SetColWidth(outputSheet, "A", "A", 60) // strId
	output.SetColWidth(outputSheet, "B", "C", 60) // EN and Italian
	output.SetColWidth(outputSheet, "D", "D", 10) // Occurrences

	// Set text wrapping and alignment for strId, EN, Italian
	styleText, err := output.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText:   true,
			Vertical:   "center",
			Horizontal: "left",
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create text wrap style: %w", err)
	}
	output.SetCellStyle(outputSheet, "A1", fmt.Sprintf("C%d", rowIndex-1), styleText)

	// Set center alignment for Occurrences column
	styleCenter, err := output.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create center alignment style: %w", err)
	}
	output.SetCellStyle(outputSheet, "D1", fmt.Sprintf("D%d", rowIndex-1), styleCenter)

	// Set "freeze first row" panes
	err = output.SetPanes(outputSheet, &excelize.Panes{
		Freeze:      true,
		TopLeftCell: "A2", // The first cell of the unfrozen area (row after the header)
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to set panes: %w", err)
	}

	// Save to a buffer
	var buffer bytes.Buffer
	if err := output.Write(&buffer); err != nil {
		return nil, "", fmt.Errorf("failed to write output: %w", err)
	}

	newFileName := fmt.Sprintf("%s_Unique.xlsx", fileName)
	return &buffer, newFileName, nil
}

// (*bytes.Buffer, error)
func RetrieveFile(bot *lark.Bot, accessToken, messageID, fileKey string) []byte {
	resp, _ := bot.GetMessageResource(messageID, fileKey)
	// fmt.Printf("Get Message RESPONSE: %x\n\n", resp.Data)
	// resp, _ := bot.DownloadFile(fileKey)
	// fmt.Printf("Get Message RESPONSE: %x\n\n", resp.Data)
	file := resp.Data
	return file

	// TODO: Check UploadFileRequest struct
	// TODO: Add getMessageResource method in api_message.go by using bot.GetAPIRequest
	// var respData lark.GetMessageResponse
	// bot.GetAPIRequest("GetMessageResource", fmt.Sprintf(getMessageResourceURL, messageID, fileKey), true, nil, &respData)
}
