package wazevo

import (
	"runtime"

	"github.com/wasilibs/wazerox/internal/engine/wazevo/backend"
	"github.com/wasilibs/wazerox/internal/engine/wazevo/backend/isa/arm64"
)

func newMachine() backend.Machine {
	switch runtime.GOARCH {
	case "arm64":
		return arm64.NewBackend()
	default:
		panic("unsupported architecture")
	}
}
