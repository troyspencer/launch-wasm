package world

import (
	"github.com/ByteArena/box2d"
)

func (worldState *WorldState) ClearPlayerJoints() {
	for jointEdge := worldState.Player.GetJointList(); jointEdge != nil; jointEdge = jointEdge.Next {
		worldState.World.DestroyJoint(jointEdge.Joint)
	}
}

func (worldState *WorldState) WeldContact(contact box2d.B2ContactInterface) {
	var worldManifold box2d.B2WorldManifold
	contact.GetWorldManifold(&worldManifold)

	worldState.StickyArray = append(worldState.StickyArray, StickyInfo{
		bodyA: contact.GetFixtureA().GetBody(),
		bodyB: contact.GetFixtureB().GetBody(),
	})
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
