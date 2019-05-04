package world

import (
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
