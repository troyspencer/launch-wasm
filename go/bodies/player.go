package bodies

type Player struct {
	Body
}

func NewPlayer() *Player {
	return &Player{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(180, 180,180,1)",
				strokeStyle: "rgba(180, 180,180,1)",
			},
			sticky: true,
			bouncy: false,
		},
	}
}

func (b *Player) FillStyle() string {
	return b.fillStyle
}

func (b *Player) StrokeStyle() string {
	return b.strokeStyle
}

func (b *Player) Sticky() bool {
	return b.sticky
}

func (b *Player) Bouncy() bool {
	return b.bouncy
}
