package util

import "context"

func IsAuthorized(ctx context.Context) bool {
	value := ctx.Value("authorized")
	return value != nil && value.(bool)
}
