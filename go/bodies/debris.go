package bodies

type Debris struct {
	jsColors JSColors
}

func NewDebris() *Debris {
	return &Debris{
		jsColors: JSColors{
			FillStyle:   "rgba(100,100,100,1)",
			StrokeStyle: "rgba(100,100,100,1)",
		},
	}
}

func (p *Debris) FillStyle() string {
	return p.jsColors.FillStyle
}

func (p *Debris) StrokeStyle() string {
	return p.jsColors.StrokeStyle
}
