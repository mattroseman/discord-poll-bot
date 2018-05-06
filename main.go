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

// newPollRegex captures commands of the form
// "!newpoll some description text [option1, option2, option3]"
// and finds group1 = (some description text) and
// group2 = (option1, option2, option3)
var newPollRegex = regexp.MustCompile(`!newpoll ([^\[]*)\[(.*)\]$`)

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
		handleShowPoll(s, m)
	}

	if m.Content[:8] == "!newpoll" {
		handleNewPoll(s, m)
	}

	if m.Content[:5] == "!vote" {
		handleNewVote(s, m)
	}

	if m.Content[:7] == "!result" {
		handleGetResults(s, m)
	}
}

func handleShowPoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := fmt.Sprintf("Current Poll: %s", currentPoll.Description)

	for option, numVotes := range currentPoll.Options {
		msg = msg + fmt.Sprintf("\nOption: %s | Votes: %d", option, numVotes)
	}

	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	if err != nil {
		panic(err)
	}
}

func handleNewPoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	matches := newPollRegex.FindStringSubmatch(m.Content)
	if matches == nil {
		_, err := s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("Sorry %s, I couldn't understand that command.", m.Author.Username))
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

func handleNewVote(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO check that there is an active poll
	matches := voteRegex.FindStringSubmatch(m.Content)
	if matches == nil {
		_, err := s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("Sorry %s, I couldn't understand that command.", m.Author.Username))
		if err != nil {
			panic(err)
		}
		return
	}
	option := matches[1]

	err := currentPoll.Vote(m.Author.ID, option)
	if err != nil {
		if _, ok := err.(*poll.AlreadyVotedError); ok {
			_, err := s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("Sorry %s, you have already voted in this poll.", m.Author.Username))
			if err != nil {
				panic(err)
			}
			return
		} else if _, ok := err.(*poll.InvalidOptionError); ok {
			_, err := s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("Sorry %s, \"%s\" is an invalid option for this poll.", m.Author.Username, option))
			if err != nil {
				panic(err)
			}
			return
		} else {
			panic(err)
		}
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has voted for %s", m.Author.Username, option))
	if err != nil {
		panic(err)
	}
}

func handleGetResults(s *discordgo.Session, m *discordgo.MessageCreate) {
	winners := currentPoll.GetResult()
	if len(winners) == 1 {
		_, err := s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("The votes are in. \"%s\" has received the most votes", winners[0]))
		if err != nil {
			panic(err)
		}
		return
	}

	_, err := s.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf("The votes are in. There is a tie between %q", winners))
	if err != nil {
		panic(err)
	}
}
