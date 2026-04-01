module github.com/kimyechan/ebiten-aio-framework/examples/go/ebiten-ui-showcase

go 1.25.0

require (
	github.com/hajimehoshi/ebiten/v2 v2.9.9
	github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug v0.0.0
	github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui v0.0.0
	github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui-debug v0.0.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20250923094054-ea854a63cce1 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/image v0.38.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui => ../../../libs/go/ebiten-ui

replace github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-debug => ../../../libs/go/ebiten-debug

replace github.com/kimyechan/ebiten-aio-framework/libs/go/ebiten-ui-debug => ../../../libs/go/ebiten-ui-debug
