package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunho/sdbx-discord-dj-bot/commands"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	str, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	config := djbot.Config{}
	err = yaml.Unmarshal(str, &config)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(file, os.Stdout)
	log.SetOutput(writer)

	dj, err := djbot.New(config)
	if err != nil {
		panic(err)
	}

	commands.Register(dj)

	err = dj.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("실행중")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dj.Close()
}
