package world

import (
	"log"
	"math"

	"github.com/ByteArena/box2d"
)

type Breaker interface {
	Breaks() bool
}

func (worldState *WorldState) LaunchPlayer(mx float64, my float64) {
	movementDx := mx - worldState.Player.GetPosition().X
	movementDy := my - worldState.Player.GetPosition().Y

	// create normalized movement vector from player to click location
	impulseVelocity := box2d.B2Vec2{X: movementDx, Y: movementDy}
	impulseVelocity.Normalize()
	impulseVelocity.OperatorScalarMulInplace(worldState.GetSmallestDimension() * worldState.WorldScale / 2)
	if worldState.AbsorbCount > 0 {
		impulseVelocity.OperatorScalarMulInplace(0.5)
	}

	worldState.ClearJoints(worldState.Player)
	worldState.ClearJoints(worldState.WeldedDebris)
	worldState.PlayerWelded = false

	worldState.PushDebris(impulseVelocity)

	// set player velocity to player desired velocity
	worldState.Player.SetLinearVelocity(impulseVelocity)
}

func (worldState *WorldState) PushDebris(impulseVelocity box2d.B2Vec2) {
	if worldState.WeldedDebris.GetType() != box2d.B2BodyType.B2_dynamicBody {
		return
	}

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

	if breaker, ok := worldState.WeldedDebris.GetUserData().(Breaker); ok {
		if breaker.Breaks() {
			worldState.Slice(worldState.WeldedDebris, debrisVelocity)
		}
	} else {
		// update debris velocity
		worldState.WeldedDebris.SetLinearVelocity(debrisVelocity)
	}

}

func (worldState *WorldState) breakerJumpedOff(fixture *box2d.B2Fixture, point box2d.B2Vec2, normal box2d.B2Vec2, fraction float64) float64 {
	affectedBody := fixture.GetBody()
	affectedPolygon, ok := fixture.GetShape().(*box2d.B2PolygonShape)
	if !ok {
		return 0
	}

	if worldState.BreaksInfo.AffectedByLaunch != affectedBody {
		worldState.BreaksInfo.AffectedByLaunch = affectedBody
		worldState.BreaksInfo.EntryPoint = point
	} else {
		entryPoint := worldState.BreaksInfo.EntryPoint
		rayCenter := box2d.B2Vec2{
			X: (point.X + entryPoint.X) / 2,
			Y: (point.Y + entryPoint.Y) / 2,
		}
		rayAngle := math.Atan2(entryPoint.Y-point.Y, entryPoint.X-point.X)

		polyVertices := affectedPolygon.M_vertices
		for _, polyVertex := range polyVertices {
			worldPoint := affectedBody.GetWorldPoint(polyVertex)
			cutAngle := math.Atan2(worldPoint.Y-rayCenter.Y, worldPoint.X-rayCenter.X) - rayAngle
			if cutAngle < math.Pi*-1 {
				cutAngle += 2 * math.Pi
			}
			if cutAngle > 0 && cutAngle <= math.Pi {
				log.Println("above")
			} else {
				log.Println("below")
			}
		}
		log.Println(entryPoint, rayCenter, point)
	}
	return 1
}

func (worldState *WorldState) Slice(body *box2d.B2Body, impulse box2d.B2Vec2) {
	sliceBeginning := worldState.Player.GetWorldCenter()
	sliceEnd := box2d.B2Vec2{
		X: sliceBeginning.X + impulse.X,
		Y: sliceBeginning.Y + impulse.Y,
	}
	worldState.World.RayCast(worldState.breakerJumpedOff, sliceBeginning, sliceEnd)
	worldState.World.RayCast(worldState.breakerJumpedOff, sliceEnd, sliceBeginning)
}
