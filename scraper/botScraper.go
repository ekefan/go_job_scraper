package scraper

import (
	"fmt"
	"log"
	"strings"
	"github.com/gocolly/colly/v2"
)

// struct to describe the fields for the Job data structure
type Job struct {
	JobTitle, JobLink string
	CompanyName, CompanyLocation string
}



func GetJobs(searchWords []string){
	//get the url for the jobsearch scrapping
	usrDescription := strings.Join(searchWords, "+")
	url := "https://ng.jooble.org/SearchResult?ukw="
	webURL := fmt.Sprintf("%s%s", url, usrDescription)

	c:= colly.NewCollector()

	jobPostings := []Job{}
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error){
		log.Println("Ops this operation was not successful", err)
	})
	c.OnHTML("div.ojoFrF", func(e *colly.HTMLElement){
		jobPosting := Job{}

		jobPosting.JobLink = e.ChildAttr("a.hyperlink_appearance_undefined", "href")
		jobPosting.JobTitle = e.ChildText("a.hyperlink_appearance_undefined")
		jobPosting.CompanyName = e.ChildText("div.VeoRvG")
		jobPosting.CompanyLocation = e.ChildText("div.nxYYVJ")
		
		jobPostings = append(jobPostings, jobPosting)
	})
	c.OnResponse(func(r *colly.Response){
		fmt.Println("Visited", r.Request.URL)
	})

	err := c.Visit(webURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(jobPostings))

}