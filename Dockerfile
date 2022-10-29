FROM golang:1.18-alpine

ENV GO111MODULE auto

ENV GOPATH /go/
ENV PORT 8080
ENV TELEGRAM_BOT_PATH /go/src/telegram-bot-service/

VOLUME ${TELEGRAM_BOT_PATH}
WORKDIR ${TELEGRAM_BOT_PATH}
COPY . ${TELEGRAM_BOT_PATH}

RUN chmod +x ./entrypoint.sh

RUN apk update && apk add curl \
                          git \
                          protobuf \
                          bash \
                          make

# Remember to change for env variable for heroku
EXPOSE 8080
CMD ["bash", "./entrypoint.sh"]
