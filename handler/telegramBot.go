package handler

import (
	"fmt"
	"os"
)

func Printenv() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		fmt.Println("TELEGRAM_BOT_TOKEN is not set or empty")
		return
	}
	fmt.Printf("Telegram API Token: %T\n%v", token, token)
}
