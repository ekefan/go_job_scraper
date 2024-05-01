## Job Panda: a telegram bot that scrapes webpages for job postings

I picked interest in Golang recently, and unlike other programming languages from the past<br>
this is different, because I will be  away from school for at least 5 months and I'd love to<br>
make international progress with programming. I have neither bought a course nor following any<br>
youtube full tutorial, I got started with a tour of Go and followed some instructions on the <br>
official Go documentation.<br>

`Job panda` is my first project and I will love see it completed :).

## edited 1 may 2024.
The Bot employes a webscraping script to fetch job postings from jooble.com when /getme command<br>
is used.
The main entry of the program runs the bots http handler<br>
```
package main

import "github.com/ekefan/go_job_scraper/handler"

func main() {
	
	handler.RunBotServer()

}
```
The function loads environment variables in the dotenv file, then listens of port:3000 to server HandlerFunction.<br>
- The handler picks up a http resquest called an update<br>
- Parses the update to convert it into a struct called update; update is designed to look like the<br>
json data in the update request.<br>
- Based on the update message, a http response is sent.
- if the update message start with the command /getme, the bot scrapes the web for the description<br>Following the command
- Other commands the bot recognizes are /help and /start
