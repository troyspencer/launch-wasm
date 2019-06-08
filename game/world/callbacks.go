package world

import (
	"log"
	"syscall/js"

	"github.com/ByteArena/box2d"
)

func (worldState *WorldState) HandleEsc(this js.Value, args []js.Value) interface{} {
	e := args[0]
	if e.Get("which").Int() == 27 {
		worldState.ResetWorld = true
	}
	return nil
}

func (worldState *WorldState) HandleClick(this js.Value, args []js.Value) interface{} {
	e := args[0]
	if e.Get("target") != worldState.Canvas {
		return nil
	}

	// only allow launch if grounded aka welded to an object
	if worldState.PlayerWelded || worldState.AbsorbCount > 0 {
		mx := e.Get("clientX").Float() * worldState.WorldScale
		my := e.Get("clientY").Float() * worldState.WorldScale
		worldState.LaunchPlayer(mx, my)
	}

	// always clear joints on click to prevent sticking to multiple debris
	worldState.ClearJoints(worldState.Player)

	return nil
}

func (worldState *WorldState) RenderFrame(this js.Value, args []js.Value) interface{} {
	now := args[0].Float()
	tdiff := now - worldState.TMark
	worldState.TMark = now

	worldState.World.Step(tdiff/1000*worldState.SimSpeed, 60, 120)

	resizing := worldState.Resize()
	if resizing {
		worldState.Resizing = true
		worldState.LastResize = now
	}
	if now-worldState.LastResize > 150 && worldState.Resizing {
		worldState.ResetWorld = true
		worldState.Resizing = false
	}

	if worldState.IsPlayerOutOfBounds() {
		worldState.ResetWorld = true
	}

	if worldState.ResetWorld {
		worldState.Reset()
	}

	// check for new weld joint and execute it
	worldState.WeldJoint()

	worldState.Context.Call("clearRect", 0, 0, worldState.Width*worldState.WorldScale, worldState.Height*worldState.WorldScale)

	draw := []*box2d.B2Body{}
	drawLast := []*box2d.B2Body{}

	// draw absorbing things last, so they can be drawn in front of others
	for curBody := worldState.World.GetBodyList(); curBody != nil; curBody = curBody.M_next {
		if absorber, ok := curBody.GetUserData().(Absorber); ok {
			if absorber.Absorbs() {
				drawLast = append(drawLast, curBody)
			} else {
				draw = append(draw, curBody)
			}
		}
	}

	for _, body := range draw {
		worldState.Draw(body)
	}

	for _, body := range drawLast {
		worldState.Draw(body)
	}

	renderFrame := js.FuncOf(worldState.RenderFrame)
	js.Global().Call("requestAnimationFrame", renderFrame)
	return nil
}

func (worldState *WorldState) HandleButton(this js.Value, args []js.Value) interface{} {
	worldState.Clicks++
	log.Println(worldState.Clicks)
	return nil
}

type Absorber interface {
	Absorbs() bool
}
