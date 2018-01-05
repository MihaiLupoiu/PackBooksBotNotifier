package crawlpackt

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

var errorRetrevingBook string
var url string

type bookInfoMessage struct {
	Title            string
	Description      string
	ImageURL         string
	TimeLeft         string
	OtherInformation string
}

// InitCrawlPackt set URL to crawl book.
func Init(recvURL string) {
	if recvURL != "" {
		url = recvURL
	} else {
		url = "https://www.packtpub.com/packt/offers/free-learning"
		log.Printf("Warning: No Packt URL provided. Using default one:" + url)
	}
}

// GetText will return the string with the book information
func GetText() string {
	bookInformation := crawlURL()
	return buildString(bookInformation)
}

func buildString(book bookInfoMessage) string {

	text := fmt.Sprintf("\n⚠️ An error ocurred in retreving today's free ebook from Packt Publishing ⚠️ \n\n" +
		"Meanwhile we fix the issue, please check Packt Publishing web page. 👉 " + url)

	if errorRetrevingBook == "" {
		text = fmt.Sprintf("\nCheck out today's free ebook from Packt Publishing 🎁 \n\n" +
			"📖 " + book.Title + "\n" +
			"🔎 " + book.Description + "\n" +
			"⌛️ " + book.TimeLeft + "\n" +
			"👉 " + url)
		return text
	}

	log.Printf("%s", errorRetrevingBook)
	return text
}

func crawlURL() bookInfoMessage {
	var book bookInfoMessage
	resp, err := soup.Get(url)

	if err != nil {
		log.Println("ERROR: Failed to crawl \"" + url + "\"")
		errorRetrevingBook = "Error executing crawl."
		return book
	}

	doc := soup.HTMLParse(resp)

	title := doc.Find("div", "class", "dotd-title").Find("h2")
	if title.Text() == "" {
		errorRetrevingBook = "Error retreving book title."
		return book
	}

	book.Title = strings.TrimSpace(title.Text())
	log.Printf("%d %q\n", len(book.Title), book.Title)

	image := doc.Find("div", "class", "dotd-main-book-image float-left").Find("noscript")
	if image.Text() == "" {
		errorRetrevingBook = "Error retreving book image."
		return book
	}
	i := strings.TrimSpace(image.Text())
	iSplit := strings.Split(i, "\"")
	book.ImageURL = strings.Trim(iSplit[1], "//")
	log.Printf("%d %q\n", len(book.ImageURL), book.ImageURL)

	description := doc.Find("div", "class", "dotd-main-book-summary float-left").Find("div").FindNextElementSibling().FindNextElementSibling().FindNextElementSibling()
	book.Description = strings.TrimSpace(description.Text())
	if description.Text() == "" {
		errorRetrevingBook = "Error retreving book description."
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
