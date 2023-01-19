package users

import (
	"context"
)

type groupUUIDKeyType struct{}

var (
	groupUUIDKey = groupUUIDKeyType{}
)

// NewContextWithGroupUUID ...
func NewContextWithGroupUUID(ctx context.Context, groupUUID string) context.Context {
	return context.WithValue(ctx, groupUUIDKey, groupUUID)
}

// GroupUUIDFromContext ...
func GroupUUIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(groupUUIDKey).(string)
	if !ok {
		return ""
	}

	return val
}
