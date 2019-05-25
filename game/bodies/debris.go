package bodies

type Debris struct {
	Body
}

func NewDebris() *Debris {
	return &Debris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(100,100,100,1)",
				strokeStyle: "rgba(100,100,100,1)",
			},
			sticky: false,
			bouncy: false,
		},
	}
}

func (b *Debris) FillStyle() string {
	return b.fillStyle
}

func (b *Debris) StrokeStyle() string {
	return b.strokeStyle
}

func (b *Debris) Sticky() bool {
	return b.sticky
}

func (b *Debris) Bouncy() bool {
	return b.bouncy
}
