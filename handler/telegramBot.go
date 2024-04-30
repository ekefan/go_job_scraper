package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ekefan/go_job_scraper/scraper"
)

const (
	startCommand  string = "/start"
	getJobCommand string = "/getme"
	helpCommand   string = "/help"

	telegramApiBaseUrl     string = "https://api.telegram.org/bot"
	telegramApiSendMessage string = "/sendMessage"
	telegramToken          string = "TELEGRAM_BOT_TOKEN"
)

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramToken) + telegramApiSendMessage

// Update is a Telegram object that the handler receives every time
// a user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int64 `json:"id"`
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	fmt.Printf("From line 46 checking what an update looks link\n %v", update)
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
	updateCommand := strings.Fields(update.Message.Text)[0]
	updateText := strings.Fields(update.Message.Text)[1:]

	switch updateCommand {
		case getJobCommand:
			jobsDescription := updateText
			sendToTelegram(update.Message.Chat.Id, "Getting Jobs...\nHold on")
			jobs := scraper.GetJobs(jobsDescription)
			fmt.Printf("%T", jobs)
			errSendingJobs := sendJobsToTelegramChat(update.Message.Chat.Id, jobs)
			if errSendingJobs != nil {
				log.Printf("Got error sending jobs to telegram:%v",
					errSendingJobs.Error())
			}
		case startCommand:
			startText := "Welcome to job panda"
			sendToTelegram(update.Message.Chat.Id, startText)
		case helpCommand:
			helpText := "The available commands are"
			sendToTelegram(update.Message.Chat.Id, helpText)
		default:
			startText := "Welcome to job panda"
			sendToTelegram(update.Message.Chat.Id, startText)
	}
}
func sendToTelegram(chatId int64, message string) error{
	_, err := sendTextToTelegramChat(chatId, message)
	if err != nil {
		log.Printf("got error %s sending Text to telegram", err.Error())
		return err
	}
	return  nil
}
type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func sendJobsToTelegramChat(chatId int64, jobs []scraper.Job) error {
	// sendTextToTelegramChat(chatId, text)
	for _, job := range jobs {
		text := job.GetJobResponseText()
		err := sendToTelegram(chatId, text)
		if err != nil {
			return err
		}
	}
	return nil
}
func sendTextToTelegramChat(chatId int64, text string) (string, error) {
	// Create the request body struct
	reqBody := &sendMessageReqBody{
		ChatID: chatId,
		Text:   text,
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
		fmt.Println(res.StatusCode)
		return "", errors.New("unexpected status" + res.Status)
	}

	return reqBody.Text, nil
}

func RunBotServer() {
	http.ListenAndServe(":3000", http.HandlerFunc(HandleTelegramWebHookTest))
}

