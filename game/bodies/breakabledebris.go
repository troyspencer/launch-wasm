package bodies

type BreakableDebris struct {
	Body
}

func NewBreakableDebris() *BreakableDebris {
	return &BreakableDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(100,100,100,1)",
				strokeStyle: "rgba(255,255,255,1)",
			},
			sticky: false,
			bouncy: false,
			breaks: true,
		},
	}
}
