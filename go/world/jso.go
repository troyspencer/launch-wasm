package world

import (
	"math"
	"syscall/js"

	"github.com/ByteArena/box2d"
)

type JSObjects struct {
	Context js.Value
	Doc     js.Value
	Canvas  js.Value
}

func (jso *JSObjects) Draw(body *box2d.B2Body) {
	if body.GetUserData() == "player" {
		// Player ball color
		jso.Context.Set("fillStyle", "rgba(180, 180,180,1)")
		jso.Context.Set("strokeStyle", "rgba(180,180,180,1)")
	} else if body.GetUserData() == "goalBlock" {
		// Goal block color
		jso.Context.Set("fillStyle", "rgba(0, 255,0,1)")
		jso.Context.Set("strokeStyle", "rgba(0,255,0,1)")
	} else if body.GetUserData() == "staticDebris" || body.GetUserData() == "launchBlock" {
		// color for other objects
		jso.Context.Set("fillStyle", "rgba(50,50,50,1)")
		jso.Context.Set("strokeStyle", "rgba(50,50,50,1)")
	} else {
		// color for other objects
		jso.Context.Set("fillStyle", "rgba(100,100,100,1)")
		jso.Context.Set("strokeStyle", "rgba(100,100,100,1)")
	}
	// Only one fixture for now
	jso.Context.Call("save")
	ft := body.M_fixtureList
	switch shape := ft.M_shape.(type) {
	case *box2d.B2PolygonShape: // Box
		// canvas translate
		jso.Context.Call("translate", body.M_xf.P.X, body.M_xf.P.Y)
		jso.Context.Call("rotate", body.M_xf.Q.GetAngle())
		jso.Context.Call("beginPath")
		jso.Context.Call("moveTo", shape.M_vertices[0].X, shape.M_vertices[0].Y)
		for _, v := range shape.M_vertices[1:shape.M_count] {
			jso.Context.Call("lineTo", v.X, v.Y)
		}
		jso.Context.Call("lineTo", shape.M_vertices[0].X, shape.M_vertices[0].Y)
		jso.Context.Call("fill")
		jso.Context.Call("stroke")
	case *box2d.B2CircleShape:
		jso.Context.Call("translate", body.M_xf.P.X, body.M_xf.P.Y)
		jso.Context.Call("rotate", body.M_xf.Q.GetAngle())
		jso.Context.Call("beginPath")
		jso.Context.Call("arc", 0, 0, shape.M_radius, 0, 2*math.Pi)
		jso.Context.Call("fill")
		jso.Context.Call("moveTo", 0, 0)
		jso.Context.Call("lineTo", 0, shape.M_radius)
		jso.Context.Call("stroke")
	}
	jso.Context.Call("restore")
}
