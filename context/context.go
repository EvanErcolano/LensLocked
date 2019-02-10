package context

import (
	"context"

	"lenslocked.com/models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

// WithUser attaches a user to the provided context
func WithUser(ctx context.Context, user *models.User) context.Context {
	// withValue looks up using the type and value of the key
	// by making our own type and no exporting it we prevent someone
	// from overwriting our user by passing in "user"
	return context.WithValue(ctx, userKey, user)
}

// User pulls a user from a context
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			// ^^ temp.(type) is a type assertion
			return user
		}
	}
	return nil
}
