package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestToken struct {
}

func (t *TestToken) Data() map[string]any {
	return make(map[string]any)
}

type FirstVoter struct {
}

func (v *FirstVoter) Supports(attr string, subject any) bool {
	return attr == "first" || attr == "last"
}

var errFirst = errors.New("first error")

func (v *FirstVoter) Vote(token Token, attr string, subject any) error {
	if attr == "first" {
		return nil
	}

	return errFirst
}

type SecondVoter struct {
}

func (v *SecondVoter) Supports(attr string, subject any) bool {
	return attr == "first" || attr == "second" || attr == "last"
}

var errSecond = errors.New("second error")

func (v *SecondVoter) Vote(token Token, attr string, subject any) error {
	return errSecond
}

func TestDecide(t *testing.T) {
	var tests = []struct {
		name string
		ctx  context.Context
		err  error
	}{
		{
			name: "no token",
			ctx:  context.Background(),
			err:  ErrAuthTokenNotFound,
		},
		{
			name: "no attr",
			ctx:  context.WithValue(context.Background(), ContextKeyAuthToken, &TestToken{}),
			err:  ErrAuthAttrNotFound,
		},
		{
			name: "unknown attr",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyAuthToken, &TestToken{}), ContextKeyAuthAttr, "unknown"),
			err:  ErrAuthNoVoters,
		},
		{
			name: "allow first",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyAuthToken, &TestToken{}), ContextKeyAuthAttr, "first"),
			err:  nil,
		},
		{
			name: "deny second",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyAuthToken, &TestToken{}), ContextKeyAuthAttr, "second"),
			err:  errSecond,
		},
		{
			name: "deny last 1",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyAuthToken, &TestToken{}), ContextKeyAuthAttr, "last"),
			err:  errFirst,
		},
		{
			name: "deny last 2",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyAuthToken, &TestToken{}), ContextKeyAuthAttr, "last"),
			err:  errSecond,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Example %s", test.name), func(t *testing.T) {
			m := NewManager()
			m.AddVoters(&FirstVoter{}, &SecondVoter{})

			err := m.DecideContext(test.ctx, nil)

			assert.ErrorIs(t, err, test.err)
		})
	}
}
