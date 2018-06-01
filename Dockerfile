FROM ksunhokim/djbot-base

WORKDIR /go/src/github.com/sunho/sdbx-discord-dj-bot
COPY . .

RUN go get -d -v ./...
RUN go get github.com/bwmarrin/dca/cmd/dca
RUN go install

CMD ["sdbx-discord-dj-bot"]