package world

import (
	"log"
	"math"

	"github.com/ByteArena/box2d"
	"github.com/troyspencer/launch-wasm/game/bodies"
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

func (worldState *WorldState) Slice(body *box2d.B2Body, impulse box2d.B2Vec2) {
	sliceBeginning := worldState.Player.GetPosition()
	sliceEnd := box2d.B2Vec2{
		X: sliceBeginning.X + 100*impulse.X,
		Y: sliceBeginning.Y + 100*impulse.Y,
	}

	log.Println(sliceBeginning, impulse, sliceEnd)
	worldState.World.RayCast(worldState.breakerJumpedOff, sliceBeginning, sliceEnd)
	worldState.World.RayCast(worldState.breakerJumpedOff, sliceEnd, sliceBeginning)
}

func (worldState *WorldState) breakerJumpedOff(fixture *box2d.B2Fixture, point box2d.B2Vec2, normal box2d.B2Vec2, fraction float64) float64 {

	affectedBody := fixture.GetBody()
	if _, isPlayer := affectedBody.GetUserData().(*bodies.Player); isPlayer {
		return 0
	}
	affectedPolygon, ok := fixture.GetShape().(*box2d.B2PolygonShape)
	if !ok {
		return 0
	}

	if worldState.BreaksInfo.AffectedByLaunch != nil {
		log.Println(worldState.BreaksInfo.AffectedByLaunch.GetUserData())
	}
	if worldState.BreaksInfo.AffectedByLaunch != affectedBody {
		log.Println("Start")
		worldState.BreaksInfo.AffectedByLaunch = affectedBody
		worldState.BreaksInfo.EntryPoint = point
	} else {
		log.Println("End")
		entryPoint := worldState.BreaksInfo.EntryPoint
		rayCenter := box2d.B2Vec2{
			X: (point.X + entryPoint.X) / 2,
			Y: (point.Y + entryPoint.Y) / 2,
		}
		rayAngle := math.Atan2(entryPoint.Y-point.Y, entryPoint.X-point.X)

		polyVertices := affectedPolygon.M_vertices
		newPolyVertices1 := []*box2d.B2Vec2{}
		newPolyVertices2 := []*box2d.B2Vec2{}
		currentPoly := 0
		cutPlaced1 := false
		cutPlaced2 := false
		for _, polyVertex := range polyVertices {
			worldPoint := affectedBody.GetWorldPoint(polyVertex)
			cutAngle := math.Atan2(worldPoint.Y-rayCenter.Y, worldPoint.X-rayCenter.X) - rayAngle
			if cutAngle < math.Pi*-1 {
				cutAngle += 2 * math.Pi
			}
			if cutAngle > 0 && cutAngle <= math.Pi {
				if currentPoly == 2 {
					cutPlaced1 = true
					newPolyVertices1 = append(newPolyVertices1, &point)
					newPolyVertices1 = append(newPolyVertices1, &entryPoint)
				}
				newPolyVertices1 = append(newPolyVertices1, &worldPoint)
				currentPoly = 1
			} else {
				if currentPoly == 1 {
					cutPlaced2 = true
					newPolyVertices2 = append(newPolyVertices2, &entryPoint)
					newPolyVertices2 = append(newPolyVertices2, &point)
				}
				newPolyVertices2 = append(newPolyVertices2, &worldPoint)
				currentPoly = 2
			}
		}

		if !cutPlaced1 {
			newPolyVertices1 = append(newPolyVertices1, &point)
			newPolyVertices1 = append(newPolyVertices1, &entryPoint)
		}
		if !cutPlaced2 {
			newPolyVertices2 = append(newPolyVertices2, &entryPoint)
			newPolyVertices2 = append(newPolyVertices2, &point)
		}
		log.Println("Affected body:", affectedBody.GetUserData())
		worldState.createSlice(newPolyVertices1, len(newPolyVertices1))
		worldState.createSlice(newPolyVertices2, len(newPolyVertices2))

		worldState.World.DestroyBody(affectedBody)
	}
	return 1
}

func (worldState *WorldState) createSlice(vertices []*box2d.B2Vec2, numVertices int) {
	center := findCentroid(vertices)
	for _, vertex := range vertices {
		vertex.OperatorMinusInplace(center)
	}
	verticesDereference := []box2d.B2Vec2{}
	for _, vertex := range vertices {
		verticesDereference = append(verticesDereference, *vertex)
	}
	log.Println(verticesDereference)
	sliceBody := &box2d.B2BodyDef{}
	sliceBody.Position.Set(center.X, center.Y)
	sliceBody.Type = box2d.B2BodyType.B2_dynamicBody

	slicePoly := &box2d.B2PolygonShape{}

	slicePoly.Set(verticesDereference, len(vertices))
	worldSlice := worldState.World.CreateBody(sliceBody)
	worldSlice.CreateFixture(slicePoly, 1)
	for _, vertex := range vertices {
		vertex.OperatorPlusInplace(center)
	}
}

func findCentroid(vs []*box2d.B2Vec2) box2d.B2Vec2 {
	c := box2d.B2Vec2{}
	area := 0.0
	p1X := 0.0
	p1Y := 0.0
	inv3 := 1.0 / 3.0
	for i := 1; i < len(vs); i++ {
		p2 := vs[i]
		p3 := &box2d.B2Vec2{}
		if i+1 < len(vs) {
			p3 = vs[i+1]
		} else {
			p3 = vs[0]
		}
		e1X := p2.X - p1X
		e1Y := p2.Y - p1Y
		e2X := p3.X - p1X
		e2Y := p3.Y - p1Y
		D := e1X*e2Y - e1Y*e2X
		triangleArea := 0.5 * D
		area += triangleArea
		c.X += triangleArea * inv3 * (p1X + p2.X + p3.X)
		c.Y += triangleArea * inv3 * (p1Y + p2.Y + p3.Y)
	}
	c.X *= 1.0 / area
	c.Y *= 1.0 / area
	return c
}
