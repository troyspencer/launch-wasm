package bodies

type Player struct {
	jsColors JSColors
}

func NewPlayer() *Player {
	return &Player{
		jsColors: JSColors{
			FillStyle:   "rgba(180, 180,180,1)",
			StrokeStyle: "rgba(180, 180,180,1)",
		},
	}
}

func (p *Player) FillStyle() string {
	return p.jsColors.FillStyle
}

func (p *Player) StrokeStyle() string {
	return p.jsColors.StrokeStyle
}
