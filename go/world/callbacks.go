package world

import (
	"math"
	"syscall/js"

	"github.com/ByteArena/box2d"
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

		movementDx := mx - worldState.Player.GetPosition().X
		movementDy := my - worldState.Player.GetPosition().Y

		// create normalized movement vector from player to click location
		impulseVelocity := box2d.B2Vec2{X: movementDx, Y: movementDy}
		impulseVelocity.Normalize()
		impulseVelocity.OperatorScalarMulInplace(worldState.GetSmallestDimension() * worldState.WorldScale / 2)

		worldState.ClearPlayerJoints()
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

	worldState.CheckPlayerOutOfBounds()

	if worldState.ResetWorld {
		worldState.Clear()
		worldState.Populate()
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
