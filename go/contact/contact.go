package contact

import (
	"github.com/ByteArena/box2d"
	"github.com/troyspencer/launch-wasm/go/world"
)

type PlayerContactListener struct {
	*world.WorldState
}

func (listener PlayerContactListener) BeginContact(contact box2d.B2ContactInterface) {
	worldState := listener.WorldState

	// wait for bodies to actually contact
	if contact.IsTouching() {

		// detect player collision
		if contact.GetFixtureB().GetBody().GetUserData() == "player" || contact.GetFixtureA().GetBody().GetUserData() == "player" {

			// check which fixture is the debris
			if contact.GetFixtureA().GetBody().GetUserData() == "player" {
				worldState.WeldedDebris = contact.GetFixtureB().GetBody()
			} else {
				worldState.WeldedDebris = contact.GetFixtureA().GetBody()
			}

			if worldState.WeldedDebris.GetUserData() == "goalBlock" {
				worldState.ResetWorld = true
				return
			}

			// If player has already collided with another object this frame
			// ignore this collision
			if !worldState.PlayerCollisionDetected &&
				!worldState.PlayerWelded &&
				worldState.WeldedDebris.GetUserData() != "bouncyDebris" &&
				worldState.WeldedDebris.GetUserData() != "staticBouncyDebris" {
				worldState.PlayerCollisionDetected = true
				worldState.WeldContact(contact)
			}
		}
	}
}

func (listener PlayerContactListener) EndContact(contact box2d.B2ContactInterface) {}

func (listener PlayerContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
}

func (listener PlayerContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}
