package world

import (
	"math/rand"

	"github.com/ByteArena/box2d"
	"github.com/troyspencer/launch-wasm/go/bodies"
)

func (worldState *WorldState) Populate() {
	worldState.CreateLaunchBlock()
	worldState.CreatePlayer()
	worldState.CreateGoalBlock()
	worldState.CreateDebris()
	worldState.CreateStaticDebris()
	worldState.CreateBouncyDebris()
	worldState.CreateStaticBouncyDebris()
	worldState.CreateStickyDebris()
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

func (worldState *WorldState) CreateStickyDebris() {
	smallestDimension := worldState.GetSmallestDimension()

	// Some Random sticky debris
	for i := 0; i < 5; i++ {
		newStickyDebris := bodies.NewStickyDebris()
		obj1 := worldState.World.CreateBody(&box2d.B2BodyDef{
			Type: box2d.B2BodyType.B2_dynamicBody,
			Position: box2d.B2Vec2{
				X: rand.Float64() * worldState.Width * worldState.WorldScale,
				Y: rand.Float64() * worldState.Height * worldState.WorldScale},
			Angle:        rand.Float64() * 100,
			Awake:        true,
			Active:       true,
			GravityScale: 1.0,
			UserData:     newStickyDebris,
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
