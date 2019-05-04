//Wasming
// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/ByteArena/box2d" // this box2d throws some unexpected panics
)

type JSObjects struct {
	Context js.Value
	Doc     js.Value
	Canvas  js.Value
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
	TMark float64
}

func main() {
	// seed the random generator
	rand.Seed(time.Now().UnixNano())

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
	worldState.Canvas.Call("setAttribute", "width", worldSettings.Width)
	worldState.Canvas.Call("setAttribute", "height", worldSettings.Height)

	worldState.Context = worldState.Canvas.Call("getContext", "2d")
	worldState.Context.Call("scale", 1/worldSettings.WorldScale, 1/worldSettings.WorldScale)

	done := make(chan struct{}, 0)

	worldState.World.SetContactListener(&playerContactListener{WorldState: worldState})

	populateWorld(worldState)

	// handle player clicks
	mouseDownEvt := js.NewCallback(func(args []js.Value) {

		e := args[0]
		if e.Get("target") != worldState.Canvas {
			return
		}

		// only allow launch if grounded aka welded to an object
		if worldState.PlayerWelded {
			mx := e.Get("clientX").Float() * worldState.WorldScale
			my := e.Get("clientY").Float() * worldState.WorldScale

			movementDx := mx - worldState.Player.GetPosition().X
			movementDy := my - worldState.Player.GetPosition().Y

			// create normalized movement vector from player to click location
			impulseVelocity := box2d.B2Vec2{X: movementDx, Y: movementDy}
			impulseVelocity.Normalize()
			impulseVelocity.OperatorScalarMulInplace(getSmallestDimension(worldState) * worldState.WorldScale / 2)

			clearPlayerJoints(worldState)
			worldState.PlayerWelded = false

			if worldState.WeldedDebris.GetType() == box2d.B2BodyType.B2_dynamicBody {

				// get current player velocity
				playerCurrentVelocity := worldState.Player.GetLinearVelocity()

				// calculate difference between current player velocity and player desired velocity
				velocityDisplacement := box2d.B2Vec2{
					X: impulseVelocity.X - playerCurrentVelocity.X,
					Y: impulseVelocity.Y - playerCurrentVelocity.Y}

				// calculate momentum of player
				momentum := velocityDisplacement.Length() * worldState.Player.GetMass()

				// calculate magnitude of debris velocity
				debrisVelocityDisplacementMagnitude := momentum / worldState.WeldedDebris.GetMass()

				// calculate velocity displacement of debris from momentum
				debrisVelocityDisplacement := velocityDisplacement
				debrisVelocityDisplacement.Normalize()
				debrisVelocityDisplacement.OperatorScalarMulInplace(debrisVelocityDisplacementMagnitude)
				debrisVelocityDisplacement = debrisVelocityDisplacement.OperatorNegate()

				// get debris current velocity, which should match the players current velocity due to welding
				debrisCurrentVelocity := worldState.WeldedDebris.GetLinearVelocity()

				// calculate resultant from debris current velocity and debris velocity displacement
				debrisVelocity := box2d.B2Vec2{
					X: debrisCurrentVelocity.X + debrisVelocityDisplacement.X,
					Y: debrisCurrentVelocity.Y + debrisVelocityDisplacement.Y,
				}

				// update debris velocity
				worldState.WeldedDebris.SetLinearVelocity(debrisVelocity)
			}

			// set player velocity to player desired velocity
			worldState.Player.SetLinearVelocity(impulseVelocity)
		}

	})
	defer mouseDownEvt.Release()

	keyUpEvt := js.NewCallback(func(args []js.Value) {
		e := args[0]
		if e.Get("which").Int() == 27 {
			worldState.ResetWorld = true
		}
	})
	defer keyUpEvt.Release()

	worldState.Doc.Call("addEventListener", "keyup", keyUpEvt)
	worldState.Doc.Call("addEventListener", "mousedown", mouseDownEvt)

	// overall style
	worldState.Context.Set("fillStyle", "rgba(100,100,100,1)")
	worldState.Context.Set("strokeStyle", "rgba(100,100,100,1)")
	worldState.Context.Set("lineWidth", 2*worldState.WorldScale)

	// Start running
	js.Global().Call("requestAnimationFrame", js.NewCallback(worldState.RenderFrame))
	<-done
}

func (worldState *WorldState) RenderFrame(args []js.Value) {
	now := args[0].Float()
	tdiff := now - worldState.TMark
	worldState.TMark = now

	// Poll window size to handle resize
	curBodyW := worldState.Doc.Get("body").Get("clientWidth").Float()
	curBodyH := worldState.Doc.Get("body").Get("clientHeight").Float()
	if curBodyW != worldState.Width || curBodyH != worldState.Height {
		worldState.Width, worldState.Height = curBodyW, curBodyH
		worldState.Canvas.Set("width", worldState.Width)
		worldState.Canvas.Set("height", worldState.Height)
	}

	worldState.World.Step(tdiff/1000*worldState.SimSpeed, 60, 120)

	checkPlayerOutOfBounds(worldState)

	if worldState.ResetWorld {
		clearWorld(worldState)
		populateWorld(worldState)
		worldState.ResetWorld = false
	}

	// check for new weld joint and execute it
	for len(worldState.StickyArray) > 0 {
		stickyBody := worldState.StickyArray[0]
		worldState.StickyArray[0] = worldState.StickyArray[len(worldState.StickyArray)-1]
		worldState.StickyArray = worldState.StickyArray[:len(worldState.StickyArray)-1]

		worldCoordsAnchorPoint := stickyBody.bodyB.GetWorldPoint(box2d.B2Vec2{X: 0, Y: 0})

		weldJointDef := box2d.MakeB2WeldJointDef()
		weldJointDef.BodyA = stickyBody.bodyA
		weldJointDef.BodyB = stickyBody.bodyB
		weldJointDef.ReferenceAngle = weldJointDef.BodyB.GetAngle() - weldJointDef.BodyA.GetAngle()
		weldJointDef.LocalAnchorA = weldJointDef.BodyA.GetLocalPoint(worldCoordsAnchorPoint)
		weldJointDef.LocalAnchorB = weldJointDef.BodyB.GetLocalPoint(worldCoordsAnchorPoint)

		if worldState.PlayerCollisionDetected {
			worldState.PlayerJoint = worldState.World.CreateJoint(&weldJointDef)
			worldState.PlayerCollisionDetected = false
			worldState.PlayerWelded = true
		} else {
			worldState.World.CreateJoint(&weldJointDef)
		}

	}

	worldState.Context.Call("clearRect", 0, 0, worldState.Width*worldState.WorldScale, worldState.Height*worldState.WorldScale)

	for curBody := worldState.World.GetBodyList(); curBody != nil; curBody = curBody.M_next {
		// ignore player and goal block, as they are styled differently
		if curBody.GetUserData() == "player" {
			// Player ball color
			worldState.Context.Set("fillStyle", "rgba(180, 180,180,1)")
			worldState.Context.Set("strokeStyle", "rgba(180,180,180,1)")
		} else if curBody.GetUserData() == "goalBlock" {
			// Goal block color
			worldState.Context.Set("fillStyle", "rgba(0, 255,0,1)")
			worldState.Context.Set("strokeStyle", "rgba(0,255,0,1)")
		} else {
			// color for other objects
			worldState.Context.Set("fillStyle", "rgba(100,100,100,1)")
			worldState.Context.Set("strokeStyle", "rgba(100,100,100,1)")
		}

		// Only one fixture for now
		worldState.Context.Call("save")
		ft := curBody.M_fixtureList
		switch shape := ft.M_shape.(type) {
		case *box2d.B2PolygonShape: // Box
			// canvas translate
			worldState.Context.Call("translate", curBody.M_xf.P.X, curBody.M_xf.P.Y)
			worldState.Context.Call("rotate", curBody.M_xf.Q.GetAngle())
			worldState.Context.Call("beginPath")
			worldState.Context.Call("moveTo", shape.M_vertices[0].X, shape.M_vertices[0].Y)
			for _, v := range shape.M_vertices[1:shape.M_count] {
				worldState.Context.Call("lineTo", v.X, v.Y)
			}
			worldState.Context.Call("lineTo", shape.M_vertices[0].X, shape.M_vertices[0].Y)
			worldState.Context.Call("fill")
			worldState.Context.Call("stroke")
		case *box2d.B2CircleShape:
			worldState.Context.Call("translate", curBody.M_xf.P.X, curBody.M_xf.P.Y)
			worldState.Context.Call("rotate", curBody.M_xf.Q.GetAngle())
			worldState.Context.Call("beginPath")
			worldState.Context.Call("arc", 0, 0, shape.M_radius, 0, 2*math.Pi)
			worldState.Context.Call("fill")
			worldState.Context.Call("moveTo", 0, 0)
			worldState.Context.Call("lineTo", 0, shape.M_radius)
			worldState.Context.Call("stroke")
		}
		worldState.Context.Call("restore")
	}
	js.Global().Call("requestAnimationFrame", js.NewCallback(worldState.RenderFrame))
}

type playerContactListener struct {
	*WorldState
}

func (listener playerContactListener) BeginContact(contact box2d.B2ContactInterface) {

	worldState := listener.WorldState

	// wait for bodies to actually contact
	if contact.IsTouching() {

		// detect player collision
		if contact.GetFixtureB().GetBody().GetUserData() == "player" || contact.GetFixtureA().GetBody().GetUserData() == "player" {

			// check which fixture is the debris
			if contact.GetFixtureA().GetBody().GetUserData() == "player" {
				worldState.WeldedDebris = contact.GetFixtureB().GetBody()
			} else {
				worldState.WeldedDebris = contact.GetFixtureA().GetBody()
			}

			if worldState.WeldedDebris.GetUserData() == "goalBlock" {
				worldState.ResetWorld = true
				return
			}

			// If player has already collided with another object this frame
			// ignore this collision
			if !worldState.PlayerCollisionDetected && !worldState.PlayerWelded {
				worldState.PlayerCollisionDetected = true
				weldContact(worldState, contact)
			}
		}
	}
}

func (listener playerContactListener) EndContact(contact box2d.B2ContactInterface) {}

func (listener playerContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {

}

func (listener playerContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {

}

func weldContact(worldState *WorldState, contact box2d.B2ContactInterface) {
	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)

	worldState.StickyArray = append(worldState.StickyArray, StickyInfo{
		bodyA: contact.GetFixtureA().GetBody(),
		bodyB: contact.GetFixtureB().GetBody(),
	})
}

// StickyInfo stores the two objects to be welded together after worldState.World.Step()
type StickyInfo struct {
	bodyA *box2d.B2Body
	bodyB *box2d.B2Body
}

func checkPlayerOutOfBounds(worldState *WorldState) {
	if worldState.Player.GetPosition().X < 0 || worldState.Player.GetPosition().X > worldState.Width*worldState.WorldScale || worldState.Player.GetPosition().Y < 0 || worldState.Player.GetPosition().Y > worldState.Height*worldState.WorldScale {
		worldState.ResetWorld = true
	}
}

func clearWorld(worldState *WorldState) {
	// clear out world of any elements
	for joint := worldState.World.GetJointList(); joint != nil; joint = joint.GetNext() {
		worldState.World.DestroyJoint(joint)
	}

	for body := worldState.World.GetBodyList(); body != nil; body = body.GetNext() {
		worldState.World.DestroyBody(body)
	}
}

func clearPlayerJoints(worldState *WorldState) {
	for jointEdge := worldState.Player.GetJointList(); jointEdge != nil; jointEdge = jointEdge.Next {
		worldState.World.DestroyJoint(jointEdge.Joint)
	}
}

func getSmallestDimension(worldState *WorldState) float64 {
	if worldState.Width > worldState.Height {
		return worldState.Height
	}
	return worldState.Width
}

func populateWorld(worldState *WorldState) {
	smallestDimension := getSmallestDimension(worldState)

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
