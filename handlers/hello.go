package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("helloHandler was called!")
	if r.URL.Path != "/hello" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method Not Supported", http.StatusMethodNotAllowed)
		return
	} else {
		_, err := fmt.Fprintf(w, "Hello Larky!")
		if err != nil {
			log.Fatal(err)
		}
	}
}
