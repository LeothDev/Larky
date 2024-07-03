### Key Concepts

- Handle Events by Myself
- Send Bot Messages through go-lark

### General - To Implement (Not in Order) 
- [ ] DB for active chats
- [ ] Button Actions
- [ ] Redirect for non-Lark Users

### Bot Features (Not in Order)
- [ ] Excel Duplicates Cleaning (Based on Language)
- [ ] Set a timer
- [ ] Possible Base Integration

### *To Consider*
- Full JSON request instead of only the *"event"* parameter [DONE]

### *Removed*
- Redirect Event Webhook and Handle the Endpoint (Reconsider?)


### *TODO*
- [ ] Excel File Handling - Remove duplicates
- [ ] Drop the connection after X failed attempts
- [ ] Enforce more protection from unknown requests
- [ ] Implement Testing Dashboard in the web app

### *To Fix*

### *WIP*
- [ ] Add specific user commands to which the bot responds (Starting with '!')

### *DONE*
- [X] Token Encryption
- [X] Encryption & Validation Webhook
- [X] Signature Validation
- [X] Split between Event and Verification Requests
- [X] EventStep Handling
- [X] Implement ~~Interface~~ map with signature function for mapping event_type with function to call
