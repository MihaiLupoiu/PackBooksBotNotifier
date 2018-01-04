package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"

	"gopkg.in/telegram-bot-api.v4"
)

var errorRetrevingBoot string

type bookInfoMessage struct {
	Title       string
	Description string
	ImageURL    string
	TimeLeft    string
}

const (
	url = "https://www.packtpub.com/packt/offers/free-learning"
)

func crawlURL() bookInfoMessage {
	var book bookInfoMessage
	resp, err := soup.Get(url)

	if err != nil {
		log.Println("ERROR: Failed to crawl \"" + url + "\"")
		errorRetrevingBoot = "Error executing crawl."
		return book
	}

	doc := soup.HTMLParse(resp)

	title := doc.Find("div", "class", "dotd-title").Find("h2")
	if title.Text() == "" {
		errorRetrevingBoot = "Error retreving book title."
		return book
	}

	book.Title = strings.TrimSpace(title.Text())
	log.Printf("%d %q\n", len(book.Title), book.Title)

	image := doc.Find("div", "class", "dotd-main-book-image float-left").Find("noscript")
	if image.Text() == "" {
		errorRetrevingBoot = "Error retreving book image."
		return book
	}
	i := strings.TrimSpace(image.Text())
	iSplit := strings.Split(i, "\"")
	book.ImageURL = strings.Trim(iSplit[1], "//")
	log.Printf("%d %q\n", len(book.ImageURL), book.ImageURL)

	description := doc.Find("div", "class", "dotd-main-book-summary float-left").Find("div").FindNextElementSibling().FindNextElementSibling().FindNextElementSibling()
	book.Description = strings.TrimSpace(description.Text())
	if description.Text() == "" {
		errorRetrevingBoot = "Error retreving book description."
		return book
	}
	log.Printf("%d %q\n", len(book.Description), book.Description)

	timeNow := time.Now().Unix()
	timeLimitMap := doc.Find("span", "class", "packt-js-countdown")

	timeLimit, err := strconv.ParseInt(timeLimitMap.Attrs()["data-countdown-to"], 10, 64)
	if err != nil {
		log.Println("ERROR: Failed to convert timeLimit.", err)
		return book
	}
	timeLeft := time.Unix(timeLimit-timeNow, 0)
	book.TimeLeft = fmt.Sprintf("%02d:%02d:%02d", timeLeft.Hour(), timeLeft.Minute(), timeLeft.Second())

	return book
}

func checkBook(chatID int64) {

}

func main() {

	telegramBotID := os.Getenv("TELEGRAM_BOT_ID")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	errorRetrevingBoot = ""

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
	var text string
	book := crawlURL()

	if errorRetrevingBoot == "" {
		text = fmt.Sprintf("\nCheck out today's free ebook from Packt Publishing 🎁 \n\n" +
			"📖 " + book.Title + "\n" +
			"🔎 " + book.Description + "\n" +
			"⌛️ " + book.TimeLeft + "\n" +
			"👉 " + url)
	} else {
		log.Printf("%s", errorRetrevingBoot)
		text = fmt.Sprintf("\n⚠️ An error ocurred in retreving today's free ebook from Packt Publishing ⚠️ \n\n" +
			"For more details, please check Packt Publishing web page. 👉 " + url)
	}

	msg := tgbotapi.NewMessage(chatID, text)
	msg.DisableWebPagePreview = true
	bot.Send(msg)

}
