package handlers

/*
// TestHandler to test if Bot is UP Running and serves the test.html file
func TestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("static", "test.html"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func SendTestMessageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The Test Button has been Clicked!")
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	larkyBot := lark.NewChatBot(bot.Init())
	_ = larkyBot.StartHeartbeat()
	err := bot.MsgTest(larkyBot)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"success": "Message Sent Successfully!"}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}
*/
