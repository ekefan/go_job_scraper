package handler

import (
	"fmt"
	"os"
	"log"
	"net/http"
	"encoding/json"
	"errors"
	"bytes"
)

const startCommand string = "/start"
const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramToken string = "TELEGRAM_BOT_TOKEN"

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramToken) + telegramApiSendMessage
// Update is a Telegram object that the handler receives every time 
// a user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int64 `json:"id"`
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func HandleTelegramWebHookTest(w http.ResponseWriter, r *http.Request) {
	update, err := parseTelegramRequest(r) 
	if err != nil {
		log.Printf("Error parsing update, %s", err.Error())
		return
	}
	fmt.Println("received update")
	updateMessage := update.Message
	if updateMessage.Text == startCommand {
		telegramResponseBody, errTelegram := sendTextToTelegramChat(updateMessage.Chat.Id)
		if errTelegram != nil {
			log.Printf("Got error %s sending message to telegram:%s", errTelegram.Error(), telegramResponseBody)
			} else {
				log.Printf("successfully distributed to chat id %d", update.Message.Chat.Id)
			}
	}
}

type sendMessageReqBody struct {
	ChatID int64 `json:"chat_id"`
	Text string  `json:"text"`
}
func sendTextToTelegramChat(chatId int64) (string, error) {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: chatId,
		Text: "Welcome Bot Panda",
	}
	//Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	res, err := http.Post(telegramApi, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("unexpected status" + res.Status)
	}

	return reqBody.Text, nil
}

func RunBotServer() {
	http.ListenAndServe(":3000", http.HandlerFunc(HandleTelegramWebHookTest))
}