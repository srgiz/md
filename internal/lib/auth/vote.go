package auth

import (
	"context"
	"errors"
)

const ContextKeyAuthToken = "auth.Token"
const ContextKeyAuthAttr = "auth.Attr"

var ErrAuthNoVoters = errors.New("auth: no voters")
var ErrAuthTokenNotFound = errors.New("auth: token not found")
var ErrAuthAttrNotFound = errors.New("auth: attr not found")

var DefaultManager = NewManager()

type Voter interface {
	Supports(attr string, subject any) bool
	Vote(token Token, attr string, subject any) error
}

type Token interface {
	Data() map[string]any
}

type Manager struct {
	voters []Voter
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) AddVoters(voters ...Voter) {
	m.voters = append(m.voters, voters...)
}

func (m *Manager) DecideContext(ctx context.Context, subject any) error {
	token, hasToken := ctx.Value(ContextKeyAuthToken).(Token)

	if !hasToken {
		return ErrAuthTokenNotFound
	}

	attr, hasAttr := ctx.Value(ContextKeyAuthAttr).(string)

	if !hasAttr {
		return ErrAuthAttrNotFound
	}

	return m.Decide(token, attr, subject)
}

func (m *Manager) Decide(token Token, attr string, subject any) error {
	var err *errVoteStack

	for _, voter := range m.voters {
		if !voter.Supports(attr, subject) {
			continue
		}

		errVote := voter.Vote(token, attr, subject)

		if errVote == nil {
			return nil // если хотя бы один разрешил
		}

		if err == nil {
			err = &errVoteStack{errs: []error{errVote}}
		} else {
			err.errs = append(err.errs, errVote)
		}
	}

	if err != nil {
		return err
	}

	return ErrAuthNoVoters
}

type errVoteStack struct {
	errs []error
}

func (e *errVoteStack) Error() string {
	return e.errs[0].Error()
}

func (e *errVoteStack) Unwrap() []error {
	return e.errs
}
