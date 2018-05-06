package poll

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNewPoll(t *testing.T) {
	tests := []struct {
		options  []string
		ok       bool
		err      error
		expected *Poll
	}{
		{
			[]string{"yes", "no"},
			true,
			nil,
			&Poll{[]string{"yes", "no"}, make(map[string][]Vote)},
		},
		{
			[]string{},
			false,
			errors.New("must supply at least two options"),
			nil,
		},
		{
			[]string{"yes"},
			false,
			errors.New("must supply at least two options"),
			nil,
		},
	}

	for _, test := range tests {
		got, err := NewPoll(test.options)

		if err != nil {
			if test.ok {
				t.Errorf("NewPoll returned unexpected error: %v", err)
			}
		} else {
			if !test.ok {
				t.Errorf("NewPoll didn't return an expected error: %v", test.err)
			}

			if !reflect.DeepEqual(*test.expected, *got) {
				t.Errorf("NewPoll returned incorect Poll.\nGot: %+v\nWant: %+v",
					got, test.expected)
			}
		}
	}
}

func TestVote(t *testing.T) {
	tests := []struct {
		poll     *Poll
		option   string
		voter    string
		ok       bool
		err      error
		expected *Poll
	}{
		{
			&Poll{[]string{"yes", "no"}, make(map[string][]Vote)},
			"yes",
			"testuser",
			true,
			nil,
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser"}},
				},
			},
		},
		{
			&Poll{[]string{"yes", "no"}, make(map[string][]Vote)},
			"no",
			"testuser",
			true,
			nil,
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"no": []Vote{Vote{"no", "testuser"}},
				},
			},
		},
		{
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser1"}},
				},
			},
			"yes",
			"testuser2",
			true,
			nil,
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser1"}, Vote{"yes", "testuser2"}},
				},
			},
		},
		{
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser1"}, Vote{"yes", "testuser4"}},
					"no":  []Vote{Vote{"no", "testuser2"}, Vote{"no", "testuser3"}},
				},
			},
			"yes",
			"testuser5",
			true,
			nil,
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{
						Vote{"yes", "testuser1"},
						Vote{"yes", "testuser4"},
						Vote{"yes", "testuser5"},
					},
					"no": []Vote{Vote{"no", "testuser2"}, Vote{"no", "testuser3"}},
				},
			},
		},
		{
			&Poll{[]string{"yes", "no"}, make(map[string][]Vote)},
			"y",
			"testuser",
			false,
			errors.New("unknown option for this poll"),
			&Poll{[]string{"yes", "no"}, make(map[string][]Vote)},
		},
		{
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser"}},
				},
			},
			"no",
			"testuser",
			false,
			errors.New("this voter already voted on this poll"),
			&Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser"}},
				},
			},
		},
	}

	for _, test := range tests {
		err := test.poll.Vote(test.option, test.voter)

		if err != nil {
			if test.ok {
				t.Errorf("Vote returned unexpected error: %v", err)
			}
		} else {
			if !test.ok {
				t.Errorf("Vote didn't return expected error: %v", test.err)
			}

			if !reflect.DeepEqual(*test.expected, *test.poll) {
				t.Errorf("Vote didn't update poll correctly.\nGot: %+v\nWant: %+v",
					test.poll, test.expected)
			}
		}
	}
}

func TestGetResult(t *testing.T) {
	tests := []struct {
		poll     Poll
		expected []string
	}{
		{
			Poll{
				[]string{"yes", "no"},
				map[string][]Vote{
					"yes": []Vote{Vote{"yes", "testuser"}},
				},
			},
			[]string{"yes"},
		},
		{
			Poll{[]string{"yes", "no"}, make(map[string][]Vote)},
			[]string{},
		},
	}

	for _, test := range tests {
		result := test.poll.GetResult()

		if fmt.Sprintf("%q", result) != fmt.Sprintf("%q", test.expected) {
			t.Errorf("GetResult didn't return correct result\nGot: %s\nWant: %s\n",
				result, test.expected)
		}
	}
}
