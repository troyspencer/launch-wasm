package bodies

type BouncyDebris struct {
	Body
}

func NewBouncyDebris() *BouncyDebris {
	return &BouncyDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(200,0,0,1)",
				strokeStyle: "rgba(200,0,0,1)",
			},
			sticky:  false,
			bouncy:  true,
			breaks:  false,
			absorbs: false,
		},
	}
}
