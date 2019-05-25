package bodies

type StaticBouncyDebris struct {
	jsColors JSColors
}

func NewStaticBouncyDebris() *StaticBouncyDebris {
	return &StaticBouncyDebris{
		jsColors: JSColors{
			FillStyle:   "rgba(100,0,0,1)",
			StrokeStyle: "rgba(100,0,0,1)",
		},
	}
}

func (p *StaticBouncyDebris) FillStyle() string {
	return p.jsColors.FillStyle
}

func (p *StaticBouncyDebris) StrokeStyle() string {
	return p.jsColors.StrokeStyle
}
