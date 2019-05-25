package bodies

type JSColors struct {
	fillStyle   string
	strokeStyle string
}

type Body struct {
	JSColors
	sticky bool
	bouncy bool
}

type Bodier interface {
	JSDrawable
	Sticker
	Bouncer
}

type JSDrawable interface {
	FillStyle() string
	StrokeStyle() string
}

type Sticker interface {
	Sticky() bool
}

type Bouncer interface {
	Bouncy() bool
}
