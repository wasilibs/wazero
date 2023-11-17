package frontend

import (
	"testing"
	"unsafe"

	"github.com/wasilibs/wazerox/internal/testing/require"
	"github.com/wasilibs/wazerox/internal/wasm"
)

func Test_Offsets(t *testing.T) {
	var globalInstance wasm.GlobalInstance
	require.Equal(t, int(unsafe.Offsetof(globalInstance.Val)), globalInstanceValueOffset)
	var memInstance wasm.MemoryInstance
	require.Equal(t, int(unsafe.Offsetof(memInstance.Buffer)), memoryInstanceBufOffset)
	var tableInstance wasm.TableInstance
	require.Equal(t, int(unsafe.Offsetof(tableInstance.References)), tableInstanceBaseAddressOffset)
}
