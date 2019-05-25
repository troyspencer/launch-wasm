package contact

import (
	"github.com/ByteArena/box2d"
	"github.com/troyspencer/launch-wasm/go/bodies"
	"github.com/troyspencer/launch-wasm/go/world"
)

type Sticker interface {
	Sticky() bool
}

type PlayerContactListener struct {
	*world.WorldState
}

func (listener PlayerContactListener) BeginContact(contact box2d.B2ContactInterface) {
	worldState := listener.WorldState

	// wait for bodies to actually contact
	if contact.IsTouching() {

		bodyA := contact.GetFixtureA().GetBody()
		bodyB := contact.GetFixtureB().GetBody()
		// check for generally sticky objects
		_, stickyA := bodyA.GetUserData().(Sticker)
		_, stickyB := bodyA.GetUserData().(Sticker)

		if stickyA || stickyB {
			worldState.WeldContact(contact)
		}

		_, playerIsA := bodyA.GetUserData().(*bodies.Player)
		_, playerIsB := bodyB.GetUserData().(*bodies.Player)

		// detect player collision
		if playerIsA || playerIsB {

			if playerIsA {
				worldState.WeldedDebris = bodyB
			} else {
				worldState.WeldedDebris = bodyA
			}

			if _, touchingGoal := worldState.WeldedDebris.GetUserData().(*bodies.GoalBlock); touchingGoal {
				worldState.ResetWorld = true
				return
			}

			_, touchingBouncyDebris := worldState.WeldedDebris.GetUserData().(*bodies.BouncyDebris)

			_, touchingStaticBouncyDebris := worldState.WeldedDebris.GetUserData().(*bodies.StaticBouncyDebris)

			// If player has already collided with another object this frame
			// ignore this collision
			if !worldState.PlayerCollisionDetected &&
				!worldState.PlayerWelded &&
				!touchingBouncyDebris &&
				!touchingStaticBouncyDebris {
				worldState.PlayerCollisionDetected = true
			}
		}
	}
}

func (listener PlayerContactListener) EndContact(contact box2d.B2ContactInterface) {}

func (listener PlayerContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
}

func (listener PlayerContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}
