package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

	telegramBotID := os.Getenv("TELEGRAM_BOT_ID")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if telegramBotID == "" || telegramChatID == "" {
		log.Println("ERROR: Missing TELEGRAM_BOT_ID or TELEGRAM_CHAT_ID environment variables")
	}

	chatID, err := strconv.ParseInt(telegramChatID, 10, 64)
	if err != nil {
		log.Println("ERROR: TELEGRAM_CHAT_ID not a valid environment variable", err)
	}
	log.Println(chatID)

	bot, err := tgbotapi.NewBotAPI(telegramBotID)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	book := crawlURL()
	var text = fmt.Sprintf("Check out today's free ebook from Packt Publishing: 🎁 \n\n" +
		"📖 " + book.Title + "\n" +
		"🔎 " + book.Description + "\n" +
		"👉 " + url)

	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)

}
