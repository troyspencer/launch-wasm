// Click on canvas to start a polygon
// - Max 8 vertices
// - Only convex polygons
// - Esc cancel polygon

//Wasming
// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"math"
	"math/rand"
	"syscall/js"

	// this box2d throws some unexpected panics
	"github.com/ByteArena/box2d"
)

var (
	width                   float64
	height                  float64
	ctx                     js.Value
	simSpeed                float64 = 1
	worldScale                      = 0.0125 // 1/8
	player                  *box2d.B2Body
	weldedDebris            *box2d.B2Body
	goalBlock               *box2d.B2Body
	playerJoint             box2d.B2JointInterface
	playerCollisionDetected bool
	playerWelded            bool
	stickyArray             []StickyInfo
	world                   box2d.B2World
	resetWorld              bool
)

func main() {

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "mycanvas")
	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float()
	canvasEl.Call("setAttribute", "width", width)
	canvasEl.Call("setAttribute", "height", height)
	ctx = canvasEl.Call("getContext", "2d")
	ctx.Call("scale", 1/worldScale, 1/worldScale)

	done := make(chan struct{}, 0)

	world = box2d.MakeB2World(box2d.B2Vec2{X: 0, Y: 0})

	world.SetContactListener(&playerContactListener{})

	populateWorld()

	// handle player clicks
	mouseDownEvt := js.NewCallback(func(args []js.Value) {

		e := args[0]
		if e.Get("target") != canvasEl {
			return
		}

		// only allow launch if grounded aka welded to an object
		if playerWelded {

			mx := e.Get("clientX").Float() * worldScale
			my := e.Get("clientY").Float() * worldScale

			movementDx := mx - player.GetPosition().X
			movementDy := my - player.GetPosition().Y

			// create normalized movement vector from player to click location
			impulseVelocity := box2d.B2Vec2{X: movementDx, Y: movementDy}
			impulseVelocity.Normalize()
			impulseVelocity.OperatorScalarMulInplace(5)

			clearPlayerJoints()
			playerWelded = false

			if weldedDebris.GetType() == box2d.B2BodyType.B2_dynamicBody {

				// get current player velocity
				playerCurrentVelocity := player.GetLinearVelocity()

				// calculate difference between current player velocity and player desired velocity
				velocityDisplacement := box2d.B2Vec2{
					X: impulseVelocity.X - playerCurrentVelocity.X,
					Y: impulseVelocity.Y - playerCurrentVelocity.Y}

				// calculate momentum of player
				momentum := velocityDisplacement.Length() * player.GetMass()

				// calculate magnitude of debris velocity
				debrisVelocityDisplacementMagnitude := momentum / weldedDebris.GetMass()

				// calculate velocity displacement of debris from momentum
				debrisVelocityDisplacement := velocityDisplacement
				debrisVelocityDisplacement.Normalize()
				debrisVelocityDisplacement.OperatorScalarMulInplace(debrisVelocityDisplacementMagnitude)
				debrisVelocityDisplacement = debrisVelocityDisplacement.OperatorNegate()

				// get debris current velocity, which should match the players current velocity due to welding
				debrisCurrentVelocity := weldedDebris.GetLinearVelocity()

				// calculate resultant from debris current velocity and debris velocity displacement
				debrisVelocity := box2d.B2Vec2{
					X: debrisCurrentVelocity.X + debrisVelocityDisplacement.X,
					Y: debrisCurrentVelocity.Y + debrisVelocityDisplacement.Y,
				}

				// update debris velocity
				weldedDebris.SetLinearVelocity(debrisVelocity)
			}

			// set player velocity to player desired velocity
			player.SetLinearVelocity(impulseVelocity)
		}

	})
	defer mouseDownEvt.Release()

	keyUpEvt := js.NewCallback(func(args []js.Value) {
		e := args[0]
		if e.Get("which").Int() == 27 {
			resetWorld = true
		}
	})
	defer keyUpEvt.Release()

	doc.Call("addEventListener", "keyup", keyUpEvt)
	doc.Call("addEventListener", "mousedown", mouseDownEvt)

	// Draw things
	var renderFrame js.Callback
	var tmark float64

	// overall style
	ctx.Set("fillStyle", "rgba(100,100,100,1)")
	ctx.Set("strokeStyle", "rgba(100,100,100,1)")
	ctx.Set("lineWidth", 2*worldScale)

	renderFrame = js.NewCallback(func(args []js.Value) {
		now := args[0].Float()
		tdiff := now - tmark
		tmark = now

		// Poll window size to handle resize
		curBodyW := doc.Get("body").Get("clientWidth").Float()
		curBodyH := doc.Get("body").Get("clientHeight").Float()
		if curBodyW != width || curBodyH != height {
			width, height = curBodyW, curBodyH
			canvasEl.Set("width", width)
			canvasEl.Set("height", height)
		}

		world.Step(tdiff/1000*simSpeed, 60, 120)

		checkPlayerOutOfBounds()

		if resetWorld {
			clearWorld()
			populateWorld()
			resetWorld = false
		}

		// check for new weld joint and execute it
		for len(stickyArray) > 0 {
			stickyBody := stickyArray[0]
			stickyArray[0] = stickyArray[len(stickyArray)-1]
			stickyArray = stickyArray[:len(stickyArray)-1]

			worldCoordsAnchorPoint := stickyBody.bodyB.GetWorldPoint(box2d.B2Vec2{X: 0, Y: 0})

			weldJointDef := box2d.MakeB2WeldJointDef()
			weldJointDef.BodyA = stickyBody.bodyA
			weldJointDef.BodyB = stickyBody.bodyB
			weldJointDef.ReferenceAngle = weldJointDef.BodyB.GetAngle() - weldJointDef.BodyA.GetAngle()
			weldJointDef.LocalAnchorA = weldJointDef.BodyA.GetLocalPoint(worldCoordsAnchorPoint)
			weldJointDef.LocalAnchorB = weldJointDef.BodyB.GetLocalPoint(worldCoordsAnchorPoint)

			if playerCollisionDetected {
				playerJoint = world.CreateJoint(&weldJointDef)
				playerCollisionDetected = false
				playerWelded = true
			} else {
				world.CreateJoint(&weldJointDef)
			}

		}

		ctx.Call("clearRect", 0, 0, width*worldScale, height*worldScale)

		for curBody := world.GetBodyList(); curBody != nil; curBody = curBody.M_next {
			// ignore player and goal block, as they are styled differently
			if curBody.GetUserData() == "player" {
				// Player ball color
				ctx.Set("fillStyle", "rgba(180, 180,180,1)")
				ctx.Set("strokeStyle", "rgba(180,180,180,1)")
			} else if curBody.GetUserData() == "goalBlock" {
				// Goal block color
				ctx.Set("fillStyle", "rgba(0, 255,0,1)")
				ctx.Set("strokeStyle", "rgba(0,255,0,1)")
			} else {
				// color for other objects
				ctx.Set("fillStyle", "rgba(100,100,100,1)")
				ctx.Set("strokeStyle", "rgba(100,100,100,1)")
			}

			// Only one fixture for now
			ctx.Call("save")
			ft := curBody.M_fixtureList
			switch shape := ft.M_shape.(type) {
			case *box2d.B2PolygonShape: // Box
				// canvas translate
				ctx.Call("translate", curBody.M_xf.P.X, curBody.M_xf.P.Y)
				ctx.Call("rotate", curBody.M_xf.Q.GetAngle())
				ctx.Call("beginPath")
				ctx.Call("moveTo", shape.M_vertices[0].X, shape.M_vertices[0].Y)
				for _, v := range shape.M_vertices[1:shape.M_count] {
					ctx.Call("lineTo", v.X, v.Y)
				}
				ctx.Call("lineTo", shape.M_vertices[0].X, shape.M_vertices[0].Y)
				ctx.Call("fill")
				ctx.Call("stroke")
			case *box2d.B2CircleShape:
				ctx.Call("translate", curBody.M_xf.P.X, curBody.M_xf.P.Y)
				ctx.Call("rotate", curBody.M_xf.Q.GetAngle())
				ctx.Call("beginPath")
				ctx.Call("arc", 0, 0, shape.M_radius, 0, 2*math.Pi)
				ctx.Call("fill")
				ctx.Call("moveTo", 0, 0)
				ctx.Call("lineTo", 0, shape.M_radius)
				ctx.Call("stroke")
			}
			ctx.Call("restore")

		}

		js.Global().Call("requestAnimationFrame", renderFrame)
	})

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done

}

type playerContactListener struct {
}

func (listener playerContactListener) BeginContact(contact box2d.B2ContactInterface) {

	// wait for bodies to actually contact
	if contact.IsTouching() {

		// detect player collision
		if contact.GetFixtureB().GetBody() == player || contact.GetFixtureA().GetBody() == player {

			// check which fixture is the debris
			if contact.GetFixtureA().GetBody() == player {
				weldedDebris = contact.GetFixtureB().GetBody()
			} else {
				weldedDebris = contact.GetFixtureA().GetBody()
			}

			if weldedDebris.GetUserData() == "goalBlock" {
				resetWorld = true
				return
			}

			// If player has already collided with another object this frame
			// ignore this collision
			if !playerCollisionDetected && !playerWelded {
				playerCollisionDetected = true
				weldContact(contact)
			}

		} else if contact.GetFixtureA().GetBody().GetLinearVelocity().Length() > 5 || contact.GetFixtureB().GetBody().GetLinearVelocity().Length() > 5 {
			// detect fast debris
			weldContact(contact)
		}
	}

}

func (listener playerContactListener) EndContact(contact box2d.B2ContactInterface) {

}

func (listener playerContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {

}

func (listener playerContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {

}

func weldContact(contact box2d.B2ContactInterface) {
	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)

	stickyArray = append(stickyArray, StickyInfo{
		bodyA: contact.GetFixtureA().GetBody(),
		bodyB: contact.GetFixtureB().GetBody(),
	})
}

// StickyInfo stores the two objects to be welded together after world.Step()
type StickyInfo struct {
	bodyA *box2d.B2Body
	bodyB *box2d.B2Body
}

func checkPlayerOutOfBounds() {
	if player.GetPosition().X < 0 || player.GetPosition().X > width*worldScale || player.GetPosition().Y < 0 || player.GetPosition().Y > height*worldScale {
		resetWorld = true
	}
}

func clearWorld() {
	// clear out world of any elements
	for joint := world.GetJointList(); joint != nil; joint = joint.GetNext() {
		world.DestroyJoint(joint)
	}

	for body := world.GetBodyList(); body != nil; body = body.GetNext() {
		world.DestroyBody(body)
	}
}

func clearPlayerJoints() {
	for jointEdge := player.GetJointList(); jointEdge != nil; jointEdge = jointEdge.Next {
		world.DestroyJoint(jointEdge.Joint)
	}
}

func populateWorld() {

	// Player Ball
	player = world.CreateBody(&box2d.B2BodyDef{
		Type:         box2d.B2BodyType.B2_dynamicBody,
		Position:     box2d.B2Vec2{X: 20 * worldScale, Y: height*worldScale - 20*worldScale},
		Awake:        true,
		Active:       true,
		GravityScale: 1.0,
		Bullet:       true,
		UserData:     "player",
	})
	shape := box2d.NewB2CircleShape()
	shape.M_radius = 10 * worldScale
	ft := player.CreateFixture(shape, 1)
	ft.M_friction = 0
	ft.M_restitution = 1

	// Create launch block
	launchBlock := world.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_dynamicBody,
		Position: box2d.B2Vec2{X: 20 * worldScale, Y: height*worldScale - 20*worldScale},
		Active:   true,
	})
	launchBlockShape := &box2d.B2PolygonShape{}
	launchBlockShape.SetAsBox(20*worldScale, 20*worldScale)
	ft = launchBlock.CreateFixture(launchBlockShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	// Create goal block
	goalBlock = world.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: width*worldScale - 20*worldScale, Y: 20 * worldScale},
		Active:   true,
		UserData: "goalBlock",
	})
	goalBlockShape := &box2d.B2PolygonShape{}
	goalBlockShape.SetAsBox(20*worldScale, 20*worldScale)
	ft = goalBlock.CreateFixture(goalBlockShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	// find smallest dimension of screen for random object sizing
	var smallestDimension float64

	if width > height {
		smallestDimension = height
	} else {
		smallestDimension = width
	}

	// Some Random debris
	for i := 0; i < 25; i++ {
		obj1 := world.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_dynamicBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * width * worldScale,
				Y: rand.Float64() * height * worldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
		})
		shape := &box2d.B2PolygonShape{}
		shape.SetAsBox(
			rand.Float64()*smallestDimension*worldScale/10,
			rand.Float64()*smallestDimension*worldScale/10)
		ft := obj1.CreateFixture(shape, 1)
		ft.M_friction = 1
		ft.M_restitution = 0 // bouncy
	}
}
