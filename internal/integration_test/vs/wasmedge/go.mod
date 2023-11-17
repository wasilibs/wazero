module github.com/wasilibs/wazerox/internal/integration_test/vs/wasmedge

go 1.19

require (
	github.com/second-state/WasmEdge-go v0.12.1
	github.com/wasilibs/wazerox v0.0.0
)

replace github.com/wasilibs/wazerox => ../../../..
