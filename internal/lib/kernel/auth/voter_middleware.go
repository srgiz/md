package auth

import (
	"context"
	"errors"
	"md/internal/lib/kernel/cmdbus"
)

const ContextKeyToken = "auth.Token"
const ContextKeyAttr = "auth.Attr"

var ErrTokenNotFound = errors.New("auth: token not found")
var ErrAttrNotFound = errors.New("auth: attr not found")

type VoterMiddleware struct {
	manager *manager
}

func NewVoterMiddleware() *VoterMiddleware {
	return &VoterMiddleware{manager: &manager{}}
}

func (m *VoterMiddleware) Handle(ctx context.Context, cmd any, next cmdbus.Handler) (any, error) {
	token, hasToken := ctx.Value(ContextKeyToken).(Token)

	if !hasToken {
		return nil, ErrTokenNotFound
	}

	attr, hasAttr := ctx.Value(ContextKeyAttr).(string)

	if !hasAttr {
		return nil, ErrAttrNotFound
	}

	if err := m.manager.allow(token, attr, cmd); err != nil {
		return nil, err
	}

	return next(ctx, cmd)
}

func (m *VoterMiddleware) AddVoters(voters ...Voter) {
	m.manager.add(voters...)
}
