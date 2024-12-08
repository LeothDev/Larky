package main

import (
	"fmt"
	"github.com/larky/bot"
	"github.com/larky/handlers"
	"log"
	"net/http"
)

// Initialize basic Larky webapp
func main() {
	commands := bot.NewCommands()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fileServer)
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	handlers.RegisterRoutes(mux, commands)
	// bot.NewBot()

	fmt.Printf("Starting server at port 8080\n")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
