package poll

import (
	"errors"
	"fmt"
)

// Poll contains information relevent to a specific poll
type Poll struct {
	Options []string
	Votes   []Vote
}

// Vote represents a vote by one person towards one option
type Vote struct {
	Option string
	Voter  string
}

// NewPoll creates a new Poll type with the given options and returns a pointer to it.
func NewPoll(options []string) (*Poll, error) {
	if len(options) < 2 {
		return nil, errors.New("must supply at least two options")
	}

	return &Poll{options, []Vote{}}, nil
}

func (p Poll) equal(q Poll) bool {
	if fmt.Sprintf("%q", p.Options) == fmt.Sprintf("%q", q.Options) {
		if fmt.Sprintf("%q", p.Votes) == fmt.Sprintf("%q", q.Votes) {
			return true
		}
	}

	return false
}

// Vote casts a vote towards one of the options in the given Poll
func (p *Poll) Vote(option, voter string) error {
	for _, v := range p.Votes {
		if v.Voter == voter {
			return errors.New("this voter already voted on this poll")
		}
	}

	for _, o := range p.Options {
		if o == option {
			p.Votes = append(p.Votes, Vote{option, voter})
			return nil
		}
	}

	return errors.New("unkown option for this poll")
}

// GetResult returns a slice of the Poll options with the most votes
func (p Poll) GetResult() []string {
	mostVotes := 0
	winningOptions := []string{}

	for _, o := range p.Options {
		numVotes := 0
		for _, v := range p.Votes {
			if v.Option == o {
				numVotes++
			}
		}

		if numVotes > mostVotes {
			winningOptions = []string{o}
			mostVotes = numVotes
		} else if numVotes == mostVotes {
			winningOptions = append(winningOptions, o)
		}
	}

	return winningOptions
}
