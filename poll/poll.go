package poll

import "fmt"

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
		return fmt.Errorf("The given option isn't a valid option for this poll: %s", option)
	}

	for _, vote := range p.Votes {
		if vote.User == user {
			return fmt.Errorf("User %s has already voted in this poll", user)
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
