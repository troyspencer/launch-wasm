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

type JSDrawable interface {
	FillStyle() string
	StrokeStyle() string
}

func (jso *JSObjects) Draw(body *box2d.B2Body) {
	if jsDrawable, ok := body.GetUserData().(JSDrawable); ok {
		jso.Context.Set("fillStyle", jsDrawable.FillStyle())
		jso.Context.Set("strokeStyle", jsDrawable.StrokeStyle())
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
