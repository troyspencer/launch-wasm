package contact

import (
	"github.com/ByteArena/box2d"
	"github.com/troyspencer/launch-wasm/game/bodies"
	"github.com/troyspencer/launch-wasm/game/world"
)

type Sticker interface {
	Sticky() bool
}

type Bouncer interface {
	Bouncy() bool
}

type Absorber interface {
	Absorbs() bool
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

		fixtureBodies := []*box2d.B2Body{bodyA, bodyB}

		// check for bounce first, nothing welds to bouncy
		bouncy := false

		for _, fixtureBody := range fixtureBodies {
			if bouncer, ok := fixtureBody.GetUserData().(Bouncer); ok {
				if bouncer.Bouncy() {
					bouncy = true
				}
			}
		}

		if bouncy {
			return
		}

		// check for sticky
		sticky := false
		for _, fixtureBody := range fixtureBodies {
			if sticker, ok := fixtureBody.GetUserData().(Sticker); ok {
				if sticker.Sticky() {
					sticky = true
				}
			}
		}

		if !sticky {
			return
		}

		// check for player
		_, playerIsA := bodyA.GetUserData().(*bodies.Player)
		_, playerIsB := bodyB.GetUserData().(*bodies.Player)

		playerContact := playerIsA || playerIsB

		// check for absorbs
		absorbs := false
		for _, fixtureBody := range fixtureBodies {
			if absorber, ok := fixtureBody.GetUserData().(Absorber); ok {
				if absorber.Absorbs() {
					absorbs = true
				}
			}
		}

		if absorbs {
			for _, fixtureBody := range fixtureBodies {
				fixtureBody.SetLinearDamping(0.4)
			}
			if playerContact {
				if worldState.AbsorbCount == 0 {
					// player entering absorber
					playerVelocity := worldState.Player.GetLinearVelocity()
					playerVelocity.OperatorScalarMulInplace(0.5)
					worldState.Player.SetLinearVelocity(playerVelocity)
				}
				worldState.AbsorbCount++
				return
			}
			return
		}

		// Prevent a welded player from welding again
		if playerContact && worldState.PlayerWelded && sticky {
			return
		}

		worldState.WeldContact(contact)

		// detect player collision
		if playerContact {
			if playerIsA {
				worldState.WeldedDebris = bodyB
			} else {
				worldState.WeldedDebris = bodyA
			}

			if _, touchingGoal := worldState.WeldedDebris.GetUserData().(*bodies.GoalBlock); touchingGoal {
				worldState.ResetWorld = true
				return
			}

			// If player has already collided with another object this frame
			// ignore this collision
			if !worldState.PlayerCollisionDetected &&
				!worldState.PlayerWelded {
				worldState.PlayerCollisionDetected = true
			}

		}

	}
}

func (listener PlayerContactListener) EndContact(contact box2d.B2ContactInterface) {
	worldState := listener.WorldState

	if !contact.IsTouching() {
		bodyA := contact.GetFixtureA().GetBody()
		bodyB := contact.GetFixtureB().GetBody()

		fixtureBodies := []*box2d.B2Body{bodyA, bodyB}

		// check for player
		_, playerIsA := bodyA.GetUserData().(*bodies.Player)
		_, playerIsB := bodyB.GetUserData().(*bodies.Player)

		playerContact := playerIsA || playerIsB

		// check for absorbs
		absorbs := false
		for _, fixtureBody := range fixtureBodies {
			if absorber, ok := fixtureBody.GetUserData().(Absorber); ok {
				if absorber.Absorbs() {
					absorbs = true
				}
			}
		}

		if absorbs && playerContact {
			worldState.AbsorbCount--
			if worldState.AbsorbCount == 0 {
				// player exiting all absorbers
				worldState.Player.SetLinearDamping(0)
				playerVelocity := worldState.Player.GetLinearVelocity()
				playerVelocity.OperatorScalarMulInplace(2)
				worldState.Player.SetLinearVelocity(playerVelocity)
			}
		}
	}
}

func (listener PlayerContactListener) PreSolve(contact box2d.B2ContactInterface, oldManifold box2d.B2Manifold) {
}

func (listener PlayerContactListener) PostSolve(contact box2d.B2ContactInterface, impulse *box2d.B2ContactImpulse) {
}
