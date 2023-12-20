package utils

import "context"

type contextKey string

var ContextKeyUID contextKey

func (c contextKey) Stirng() string {
	return string(c)
}

func GetCtxUID(ctx context.Context) (string, bool) {
	UID, ok := ctx.Value(ContextKeyUID).(string)
	return UID, ok
}
