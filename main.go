package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mroseman95/discord-poll-bot/poll"
)

var botID string

// pollRegex captures commands of the form
// "!poll some description text [option1, option2, option3]"
// and finds group1 = (some description text) and
// group2 = (option1, option2, option3)
var pollRegex = regexp.MustCompile(`^!poll ([^\[]*)\[(.*)\]$`)

// voteRegex captures commands of the form
// "!vote some option"
// and finds group1 = (some option)
var voteRegex = regexp.MustCompile(`^!vote (.*)$`)

var currentPoll poll.Poll

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

	fmt.Printf("Bot %v is running", botID)

	<-make(chan struct{})
	return
}

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
		return
	}

	if m.Content[:5] == "!poll" {
		matches := pollRegex.FindStringSubmatch(m.Content)
		if matches == nil {
			_, err := s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't understand that command.")
			if err != nil {
				panic(err)
			}
			return
		}
		description := matches[1]
		options := matches[2]

		currentPoll = poll.New(description, strings.Split(options, ", "))

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("New poll started for: %s\nwith options: %q", description, options))
		if err != nil {
			panic(err)
		}
	}

	if m.Content[:5] == "!vote" {
		// TODO check that user hasn't voted already
		// TODO check that there is an active poll
		matches := voteRegex.FindStringSubmatch(m.Content)
		if matches == nil {
			_, err := s.ChannelMessageSend(m.ChannelID, "Sorry I couldn't understand that command.")
			if err != nil {
				panic(err)
			}
			return
		}
		option := matches[1]

		err := currentPoll.Vote(m.Author.ID, option)
		if err != nil {
			// TODO check this error type
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sorry, %s, is not a valid option in the current poll", option))
			if err != nil {
				panic(err)
			}
			return
		}

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has voted for %s", m.Author.Username, option))
		if err != nil {
			panic(err)
		}
	}
}

// poll := poll.New("test poll", []string{"yes", "no"})
// fmt.Printf("%+v\n", poll)
//
// err := poll.Vote("matt", "no")
// if err != nil {
// 	panic(err)
// }
// err = poll.Vote("tony", "yes")
// if err != nil {
// 	panic(err)
// }
// err = poll.Vote("aidan", "yes")
// if err != nil {
// 	panic(err)
// }
// fmt.Printf("%+v\n", poll)
// fmt.Printf("%q\n", poll.GetResult())
