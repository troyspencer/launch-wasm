package bodies

type StaticDebris struct {
	Body
}

func NewStaticDebris() *StaticDebris {
	return &StaticDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(50,50,50,1)",
				strokeStyle: "rgba(50,50,50,1)",
			},
			sticky: false,
			bouncy: false,
		},
	}
}

func (b *StaticDebris) FillStyle() string {
	return b.fillStyle
}

func (b *StaticDebris) StrokeStyle() string {
	return b.strokeStyle
}

func (b *StaticDebris) Sticky() bool {
	return b.sticky
}

func (b *StaticDebris) Bouncy() bool {
	return b.bouncy
}
