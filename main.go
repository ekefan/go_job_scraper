package main

import (
	"log"
	"github.com/ekefan/go_job_scraper/handler"
)

// The entire programs starts running here
func main() {
	errLoadingEnv := handler.LoadDotEnv("job.env")
	if errLoadingEnv != nil {
		log.Printf("Error loading the env variables")
		return
	}
	handler.RunBotServer()
}
