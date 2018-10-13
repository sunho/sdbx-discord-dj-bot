FROM golang:stretch

RUN apt-get update && apt-get install -y \ 
            ffmpeg \
            curl
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

WORKDIR /go/src/github.com/sunho/sdbx-discord-dj-bot
COPY . .

RUN go get -d -v ./...
RUN go get github.com/bwmarrin/dca/cmd/dca
RUN go install

CMD ["sdbx-discord-dj-bot"]
