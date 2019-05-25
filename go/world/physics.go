package world

import "github.com/ByteArena/box2d"

func (worldState *WorldState) PushDebris(impulseVelocity box2d.B2Vec2) {
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
}

func (worldState *WorldState) LaunchPlayer(mx float64, my float64) {
	movementDx := mx - worldState.Player.GetPosition().X
	movementDy := my - worldState.Player.GetPosition().Y

	// create normalized movement vector from player to click location
	impulseVelocity := box2d.B2Vec2{X: movementDx, Y: movementDy}
	impulseVelocity.Normalize()
	impulseVelocity.OperatorScalarMulInplace(worldState.GetSmallestDimension() * worldState.WorldScale / 2)

	worldState.ClearPlayerJoints()
	worldState.PlayerWelded = false

	worldState.PushDebris(impulseVelocity)

	// set player velocity to player desired velocity
	worldState.Player.SetLinearVelocity(impulseVelocity)
}
