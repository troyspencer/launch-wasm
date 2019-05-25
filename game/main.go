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

	// handle player clicks
	mouseDownEvt := js.FuncOf(worldState.HandleClick)
	defer mouseDownEvt.Release()

	keyUpEvt := js.FuncOf(worldState.HandleEsc)
	defer keyUpEvt.Release()

	worldState.Doc.Call("addEventListener", "keyup", keyUpEvt)
	worldState.Doc.Call("addEventListener", "mousedown", mouseDownEvt)

	done := make(chan struct{}, 0)
	// Start running
	renderFrame := js.FuncOf(worldState.RenderFrame)
	js.Global().Call("requestAnimationFrame", renderFrame)
	<-done
}
