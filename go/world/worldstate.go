package world

import (
	"math/rand"
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

	// Init Canvas stuff
	worldState.Canvas = worldState.Doc.Call("getElementById", "mycanvas")
	worldState.Canvas.Call("setAttribute", "width", worldState.Width)
	worldState.Canvas.Call("setAttribute", "height", worldState.Height)

	worldState.Context = worldState.Canvas.Call("getContext", "2d")
	worldState.Context.Call("scale", 1/worldSettings.WorldScale, 1/worldSettings.WorldScale)

	// overall style
	worldState.Context.Set("fillStyle", "rgba(100,100,100,1)")
	worldState.Context.Set("strokeStyle", "rgba(100,100,100,1)")
	worldState.Context.Set("lineWidth", 2*worldState.WorldScale)

	worldState.Populate()

	return worldState
}

func (worldState WorldState) CheckPlayerOutOfBounds() {
	if worldState.Player.GetPosition().X < 0 || worldState.Player.GetPosition().X > worldState.Width*worldState.WorldScale || worldState.Player.GetPosition().Y < 0 || worldState.Player.GetPosition().Y > worldState.Height*worldState.WorldScale {
		worldState.ResetWorld = true
	}
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

func (worldState *WorldState) ClearPlayerJoints() {
	for jointEdge := worldState.Player.GetJointList(); jointEdge != nil; jointEdge = jointEdge.Next {
		worldState.World.DestroyJoint(jointEdge.Joint)
	}
}

func (worldState WorldState) GetSmallestDimension() float64 {
	if worldState.Width > worldState.Height {
		return worldState.Height
	}
	return worldState.Width
}

func (worldState *WorldState) Populate() {
	smallestDimension := worldState.GetSmallestDimension()

	// Player Ball
	worldState.Player = worldState.World.CreateBody(&box2d.B2BodyDef{
		Type:         box2d.B2BodyType.B2_dynamicBody,
		Position:     box2d.B2Vec2{X: smallestDimension * worldState.WorldScale / 32, Y: worldState.Height*worldState.WorldScale - smallestDimension*worldState.WorldScale/32},
		Awake:        true,
		Active:       true,
		GravityScale: 1.0,
		Bullet:       true,
		UserData:     "player",
	})
	shape := box2d.NewB2CircleShape()
	shape.M_radius = smallestDimension * worldState.WorldScale / 64
	ft := worldState.Player.CreateFixture(shape, 1)
	ft.M_friction = 0
	ft.M_restitution = 1

	// Create launch block
	launchBlock := worldState.World.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_dynamicBody,
		Position: box2d.B2Vec2{X: smallestDimension * worldState.WorldScale / 32, Y: worldState.Height*worldState.WorldScale - smallestDimension*worldState.WorldScale/32},
		Active:   true,
		UserData: "launchBlock",
	})
	launchBlockShape := &box2d.B2PolygonShape{}
	launchBlockShape.SetAsBox(smallestDimension*worldState.WorldScale/32, smallestDimension*worldState.WorldScale/32)
	ft = launchBlock.CreateFixture(launchBlockShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	// Create goal block
	worldState.GoalBlock = worldState.World.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: worldState.Width*worldState.WorldScale - smallestDimension*worldState.WorldScale/32, Y: smallestDimension * worldState.WorldScale / 32},
		Active:   true,
		UserData: "goalBlock",
	})
	goalBlockShape := &box2d.B2PolygonShape{}
	goalBlockShape.SetAsBox(smallestDimension*worldState.WorldScale/32, smallestDimension*worldState.WorldScale/32)
	ft = worldState.GoalBlock.CreateFixture(goalBlockShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	// Some Random debris
	for i := 0; i < 25; i++ {
		obj1 := worldState.World.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_dynamicBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * worldState.Width * worldState.WorldScale,
				Y: rand.Float64() * worldState.Height * worldState.WorldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
			UserData:     "weldableDebris",
		})
		shape := &box2d.B2PolygonShape{}
		shape.SetAsBox(
			rand.Float64()*smallestDimension*worldState.WorldScale/10,
			rand.Float64()*smallestDimension*worldState.WorldScale/10)
		ft := obj1.CreateFixture(shape, 1)
		ft.M_friction = 1
		ft.M_restitution = 0 // bouncy
	}
}

func (worldState *WorldState) WeldContact(contact box2d.B2ContactInterface) {
	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)

	worldState.StickyArray = append(worldState.StickyArray, StickyInfo{
		bodyA: contact.GetFixtureA().GetBody(),
		bodyB: contact.GetFixtureB().GetBody(),
	})
}
