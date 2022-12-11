package cmd

import (
	"context"
)

type CmdHandler interface {
	HandleCmd(ctx context.Context, payload CmdPayload) (string, error)
}

type CmdPayload struct {
	Args     []string
	UserName string
	ChatID   int64
}
