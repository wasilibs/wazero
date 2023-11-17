package wazevo

import (
	"context"
	"testing"

	"github.com/wasilibs/wazerox/api"
	"github.com/wasilibs/wazerox/internal/testing/require"
)

func Test_writeIface_readIface(t *testing.T) {
	buf := make([]byte, 100)

	var called bool
	var goFn api.GoFunction = api.GoFunc(func(context.Context, []uint64) {
		called = true
	})
	writeIface(goFn, buf)
	got := readIface(buf).(api.GoFunction)
	got.Call(context.Background(), nil)
	require.True(t, called)
}
