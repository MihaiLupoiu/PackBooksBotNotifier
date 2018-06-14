# PackBooksBotNotifier
Notifier about the free book of the day offered by Pack. 

[![Build Status](https://travis-ci.org/MihaiLupoiu/PackBooksBotNotifier.svg?branch=master)](https://travis-ci.org/MihaiLupoiu/PackBooksBotNotifier)


Build:

CGO_ENABLED=0 go build -a -installsuffix cgo

Set variables:

export TELEGRAM_BOT_ID=telegram id
export TELEGRAM_CHAT_ID=chat id or channel id
