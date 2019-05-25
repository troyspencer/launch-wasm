package bodies

type BouncyDebris struct {
	jsColors JSColors
}

func NewBouncyDebris() *BouncyDebris {
	return &BouncyDebris{
		jsColors: JSColors{
			FillStyle:   "rgba(200,0,0,1)",
			StrokeStyle: "rgba(200,0,0,1)",
		},
	}
}

func (p *BouncyDebris) FillStyle() string {
	return p.jsColors.FillStyle
}

func (p *BouncyDebris) StrokeStyle() string {
	return p.jsColors.StrokeStyle
}
