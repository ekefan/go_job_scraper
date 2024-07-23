package scraper

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly/v2"
)

// struct to describe the fields for the Job data structure
type Job struct {
	JobTitle, JobLink string
	CompanyName, CompanyLocation string
}


func (j *Job) GetJobResponseText() string {
	jobText := fmt.Sprintf("Job Title: %v\nCompany:%v\nLocated at:%v\nKnow more:%v", 
	j.JobTitle, 
	j.CompanyName, 
	j.CompanyLocation,
	j.JobLink)
	return jobText
}
func GetJobs(searchWords []string) []Job {
	//get the url for the jobsearch scrapping
	usrDescription := strings.Join(searchWords, "%20")
	url := "https://ng.jooble.org/SearchResult?ukw="
	webURL := fmt.Sprintf("%s%s", url, usrDescription)
	c:= colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	jobPostings := []Job{}
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		fmt.Println(r.Headers)
	})
	c.OnError(func(resp *colly.Response, err error){
		fmt.Println(resp.Headers)
		fmt.Println("Ops this operation was not successful", err)
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
		fmt.Println(err)
	}
	return jobPostings //slice containing all job postings from this scraping session

}