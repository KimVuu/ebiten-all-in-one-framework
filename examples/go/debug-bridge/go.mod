module github.com/kimyechan/ebiten-aio-framework/examples/go/debug-bridge

go 1.25.0

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/kimyechan/ebiten-aio-framework/libs/go/ebitendebug v0.0.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace github.com/kimyechan/ebiten-aio-framework/libs/go/ebitendebug => ../../../libs/go/ebitendebug
