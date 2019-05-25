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
			breaks: false,
		},
	}
}
