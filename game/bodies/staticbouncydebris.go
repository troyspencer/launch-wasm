package bodies

type StaticBouncyDebris struct {
	Body
}

func NewStaticBouncyDebris() *StaticBouncyDebris {
	return &StaticBouncyDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(100,0,0,1)",
				strokeStyle: "rgba(100,0,0,1)",
			},
			sticky: false,
			bouncy: true,
		},
	}
}

func (b *StaticBouncyDebris) FillStyle() string {
	return b.fillStyle
}

func (b *StaticBouncyDebris) StrokeStyle() string {
	return b.strokeStyle
}

func (b *StaticBouncyDebris) Sticky() bool {
	return b.sticky
}

func (b *StaticBouncyDebris) Bouncy() bool {
	return b.bouncy
}
