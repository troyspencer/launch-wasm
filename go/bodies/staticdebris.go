package bodies

type StaticDebris struct {
	jsColors JSColors
}

func NewStaticDebris() *StaticDebris {
	return &StaticDebris{
		jsColors: JSColors{
			FillStyle:   "rgba(50,50,50,1)",
			StrokeStyle: "rgba(50,50,50,1)",
		},
	}
}

func (p *StaticDebris) FillStyle() string {
	return p.jsColors.FillStyle
}

func (p *StaticDebris) StrokeStyle() string {
	return p.jsColors.StrokeStyle
}
