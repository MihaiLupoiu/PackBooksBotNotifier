package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anaskhan96/soup"

	"gopkg.in/telegram-bot-api.v4"
)

type bookInfoMessage struct {
	Title       string
	Description string
	ImageURL    string
}

const (
	url = "https://www.packtpub.com/packt/offers/free-learning"
)

func crawlURL() bookInfoMessage {
	var book bookInfoMessage
	resp, err := soup.Get(url)

	if err != nil {
		log.Println("ERROR: Failed to crawl \"" + url + "\"")
		return book
	}

	doc := soup.HTMLParse(resp)

	title := doc.Find("div", "class", "dotd-title").Find("h2")
	book.Title = strings.TrimSpace(title.Text())
	log.Printf("%d %q", len(book.Title), book.Title)

	image := doc.Find("div", "class", "dotd-main-book-image float-left").Find("noscript")
	i := strings.TrimSpace(image.Text())
	iSplit := strings.Split(i, "\"")
	book.ImageURL = strings.Trim(iSplit[1], "//")
	log.Printf("%d %q", len(book.ImageURL), book.ImageURL)

	description := doc.Find("div", "class", "dotd-main-book-summary float-left").Find("div").FindNextElementSibling().FindNextElementSibling().FindNextElementSibling()
	book.Description = strings.TrimSpace(description.Text())
	log.Printf("%d %q", len(book.Description), book.Description)
	return book
}

func checkBook(chatID int64) {

}

func main() {
	bot, err := tgbotapi.NewBotAPI("asdf")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 //24*60* 60
	//6234638

	var chatID int64

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID = update.Message.Chat.ID

		ticker := time.NewTicker(time.Millisecond * 60 * 1000)

		go func() {
			for t := range ticker.C {
				fmt.Println("Tick at", t)
				log.Println("-------------------------------")
				book := crawlURL()
				log.Println("-------------------------------")
				//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				var text = fmt.Sprintf("Check out today's free ebook from Packt Publishing: \n\n" +
					"📖 " + book.Title + "\n" +
					"🔎 " + book.Description + "\n" +
					"➡️ " + url)

				msg := tgbotapi.NewMessage(chatID, text)
				bot.Send(msg)
			}
		}()

	}
}
