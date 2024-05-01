## Job Panda: a telegram bot that scrapes webpages for job postings

I picked interest in Golang recently, and unlike other programming languages from the past<br>
this is different, because I will be  away from school for at least 5 months and I'd love to<br>
make international progress with programming. I have neither bought a course nor am I <br>
following any youtube full tutorial, I got started with a tour of Go and followed some <br>instructions on the official Go documentation.<br>

**Job panda** is my first project and I will love see it completed :).

## edited 1 may 2024.
The Bot employes a webscraping script to fetch job postings from jooble.com when /getme command<br>
is used.
The main entry of the program runs the bot's http handler<br>
```
package main

import "github.com/ekefan/go_job_scraper/handler"

func main() {
	
	handler.RunBotServer()

}
```
The function loads environment variables in the dotenv file, then listens of port:3000 to serve the HandlerFunction.<br>
- The handler picks up a http resquest called an update<br>
- Parses the update to convert it into a struct called Update; Update is designed to look like the<br>
json data in the update request.<br>
- Based on the update message, a http response is generated and sent.
- if the update message starts with the command /getme, the bot scrapes the web for the description<br>
following the command, like /getme backend developer jobs
- Other commands the bot recognizes are /help and /start and basic string responses are sent 

To test the functionality of the bot, I used ngok (a service that makes my localhost publicly available to the internet)<br>
I then curled the URL provided by the ngok to the bots URL so that Updates are received accordingly

## After all
Building the project has help me gain a clearer understanding of golang's syntax and a http server client

`#Building projects till I become proficient from basic to...`   as it gets complex I will know :)
