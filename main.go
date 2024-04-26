package main

import "github.com/ekefan/go_job_scraper/handler"

// struct to describe the fields for the Job data structure
type Job struct {
	JobTitle, JobLink            string
	CompanyName, CompanyLocation string
}

// The entire programs starts running here
func main() {

	handler.RunBotServer()

}
