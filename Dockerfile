FROM alpine:3.6
MAINTAINER Mihai Lupoiu <mihai.alexandru.lupoiu@gmail.com>

ENV TELEGRAM_CHAT_ID=@telegram-chat-id
ENV TELEGRAM_BOT_ID=@telegram-bot-id
ENV PACKT_URL=https://www.packtpub.com/packt/offers/free-learning


RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY PackBooksBotNotifier /bin/PackBooksBotNotifier
COPY cron /var/spool/cron/crontabs/root

RUN chmod +x /bin/PackBooksBotNotifier

CMD crond -l 2 -f

#For deploying only the image.
#FROM myhay/packt-free-learning:1