package bodies

type Water struct {
	Body
}

func NewWater() *Water {
	return &Water{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(0,0,180,0.5)",
				strokeStyle: "rgba(0,0,50,0.5)",
			},
			sticky:  false,
			bouncy:  false,
			breaks:  false,
			absorbs: true,
		},
	}
}
