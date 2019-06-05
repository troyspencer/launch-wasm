package bodies

type Debris struct {
	Body
}

func NewDebris() *Debris {
	return &Debris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(100,100,100,1)",
				strokeStyle: "rgba(100,100,100,1)",
			},
			sticky:  false,
			bouncy:  false,
			breaks:  false,
			absorbs: false,
		},
	}
}
