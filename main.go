package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var botID string

func main() {
	rand.Seed(time.Now().Unix())
	token := os.Getenv("DISCORD_BOT_AUTH_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	u, err := dg.User("@me")
	if err != nil {
		panic(err)
	}

	botID = u.ID

	dg.AddHandler(handleMessage)

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Bot %v is running\n", botID)

	<-make(chan struct{})
	return
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
		return
	}

	ok, err := regexp.MatchString("^!choose ", m.Content)
	if err != nil {
		panic(err)
	}
	if ok {
		handleChoose(s, m)
	}
}

func handleChoose(s *discordgo.Session, m *discordgo.MessageCreate) {
	options := strings.Split(m.Content[8:], ", ")

	option := options[rand.Intn(len(options))]

	fmt.Printf("%q chose: %s\n", options, option)

	msg := "How about " + option

	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		panic(err)
	}
}
