package world

import (
	"syscall/js"

	"github.com/ByteArena/box2d"
)

// StickyInfo stores the two objects to be welded together after worldState.World.Step()
type StickyInfo struct {
	bodyA *box2d.B2Body
	bodyB *box2d.B2Body
}

type WorldSettings struct {
	SimSpeed   float64
	Height     float64
	Width      float64
	WorldScale float64
	Resizing   bool
	LastResize float64
}

type WorldState struct {
	*WorldSettings
	*JSObjects
	Player                  *box2d.B2Body
	WeldedDebris            *box2d.B2Body
	GoalBlock               *box2d.B2Body
	PlayerJoint             box2d.B2JointInterface
	PlayerCollisionDetected bool
	PlayerWelded            bool
	StickyArray             []StickyInfo
	World                   *box2d.B2World
	ResetWorld              bool
	TMark                   float64
}

func Initialize() *WorldState {
	worldState := &WorldState{
		WorldSettings: &WorldSettings{},
		JSObjects:     &JSObjects{},
	}

	worldState.Doc = js.Global().Get("document")
	worldState.Canvas = worldState.Doc.Call("getElementById", "mycanvas")
	worldState.Context = worldState.Canvas.Call("getContext", "2d")

	// create WorldSettings
	worldSettings := &WorldSettings{
		SimSpeed:   1,
		WorldScale: 0.0125,
		Width:      worldState.Doc.Get("body").Get("clientWidth").Float(),
		Height:     worldState.Doc.Get("body").Get("clientHeight").Float(),
	}

	world := box2d.MakeB2World(box2d.B2Vec2{X: 0, Y: 0})

	// create WorldState
	worldState.WorldSettings = worldSettings
	worldState.World = &world

	worldState.Size()

	worldState.Populate()

	return worldState
}

func (worldState *WorldState) Reset() {
	worldState.Clear()
	worldState.Populate()
	worldState.ResetWorld = false
}

func (worldState *WorldState) Clear() {
	// clear out world of any elements
	for joint := worldState.World.GetJointList(); joint != nil; joint = joint.GetNext() {
		worldState.World.DestroyJoint(joint)
	}

	for body := worldState.World.GetBodyList(); body != nil; body = body.GetNext() {
		worldState.World.DestroyBody(body)
	}
}

func (worldState WorldState) IsPlayerOutOfBounds() bool {
	return worldState.Player.GetPosition().X < 0 ||
		worldState.Player.GetPosition().X > worldState.Width*worldState.WorldScale ||
		worldState.Player.GetPosition().Y < 0 ||
		worldState.Player.GetPosition().Y > worldState.Height*worldState.WorldScale
}

func (worldState *WorldState) Resize() bool {
	// Poll window size to handle resize
	curBodyW := worldState.Doc.Get("body").Get("clientWidth").Float()
	curBodyH := worldState.Doc.Get("body").Get("clientHeight").Float()
	if curBodyW != worldState.Width || curBodyH != worldState.Height {
		worldState.Width, worldState.Height = curBodyW, curBodyH
		worldState.Size()
		return true
	}
	return false
}

func (worldState *WorldState) Size() {
	// size
	worldState.Canvas.Set("width", worldState.Width)
	worldState.Canvas.Set("height", worldState.Height)

	// scale
	worldState.Context.Call("scale", 1/worldState.WorldScale, 1/worldState.WorldScale)

	// style
	worldState.Context.Set("fillStyle", "rgba(100,100,100,1)")
	worldState.Context.Set("strokeStyle", "rgba(100,100,100,1)")
	worldState.Context.Set("lineWidth", 2*worldState.WorldScale)
}

func (worldState WorldState) GetSmallestDimension() float64 {
	if worldState.Width > worldState.Height {
		return worldState.Height
	}
	return worldState.Width
}
