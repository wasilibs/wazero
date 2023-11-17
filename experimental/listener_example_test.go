package experimental_test

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"sort"

	wazero "github.com/wasilibs/wazerox"
	"github.com/wasilibs/wazerox/api"
	"github.com/wasilibs/wazerox/experimental"
	"github.com/wasilibs/wazerox/imports/wasi_snapshot_preview1"
	"github.com/wasilibs/wazerox/internal/wasm"
)

// listenerWasm was generated by the following:
//
//	cd testdata; wat2wasm --debug-names listener.wat
//
//go:embed logging/testdata/listener.wasm
var listenerWasm []byte

// uniqGoFuncs implements both FunctionListenerFactory and FunctionListener
type uniqGoFuncs map[string]struct{}

// callees returns the go functions called.
func (u uniqGoFuncs) callees() []string {
	ret := make([]string, 0, len(u))
	for k := range u {
		ret = append(ret, k)
	}
	// Sort names for consistent iteration
	sort.Strings(ret)
	return ret
}

// NewFunctionListener implements FunctionListenerFactory.NewFunctionListener
func (u uniqGoFuncs) NewFunctionListener(def api.FunctionDefinition) experimental.FunctionListener {
	if def.GoFunction() == nil {
		return nil // only track go funcs
	}
	return u
}

// Before implements FunctionListener.Before
func (u uniqGoFuncs) Before(ctx context.Context, _ api.Module, def api.FunctionDefinition, _ []uint64, _ experimental.StackIterator) {
	u[def.DebugName()] = struct{}{}
}

// After implements FunctionListener.After
func (u uniqGoFuncs) After(context.Context, api.Module, api.FunctionDefinition, []uint64) {}

// Abort implements FunctionListener.Abort
func (u uniqGoFuncs) Abort(context.Context, api.Module, api.FunctionDefinition, error) {}

// This shows how to make a listener that counts go function calls.
func Example_customListenerFactory() {
	u := uniqGoFuncs{}

	// Set context to one that has an experimental listener
	ctx := context.WithValue(context.Background(), experimental.FunctionListenerFactoryKey{}, u)

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	mod, err := r.Instantiate(ctx, listenerWasm)
	if err != nil {
		log.Panicln(err)
	}

	for i := 0; i < 5; i++ {
		if _, err = mod.ExportedFunction("rand").Call(ctx, 4); err != nil {
			log.Panicln(err)
		}
	}

	// A Go function was called multiple times, but we should only see it once.
	for _, f := range u.callees() {
		fmt.Println(f)
	}

	// Output:
	// wasi_snapshot_preview1.fd_write
	// wasi_snapshot_preview1.random_get
}

func Example_stackIterator() {
	it := &fakeStackIterator{}

	for it.Next() {
		fn := it.Function()
		pc := it.ProgramCounter()
		fmt.Println("function:", fn.Definition().DebugName())
		fmt.Println("\tparameters:", it.Parameters())
		fmt.Println("\tprogram counter:", pc)
		fmt.Println("\tsource offset:", fn.SourceOffsetForPC(pc))
	}

	// Output:
	// function: fn0
	// 	parameters: [1 2 3]
	// 	program counter: 5890831
	// 	source offset: 1234
	// function: fn1
	// 	parameters: []
	// 	program counter: 5899822
	// 	source offset: 7286
	// function: fn2
	// 	parameters: [4]
	// 	program counter: 6820312
	// 	source offset: 935891
}

type fakeStackIterator struct {
	iteration    int
	def          api.FunctionDefinition
	args         []uint64
	pc           uint64
	sourceOffset uint64
}

func (s *fakeStackIterator) Next() bool {
	switch s.iteration {
	case 0:
		s.def = &mockFunctionDefinition{debugName: "fn0"}
		s.args = []uint64{1, 2, 3}
		s.pc = 5890831
		s.sourceOffset = 1234
	case 1:
		s.def = &mockFunctionDefinition{debugName: "fn1"}
		s.args = []uint64{}
		s.pc = 5899822
		s.sourceOffset = 7286
	case 2:
		s.def = &mockFunctionDefinition{debugName: "fn2"}
		s.args = []uint64{4}
		s.pc = 6820312
		s.sourceOffset = 935891
	case 3:
		return false
	}
	s.iteration++
	return true
}

func (s *fakeStackIterator) Function() experimental.InternalFunction {
	return internalFunction{
		definition:   s.def,
		sourceOffset: s.sourceOffset,
	}
}

func (s *fakeStackIterator) Parameters() []uint64 {
	return s.args
}

func (s *fakeStackIterator) ProgramCounter() experimental.ProgramCounter {
	return experimental.ProgramCounter(s.pc)
}

var _ experimental.StackIterator = &fakeStackIterator{}

type internalFunction struct {
	definition   api.FunctionDefinition
	sourceOffset uint64
}

func (f internalFunction) Definition() api.FunctionDefinition {
	return f.definition
}

func (f internalFunction) SourceOffsetForPC(pc experimental.ProgramCounter) uint64 {
	return f.sourceOffset
}

type mockFunctionDefinition struct {
	debugName string
	*wasm.FunctionDefinition
}

func (f *mockFunctionDefinition) DebugName() string {
	return f.debugName
}

func (f *mockFunctionDefinition) ParamTypes() []wasm.ValueType {
	return []wasm.ValueType{}
}

func (f *mockFunctionDefinition) ResultTypes() []wasm.ValueType {
	return []wasm.ValueType{}
}
