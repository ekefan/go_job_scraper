package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ekefan/go_job_scraper/scraper"
	"github.com/joho/godotenv"
)

const (
	startCommand  string = "/start"
	getJobCommand string = "/getme"
	helpCommand   string = "/help"

	telegramApiBaseUrl     string = "https://api.telegram.org/bot"
	telegramApiSendMessage string = "/sendMessage"
	telegramToken          string = "TELEGRAM_BOT_TOKEN"
)

// telegramApi uses consts from telegramBotUtils.go
var telegramApi string
var startText string = fmt.Sprintf(
	"Welcome to Job Panda\n" +
		"This bot brings to your dm the latest Job postings based on your description from jooble.com" +
		"Using the command /getme <your preferred job>\n" +
		"/help shows you all commands and usage\n",
)

var helpText string = fmt.Sprintf(
	"Briefs on Job Panda\n" +
		"Job Panda has only three commands: /start, /getme and /help\n" +
		"This message comes up everytime you send a cmd or message the bot does not recognize\n" +
		"The /start command by default runs at the at the start of conversation with the bot\n" +
		"/getme command is used to search for your desired jobs from the jobble.com. Usage: /getme backend developer",
)

// Update is a Telegram object that the handler receives every time
// a user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int      `json:"message_id"`
	From      From     `json:"from"`
	Text      string   `json:"text"`
	Chat      Chat     `json:"chat"`
	Date      int      `json:"date"`
	Entities  []Entity `json:"entities,omitempty"`
}

type From struct {
	Id           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type Entity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func LoadDotEnv(filename string) error {
	err := godotenv.Load(filename)

	if err != nil {
		fmt.Printf("Error loading .env files")
		return err
	}
	return nil
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("could not read request body: %s", err.Error())
		return nil, err
	}
	defer r.Body.Close()
	// log.Printf("Raw request body: %s", string(body))
	var update Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("could not decode incoming update: %s", err.Error())
		return nil, err
	}
	return &update, nil
}

// sendToTelegram: DRY implementation for sending text to telegram
func sendToTelegram(chatId int64, message string) error {
	_, err := sendTextToTelegramChat(chatId, message)
	if err != nil {
		log.Printf("got error %s sending Text to telegram", err.Error())
		return err
	}
	return nil
}

// sendJobsToTelegramChat: generates texts from job and sends message to chat_id
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

// sendTextToTelegramChat handles the functionality of
// sending response messages to the respective chatId(dm)
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

	//make a post request to the
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

// HandleTelegramWebHookTest: main http handler for the telegram bot
func HandleTelegramWebHookTest(w http.ResponseWriter, r *http.Request) {
	update, err := parseTelegramRequest(r)
	if err != nil {
		log.Printf("Error parsing update, %s", err.Error())
		http.Error(w, "could not parse update", http.StatusInternalServerError)
		return
	}
	log.Printf("Received update: %+v", update)

	if update.Message.Text == "" {
		fmt.Println(*update)
	}
	updateCommand := strings.Fields(update.Message.Text)[0]
	updateText := strings.Fields(update.Message.Text)[1:]
	switch updateCommand {
	case getJobCommand:
		jobsDescription := updateText
		sendToTelegram(update.Message.Chat.Id, "Getting Jobs...\nHold on")
		jobs := scraper.GetJobs(jobsDescription)
		if len(jobs) != 0 {
			errSendingJobs := sendJobsToTelegramChat(update.Message.Chat.Id, jobs)
			if errSendingJobs != nil {
				log.Printf("Got error sending jobs to telegram: %v", errSendingJobs.Error())
			}
		} else {
			sendToTelegram(update.Message.Chat.Id, "Sorry no Jobs matched your description\n Maybe try refining the description")
		}
	case startCommand:
		sendToTelegram(update.Message.Chat.Id, startText)
	case helpCommand:
		sendToTelegram(update.Message.Chat.Id, helpText)
	default:
		sendToTelegram(update.Message.Chat.Id, helpText)
	}
}

func RunBotServer() {
	telegramApi = fmt.Sprintf("%s%v%s", telegramApiBaseUrl, os.Getenv("TELEGRAM_BOT_TOKEN"), telegramApiSendMessage)
	fmt.Println("listening")
	http.ListenAndServe(":8080", http.HandlerFunc(HandleTelegramWebHookTest))

}
