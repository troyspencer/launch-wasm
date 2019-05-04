package main

import (
	"math/rand"
	"syscall/js"
	"time"

	"github.com/troyspencer/launch-wasm/go/contact"
	"github.com/troyspencer/launch-wasm/go/world"
)

func main() {
	// seed the random generator
	rand.Seed(time.Now().UnixNano())

	worldState := world.Initialize()

	worldState.World.SetContactListener(&contact.PlayerContactListener{WorldState: worldState})

	// handle player clicks
	mouseDownEvt := js.NewCallback(worldState.HandleClick)
	defer mouseDownEvt.Release()

	keyUpEvt := js.NewCallback(worldState.HandleEsc)
	defer keyUpEvt.Release()

	worldState.Doc.Call("addEventListener", "keyup", keyUpEvt)
	worldState.Doc.Call("addEventListener", "mousedown", mouseDownEvt)

	done := make(chan struct{}, 0)
	// Start running
	js.Global().Call("requestAnimationFrame", js.NewCallback(worldState.RenderFrame))
	<-done
}
