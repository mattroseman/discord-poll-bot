package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var botID string

func main() {
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

	if len(m.Content) > 6 && m.Content[:7] == "!choose" {
		handleChoose(s, m)
	}
}

func handleChoose(s *discordgo.Session, m *discordgo.MessageCreate) {
	options := strings.Split(m.Content[8:], ", ")
	fmt.Printf("%q\n", options)
	fmt.Printf("%d mod %d\n", rand.Int(), len(options))

	option := options[int(math.Mod(float64(rand.Int()), float64(len(options))))]

	msg := "How about " + option

	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		panic(err)
	}
}
