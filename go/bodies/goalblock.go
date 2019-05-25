package bodies

type GoalBlock struct {
	Body
}

func NewGoalBlock() *GoalBlock {
	return &GoalBlock{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(0, 255,0,1)",
				strokeStyle: "rgba(0, 255,0,1)",
			},
			sticky: true,
			bouncy: false,
		},
	}
}

func (b *GoalBlock) FillStyle() string {
	return b.fillStyle
}

func (b *GoalBlock) StrokeStyle() string {
	return b.strokeStyle
}

func (b *GoalBlock) Sticky() bool {
	return b.sticky
}

func (b *GoalBlock) Bouncy() bool {
	return b.bouncy
}
