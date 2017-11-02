package main

import (
	"log"
	"strings"

	"github.com/anaskhan96/soup"

	"gopkg.in/telegram-bot-api.v4"
)

const (
	url = "https://www.packtpub.com/packt/offers/free-learning"
)

func crawlURL() {

	resp, err := soup.Get(url)

	if err != nil {
		log.Println("ERROR: Failed to crawl \"" + url + "\"")
		return
	}

	doc := soup.HTMLParse(resp)

	title := doc.Find("div", "class", "dotd-title").Find("h2")
	t := strings.TrimSpace(title.Text())
	log.Printf("%d %q", len(t), t)

	image := doc.Find("div", "class", "dotd-main-book-image float-left").Find("noscript")
	i := strings.TrimSpace(image.Text())
	i_split := strings.Split(i, "\"")
	log.Printf("%d %q", len(i_split[1]), i_split[1])

	description := doc.Find("div", "class", "dotd-main-book-summary float-left").Find("div").FindNextElementSibling().FindNextElementSibling().FindNextElementSibling()
	d := strings.TrimSpace(description.Text())
	log.Printf("%d %q", len(d), d)
}

func main() {
	bot, err := tgbotapi.NewBotAPI("asdf")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		log.Println("-------------------------------")
		crawlURL()
		log.Println("-------------------------------")
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
