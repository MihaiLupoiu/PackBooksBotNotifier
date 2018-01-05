package main

import (
	"log"
	"os"

	"github.com/MihaiLupoiu/PackBooksBotNotifier/crawlpackt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	url := os.Getenv("PACKTURL")
	crawlpackt.Init(url)

	telegramBotID := os.Getenv("TELEGRAM_BOT_ID")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if telegramBotID == "" || telegramChatID == "" {
		log.Println("ERROR: Missing TELEGRAM_BOT_ID or TELEGRAM_CHAT_ID environment variables")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(telegramBotID)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	text := crawlpackt.GetText()

	msg := tgbotapi.NewMessageToChannel(telegramChatID, text)
	msg.DisableWebPagePreview = true
	bot.Send(msg)
	// A free eBook every day notification. | Not affiliated with Packt Publishing | maintained by @myhay
	// Packt Free Learning - Free Programming Ebooks
}
