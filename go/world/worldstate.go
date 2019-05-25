package world

import (
	"math/rand"
	"syscall/js"

	"github.com/ByteArena/box2d"
	"github.com/troyspencer/launch-wasm/go/bodies"
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

func (worldState *WorldState) CreateLaunchBlock() {
	smallestDimension := worldState.GetSmallestDimension()
	newLaunchBlock := bodies.NewLaunchBlock()
	launchBlock := worldState.World.CreateBody(&box2d.B2BodyDef{
		Type: box2d.B2BodyType.B2_staticBody,
		Position: box2d.B2Vec2{
			X: smallestDimension * worldState.WorldScale / 32,
			Y: worldState.Height*worldState.WorldScale - smallestDimension*worldState.WorldScale/32,
		},
		Active:   true,
		UserData: newLaunchBlock,
	})
	launchBlockShape := &box2d.B2PolygonShape{}
	launchBlockShape.SetAsBox(
		smallestDimension*worldState.WorldScale/32,
		smallestDimension*worldState.WorldScale/32,
	)
	ft := launchBlock.CreateFixture(launchBlockShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0
}

func (worldState *WorldState) CreatePlayer() {
	smallestDimension := worldState.GetSmallestDimension()
	newPlayer := bodies.NewPlayer()
	worldState.Player = worldState.World.CreateBody(&box2d.B2BodyDef{
		Type:         box2d.B2BodyType.B2_dynamicBody,
		Position:     box2d.B2Vec2{X: smallestDimension * worldState.WorldScale / 32, Y: worldState.Height*worldState.WorldScale - smallestDimension*worldState.WorldScale/32},
		Awake:        true,
		Active:       true,
		GravityScale: 1.0,
		Bullet:       true,
		UserData:     newPlayer,

	})
	shape := box2d.NewB2CircleShape()
	shape.M_radius = smallestDimension * worldState.WorldScale / 64
	ft := worldState.Player.CreateFixture(shape, 1)
	ft.M_friction = 0
	ft.M_restitution = 1
}

func (worldState *WorldState) CreateGoalBlock() {
	smallestDimension := worldState.GetSmallestDimension()
	
	goalBlock := bodies.NewGoalBlock()
	worldState.GoalBlock = worldState.World.CreateBody(&box2d.B2BodyDef{
		Type:     box2d.B2BodyType.B2_kinematicBody,
		Position: box2d.B2Vec2{X: worldState.Width*worldState.WorldScale - smallestDimension*worldState.WorldScale/32, Y: smallestDimension * worldState.WorldScale / 32},
		Active:   true,
		UserData: goalBlock,
	})
	goalBlockShape := &box2d.B2PolygonShape{}
	goalBlockShape.SetAsBox(smallestDimension*worldState.WorldScale/32, smallestDimension*worldState.WorldScale/32)
	ft := worldState.GoalBlock.CreateFixture(goalBlockShape, 1)
	ft.M_friction = 1
	ft.M_restitution = 0
}

func (worldState *WorldState) CreateDebris() {
	smallestDimension := worldState.GetSmallestDimension()

	// Some Random debris
	for i := 0; i < 10; i++ {
		newDebris := bodies.NewDebris()
		obj1 := worldState.World.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_dynamicBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * worldState.Width * worldState.WorldScale,
				Y: rand.Float64() * worldState.Height * worldState.WorldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
			UserData:     newDebris,
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

func (worldState *WorldState) CreateStaticDebris() {
	smallestDimension := worldState.GetSmallestDimension()

	// Some Random debris
	for i := 0; i < 4; i++ {
		newStaticDebris := bodies.NewStaticDebris()
		obj1 := worldState.World.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_staticBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * worldState.Width * worldState.WorldScale,
				Y: rand.Float64() * worldState.Height * worldState.WorldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
			UserData:     newStaticDebris,
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

func (worldState *WorldState) CreateStaticBouncyDebris() {
	smallestDimension := worldState.GetSmallestDimension()

	// Some Random debris
	for i := 0; i < 3; i++ {
		staticBouncyDebris := bodies.NewStaticBouncyDebris()
		obj1 := worldState.World.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_staticBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * worldState.Width * worldState.WorldScale,
				Y: rand.Float64() * worldState.Height * worldState.WorldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
			UserData:     staticBouncyDebris,
		})
		shape := &box2d.B2PolygonShape{}
		shape.SetAsBox(
			rand.Float64()*smallestDimension*worldState.WorldScale/10,
			rand.Float64()*smallestDimension*worldState.WorldScale/10)
		ft := obj1.CreateFixture(shape, 1)
		ft.M_friction = 1
		ft.M_restitution = 1 // bouncy
	}
}

func (worldState *WorldState) CreateBouncyDebris() {
	smallestDimension := worldState.GetSmallestDimension()

	// Some Random debris
	for i := 0; i < 3; i++ {
		bouncyDebris := bodies.NewBouncyDebris()
		obj1 := worldState.World.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_dynamicBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * worldState.Width * worldState.WorldScale,
				Y: rand.Float64() * worldState.Height * worldState.WorldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
			UserData:     bouncyDebris,
		})
		shape := &box2d.B2PolygonShape{}
		shape.SetAsBox(
			rand.Float64()*smallestDimension*worldState.WorldScale/10,
			rand.Float64()*smallestDimension*worldState.WorldScale/10)
		ft := obj1.CreateFixture(shape, 1)
		ft.M_friction = 1
		ft.M_restitution = 1 // bouncy
	}
}

func (worldState *WorldState) Populate() {
	worldState.CreateLaunchBlock()
	worldState.CreatePlayer()
	worldState.CreateGoalBlock()
	worldState.CreateDebris()
	worldState.CreateStaticDebris()
	worldState.CreateBouncyDebris()
	worldState.CreateStaticBouncyDebris()
}

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

func (worldState *WorldState) WeldContact(contact box2d.B2ContactInterface) {
	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)

	worldState.StickyArray = append(worldState.StickyArray, StickyInfo{
		bodyA: contact.GetFixtureA().GetBody(),
		bodyB: contact.GetFixtureB().GetBody(),
	})
}

func (worldState *WorldState) Reset() {
	worldState.Clear()
	worldState.Populate()
	worldState.ResetWorld = false
}

func (worldState *WorldState) WeldJoint() {
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
}
