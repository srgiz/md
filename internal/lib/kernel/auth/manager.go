package auth

import "errors"

var ErrNoVoters = errors.New("auth: no voters")

type Voter interface {
	Supports(attr string, cmd any) bool
	Vote(token Token, attr string, cmd any) error
}

type Token interface {
	Data() map[string]any
}

type manager struct {
	voters []Voter
}

func (m *manager) add(voters ...Voter) {
	m.voters = append(m.voters, voters...)
}

func (m *manager) allow(token Token, attr string, cmd any) error {
	var err *errStackVote

	for _, voter := range m.voters {
		if !voter.Supports(attr, cmd) {
			continue
		}

		errVote := voter.Vote(token, attr, cmd)

		if errVote == nil {
			return nil // если хотя бы один разрешил
		}

		if err == nil {
			err = &errStackVote{errs: []error{errVote}}
		} else {
			err.errs = append(err.errs, errVote)
		}
	}

	if err != nil {
		return err
	}

	return ErrNoVoters
}

type errStackVote struct {
	errs []error
}

func (e *errStackVote) Error() string {
	return e.errs[0].Error()
}

func (e *errStackVote) Unwrap() []error {
	return e.errs
}
