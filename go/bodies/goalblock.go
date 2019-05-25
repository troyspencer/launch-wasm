package bodies

type GoalBlock struct {
	jsColors JSColors
}

func NewGoalBlock() *GoalBlock {
	return &GoalBlock{
		jsColors: JSColors{
			FillStyle:   "rgba(0, 255,0,1)",
			StrokeStyle: "rgba(0, 255,0,1)",
		},
	}
}

func (p *GoalBlock) FillStyle() string {
	return p.jsColors.FillStyle
}

func (p *GoalBlock) StrokeStyle() string {
	return p.jsColors.StrokeStyle
}
