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
