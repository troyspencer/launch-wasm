package bodies

type StickyDebris struct {
	Body
}

func NewStickyDebris() *StickyDebris {
	return &StickyDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgb(252, 241, 37)",
				strokeStyle: "rgb(252, 241, 37)",
			},
			sticky: true,
			bouncy: false,
		},
	}
}

func (b *StickyDebris) FillStyle() string {
	return b.fillStyle
}

func (b *StickyDebris) StrokeStyle() string {
	return b.strokeStyle
}

func (b *StickyDebris) Sticky() bool {
	return b.sticky
}

func (b *StickyDebris) Bouncy() bool {
	return b.bouncy
}
