package main

import (
	"fmt"
	"log"
	"github.com/gocolly/colly/v2"
)

// a type to describe the fields for the Job type
type Job struct {
	jobTitle, jobDetails string
}


//Callback function to handle OnRequests callbacks to the
//collector

//The entire programs starts running here
func main() {
	c:= colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error){
		log.Println("Something went wrong", err)
	})
	// c.OnHTML("")
	c.OnResponse(func(r *colly.Response){
		fmt.Println("Visited", r.Request.URL)
	})

	err := c.Visit("https://ng.jooble.org/SearchResult?ukw=backend+developer")
	if err != nil {
		log.Fatal(err)
	}
}