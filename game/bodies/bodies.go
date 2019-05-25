package bodies

type JSColors struct {
	fillStyle   string
	strokeStyle string
}

type Body struct {
	JSColors
	sticky bool
	bouncy bool
	breaks bool
}

func (b *Body) FillStyle() string {
	return b.fillStyle
}

func (b *Body) StrokeStyle() string {
	return b.strokeStyle
}

func (b *Body) Sticky() bool {
	return b.sticky
}

func (b *Body) Bouncy() bool {
	return b.bouncy
}

func (b *Body) Breaks() bool {
	return b.breaks
}

type Bodier interface {
	JSDrawable
	Sticker
	Bouncer
	Breaker
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

type Breaker interface {
	Breaks() bool
}
