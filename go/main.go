// Click on canvas to start a polygon
// - Max 8 vertices
// - Only convex polygons
// - Esc cancel polygon

//Wasming
// compile: GOOS=js GOARCH=wasm go build -o main.wasm ./main.go
package main

import (
	"log"
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
	begin                   bool    = true
	worldScale                      = 0.0125 // 1/8
	player                  *box2d.B2Body
	playerJoint             box2d.B2JointInterface
	playerCollisionDetected bool
	playerWelded            bool
	stickyArray             []StickyInfo
	defaultWorld            = box2d.B2World{}
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

	world := box2d.MakeB2World(box2d.B2Vec2{X: 0, Y: 0})

	world.SetContactListener(&playerContactListener{})

	// Player Ball
	player = world.CreateBody(&box2d.B2BodyDef{
		Type:         box2d.B2BodyType.B2_dynamicBody,
		Position:     box2d.B2Vec2{X: 0.2 * width * worldScale, Y: 0.9 * height * worldScale},
		Awake:        true,
		Active:       true,
		GravityScale: 1.0,
		Bullet:       true,
	})
	shape := box2d.NewB2CircleShape()
	shape.M_radius = 15 * worldScale
	ft := player.CreateFixture(shape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	// Boundaries
	floor := world.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: 0, Y: height*worldScale - 20*worldScale},
		Active:   true,
	})
	floorShape := &box2d.B2PolygonShape{}
	floorShape.SetAsBox(width*worldScale, 20*worldScale)
	ft = floor.CreateFixture(floorShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	ceiling := world.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: 0, Y: 20 * worldScale},
		Active:   true,
	})
	ceilingShape := &box2d.B2PolygonShape{}
	ceilingShape.SetAsBox(width*worldScale, 20*worldScale)
	ft = ceiling.CreateFixture(ceilingShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	leftWall := world.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: 20 * worldScale, Y: 0},
		Active:   true,
	})
	leftWallShape := &box2d.B2PolygonShape{}
	leftWallShape.SetAsBox(20*worldScale, height*worldScale)
	ft = leftWall.CreateFixture(leftWallShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	rightWall := world.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: width*worldScale - 20*worldScale, Y: 0},
		Active:   true,
	})
	rightWallShape := &box2d.B2PolygonShape{}
	rightWallShape.SetAsBox(20*worldScale, height*worldScale)
	ft = rightWall.CreateFixture(rightWallShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0

	// Some Random debris
	for i := 0; i < 3; i++ {
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
			(60*rand.Float64()+10)*worldScale,
			(60*rand.Float64()+20)*worldScale)
		ft := obj1.CreateFixture(shape, 1)
		ft.M_friction = 1
		ft.M_restitution = 0 // bouncy
	}

	defaultWorld = world

	mouseDownEvt := js.NewCallback(func(args []js.Value) {

		e := args[0]
		if e.Get("target") != canvasEl {
			return
		}
		mx := e.Get("clientX").Float() * worldScale
		my := e.Get("clientY").Float() * worldScale

		movementDx := mx - player.GetPosition().X
		movementDy := my - player.GetPosition().Y

		// create normalized movement vector from player to click location
		movementVector := box2d.B2Vec2{X: movementDx, Y: movementDy}
		movementVector.Normalize()
		movementVector.OperatorScalarMulInplace(5)

		if playerWelded {
			world.DestroyJoint(playerJoint)
			playerWelded = false
			player.SetLinearVelocity(movementVector)
		}

		if begin {
			begin = false
			player.SetLinearVelocity(movementVector)
		}

	})
	defer mouseDownEvt.Release()

	keyUpEvt := js.NewCallback(func(args []js.Value) {
		e := args[0]
		if e.Get("which").Int() == 27 {
			log.Println("Reset")
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

		// check for new weld joint and execute it
		for len(stickyArray) > 0 {
			stickyBody := stickyArray[0]
			stickyArray[0] = stickyArray[len(stickyArray)-1]
			stickyArray = stickyArray[:len(stickyArray)-1]

			var worldCoordsAnchorPoint box2d.B2Vec2
			worldCoordsAnchorPoint = stickyBody.bodyB.GetWorldPoint(box2d.B2Vec2{0.6, 0})

			weldJointDef := box2d.MakeB2WeldJointDef()
			weldJointDef.BodyA = stickyBody.bodyA
			weldJointDef.BodyB = stickyBody.bodyB
			weldJointDef.ReferenceAngle = weldJointDef.BodyB.GetAngle() - weldJointDef.BodyA.GetAngle()
			weldJointDef.LocalAnchorA = weldJointDef.BodyA.GetLocalPoint(worldCoordsAnchorPoint)
			weldJointDef.LocalAnchorB = weldJointDef.BodyB.GetLocalPoint(worldCoordsAnchorPoint)
			if playerWelded {
				world.DestroyJoint(playerJoint)
			}
			playerJoint = world.CreateJoint(&weldJointDef)
			playerCollisionDetected = false
			playerWelded = true
		}

		ctx.Call("clearRect", 0, 0, width*worldScale, height*worldScale)

		// color for other objects
		ctx.Set("fillStyle", "rgba(100,100,100,1)")
		ctx.Set("strokeStyle", "rgba(100,100,100,1)")

		for curBody := world.GetBodyList(); curBody != nil; curBody = curBody.M_next {
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

		// Player ball
		// color
		ctx.Set("fillStyle", "rgba(180, 180,180,1)")
		ctx.Set("strokeStyle", "rgba(180,180,180,1)")

		// draw player
		curBody := player
		ctx.Call("save")
		ctx.Call("translate", curBody.M_xf.P.X, curBody.M_xf.P.Y)
		ctx.Call("rotate", curBody.M_xf.Q.GetAngle())
		ctx.Call("beginPath")
		ctx.Call("arc", 0, 0, shape.M_radius, 0, 2*math.Pi)
		ctx.Call("fill")
		ctx.Call("moveTo", 0, 0)
		ctx.Call("lineTo", 0, shape.M_radius)
		ctx.Call("stroke")

		ctx.Call("restore")

		js.Global().Call("requestAnimationFrame", renderFrame)
	})

	// Start running
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-done

}

type playerContactListener struct {
}

func (listener playerContactListener) BeginContact(contact box2d.B2ContactInterface) {
	if contact.GetFixtureB().GetBody() == player || contact.GetFixtureA().GetBody() == player {
		if contact.IsTouching() && !playerCollisionDetected {
			playerCollisionDetected = true

			//contactPoint := contact.GetManifold().Points[0].Id
			var worldManifold box2d.B2WorldManifold
			contact.GetWorldManifold(&worldManifold)
			log.Println("World", worldManifold.Points[0])
			log.Println("Local", contact.GetManifold().LocalPoint)

			stickyArray = append(stickyArray, StickyInfo{
				bodyA:       contact.GetFixtureA().GetBody(),
				bodyB:       contact.GetFixtureB().GetBody(),
				anchorPoint: contact.GetManifold().LocalPoint,
			})
		}

	}
}

func (listener playerContactListener) EndContact(contact box2d.B2ContactInterface) {

}

func (listener playerContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {

}

func (listener playerContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {

}

// StickyInfo stores the two objects to be welded together after world.Step()
type StickyInfo struct {
	bodyA       *box2d.B2Body
	bodyB       *box2d.B2Body
	anchorPoint box2d.B2Vec2
}
