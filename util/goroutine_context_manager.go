package util

import (
	"context"
)

type GoroutineContextManager struct {
	gls GoroutineLocalStorage
}

func (gcm *GoroutineContextManager) SetContext(ctx *context.Context) {
	key := gcm.gls.GetGoroutineId()
	gcm.gls.Set(key, ctx)
}

func (gcm *GoroutineContextManager) DelContext() {
	key := gcm.gls.GetGoroutineId()
	gcm.gls.Del(key)
}

func (gcm *GoroutineContextManager) GetContext() (*context.Context, bool) {
	key := gcm.gls.GetGoroutineId()
	ctx, ok := gcm.gls.Get(key).(*context.Context)
	return ctx, ok
}
