package wasm

import (
	"github.com/wasilibs/wazerox/api"
	"github.com/wasilibs/wazerox/internal/internalapi"
)

// constantGlobal wraps GlobalInstance to implement api.Global.
type constantGlobal struct {
	internalapi.WazeroOnlyType
	g *GlobalInstance
}

// Type implements api.Global.
func (g constantGlobal) Type() api.ValueType {
	return g.g.Type.ValType
}

// Get implements api.Global.
func (g constantGlobal) Get() uint64 {
	return g.g.Val
}

// String implements api.Global.
func (g constantGlobal) String() string {
	return g.g.String()
}

// mutableGlobal extends constantGlobal to allow updates.
type mutableGlobal struct {
	internalapi.WazeroOnlyType
	g *GlobalInstance
}

// Type implements api.Global.
func (g mutableGlobal) Type() api.ValueType {
	return g.g.Type.ValType
}

// Get implements api.Global.
func (g mutableGlobal) Get() uint64 {
	return g.g.Val
}

// String implements api.Global.
func (g mutableGlobal) String() string {
	return g.g.String()
}

// Set implements the same method as documented on api.MutableGlobal.
func (g mutableGlobal) Set(v uint64) {
	g.g.Val = v
}

// compile-time check to ensure mutableGlobal is a api.Global.
var _ api.Global = mutableGlobal{}
