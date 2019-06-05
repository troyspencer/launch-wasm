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
			sticky:  true,
			bouncy:  false,
			breaks:  false,
			absorbs: false,
		},
	}
}
