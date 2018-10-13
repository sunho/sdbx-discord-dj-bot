FROM golang:stretch

RUN apt-get update && apt-get install -y \ 
            ffmpeg \
            curl
RUN curl -o /usr/bin/youtube-dl http://yt-dl.org/latest/youtube-dl
RUN chmod 755 /usr/bin/youtube-dl

WORKDIR /go/src/github.com/sunho/sdbx-discord-dj-bot
COPY . .

RUN go get -d -v ./...
RUN go get github.com/bwmarrin/dca/cmd/dca
RUN go install

CMD ["sdbx-discord-dj-bot"]
