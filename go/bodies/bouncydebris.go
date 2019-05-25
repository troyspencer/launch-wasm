package bodies

type BouncyDebris struct {
	Body
}

func NewBouncyDebris() *BouncyDebris {
	return &BouncyDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(200,0,0,1)",
				strokeStyle: "rgba(200,0,0,1)",
			},
			sticky: false,
		},
	}
}

func (b *BouncyDebris) FillStyle() string {
	return b.fillStyle
}

func (b *BouncyDebris) StrokeStyle() string {
	return b.strokeStyle
}

func (b *BouncyDebris) Sticky() bool {
	return b.sticky
}
