package poll

import "fmt"

// InvalidOptionError is returned when a vote is cast for an option that doesn't
// exist
type InvalidOptionError struct {
	msg    string
	option string
}

func (e *InvalidOptionError) Error() string {
	return fmt.Sprintf("%s\noption: %s\n", e.msg, e.option)
}

// AlreadyVotedError is returned when a user ties to vote multiple times for the
// same poll
type AlreadyVotedError struct {
	msg  string
	user string
}

func (e *AlreadyVotedError) Error() string {
	return fmt.Sprintf("%s\nuser: %s\n", e.msg, e.user)
}

// Poll contains the data relevant to a particular poll. The description, the voting options,
// and a list of votes.
type Poll struct {
	ID          int
	Description string
	Options     map[string]int
	Votes       []Vote
}

// Vote represents a users vote towards a particular option
type Vote struct {
	User   string
	Option string
}

var prevID = 0

// New creates and returns a new Poll type
func New(description string, options []string) Poll {
	optionMap := make(map[string]int)
	for _, option := range options {
		optionMap[option] = 0
	}
	return Poll{
		ID:          prevID + 1,
		Description: description,
		Options:     optionMap,
		Votes:       []Vote{},
	}
}

// Vote is a method on a Poll type that casts a vote towards a specified option.
// Returns an error if the specified option isn't available.
func (p *Poll) Vote(user string, option string) error {
	if _, ok := p.Options[option]; !ok {
		return &InvalidOptionError{"The given option isn't valid for the current poll", option}
	}

	for _, vote := range p.Votes {
		if vote.User == user {
			return &AlreadyVotedError{"The user already voted for the current poll", user}
		}
	}

	p.Options[option]++
	p.Votes = append(p.Votes, Vote{user, option})

	return nil
}

// GetResult returns a slice of the current options with the most votes
func (p *Poll) GetResult() []string {
	var maxVotes int
	var topOptions []string

	for option, voteCount := range p.Options {
		if voteCount > maxVotes {
			maxVotes = voteCount
			topOptions = []string{option}
		} else if voteCount == maxVotes {
			topOptions = append(topOptions, option)
		}
	}

	return topOptions
}
