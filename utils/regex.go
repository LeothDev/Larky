package utils

import (
	"fmt"
	"log"
	"regexp"
)

func ExtractContent(content string) string {
	pattern := `{"text":"([^"]+)"}`
	re, err := regexp.Compile(pattern)
	if err != nil {
		// fmt.Println("Error compiling regex... ", err)
		log.Fatal("Error compiling regex...")
	}
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	} else {
		fmt.Println("Couldn't extract the content from the message")
		return ""
	}
}

func ExtractFileMsgContents(content string) (string, string) {
	// Define a regex pattern with two capturing groups for file_key and file_name
	pattern := `{"file_key":"([^"]+)","file_name":"([^"]+)"}`

	// Compile the regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Error compiling regex: ", err)
	}

	// Find matches in the content string
	matches := re.FindStringSubmatch(content)

	// Check if both file_key and file_name are captured
	if len(matches) == 3 {
		return matches[1], matches[2] // Return file_key and file_name
	}

	// If no matches are found, return empty strings
	fmt.Println("Couldn't extract file_key and file_name from the message")
	return "", ""
}
