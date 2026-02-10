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

func (v *FirstVoter) Supports(attr string, cmd any) bool {
	return attr == "first" || attr == "last"
}

var errFirst = errors.New("first error")

func (v *FirstVoter) Vote(token Token, attr string, cmd any) error {
	if attr == "first" {
		return nil
	}

	return errFirst
}

type SecondVoter struct {
}

func (v *SecondVoter) Supports(attr string, cmd any) bool {
	return attr == "first" || attr == "second" || attr == "last"
}

var errSecond = errors.New("second error")

func (v *SecondVoter) Vote(token Token, attr string, cmd any) error {
	return errSecond
}

func TestValidateCreateUser(t *testing.T) {
	var tests = []struct {
		name string
		ctx  context.Context
		err  error
	}{
		{
			name: "no token",
			ctx:  context.Background(),
			err:  ErrTokenNotFound,
		},
		{
			name: "no attr",
			ctx:  context.WithValue(context.Background(), ContextKeyToken, &TestToken{}),
			err:  ErrAttrNotFound,
		},
		{
			name: "unknown attr",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyToken, &TestToken{}), ContextKeyAttr, "unknown"),
			err:  ErrNoVoters,
		},
		{
			name: "allow first",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyToken, &TestToken{}), ContextKeyAttr, "first"),
			err:  nil,
		},
		{
			name: "deny second",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyToken, &TestToken{}), ContextKeyAttr, "second"),
			err:  errSecond,
		},
		{
			name: "deny last 1",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyToken, &TestToken{}), ContextKeyAttr, "last"),
			err:  errFirst,
		},
		{
			name: "deny last 2",
			ctx:  context.WithValue(context.WithValue(context.Background(), ContextKeyToken, &TestToken{}), ContextKeyAttr, "last"),
			err:  errSecond,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Example %s", test.name), func(t *testing.T) {
			mw := NewVoterMiddleware()
			mw.AddVoters(&FirstVoter{}, &SecondVoter{})

			_, err := mw.Handle(test.ctx, 0, func(ctx context.Context, cmd any) (any, error) {
				return 0, nil
			})

			assert.ErrorIs(t, err, test.err)
		})
	}
}
