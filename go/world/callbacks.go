package world

import (
	"syscall/js"
)

func (worldState *WorldState) HandleEsc(args []js.Value) {
	e := args[0]
	if e.Get("which").Int() == 27 {
		worldState.ResetWorld = true
	}
}

func (worldState *WorldState) HandleClick(args []js.Value) {
	e := args[0]
	if e.Get("target") != worldState.Canvas {
		return
	}

	// only allow launch if grounded aka welded to an object
	if worldState.PlayerWelded {
		mx := e.Get("clientX").Float() * worldState.WorldScale
		my := e.Get("clientY").Float() * worldState.WorldScale
		worldState.LaunchPlayer(mx, my)
	}
}

func (worldState *WorldState) RenderFrame(args []js.Value) {
	now := args[0].Float()
	tdiff := now - worldState.TMark
	worldState.TMark = now

	worldState.Resize()

	worldState.World.Step(tdiff/1000*worldState.SimSpeed, 60, 120)

	if worldState.IsPlayerOutOfBounds() {
		worldState.ResetWorld = true
	}

	if worldState.ResetWorld {
		worldState.Reset()
	}

	// check for new weld joint and execute it
	worldState.WeldJoint()

	worldState.Context.Call("clearRect", 0, 0, worldState.Width*worldState.WorldScale, worldState.Height*worldState.WorldScale)

	for curBody := worldState.World.GetBodyList(); curBody != nil; curBody = curBody.M_next {
		// ignore player and goal block, as they are styled differently
		worldState.Draw(curBody)
	}
	js.Global().Call("requestAnimationFrame", js.NewCallback(worldState.RenderFrame))
}
