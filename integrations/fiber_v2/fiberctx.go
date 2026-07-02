package inalog

import (
	"context"
	"time"
)

type fiberValueCtx struct {
	fctx context.Context
}

const userContextKey = "__local_user_context__"

func WithFiberCtx(fctx context.Context) context.Context {
	return &fiberValueCtx{fctx}
}

func (c *fiberValueCtx) Deadline() (deadline time.Time, ok bool) {
	return c.fctx.Deadline()
}

func (c *fiberValueCtx) Done() <-chan struct{} {
	return c.fctx.Done()
}

func (c *fiberValueCtx) Err() error {
	return c.fctx.Err()
}

func (c *fiberValueCtx) Value(key any) any {
	v := c.fctx.Value(key)
	if v == nil {
		if userCtx, ok := c.fctx.Value(userContextKey).(context.Context); ok {
			return userCtx.Value(key)
		}
	}
	return v
}
