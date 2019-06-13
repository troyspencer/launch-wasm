package main

import (
	"math/rand"
	"syscall/js"
	"time"

	"github.com/troyspencer/launch-wasm/game/contact"
	"github.com/troyspencer/launch-wasm/game/world"
)

func main() {
	// seed the random generator
	rand.Seed(time.Now().UnixNano())

	worldState := world.Initialize()

	worldState.World.SetContactListener(&contact.PlayerContactListener{WorldState: worldState})

	// link callbacks
	pauseEvent := js.FuncOf(worldState.HandlePause)
	defer pauseEvent.Release()
	worldState.Doc.Call("addEventListener", "pause", pauseEvent)

	unpauseEvent := js.FuncOf(worldState.HandleUnpause)
	defer unpauseEvent.Release()
	worldState.Doc.Call("addEventListener", "unpause", unpauseEvent)

	mouseDownEvent := js.FuncOf(worldState.HandleClick)
	defer mouseDownEvent.Release()
	worldState.Doc.Call("addEventListener", "mousedown", mouseDownEvent)

	keyUpEvent := js.FuncOf(worldState.HandleKeys)
	defer keyUpEvent.Release()
	worldState.Doc.Call("addEventListener", "keyup", keyUpEvent)

	done := make(chan struct{}, 0)
	// Start running
	renderFrame := js.FuncOf(worldState.RenderFrame)
	js.Global().Call("requestAnimationFrame", renderFrame)
	<-done
}
