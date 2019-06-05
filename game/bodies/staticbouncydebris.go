package bodies

type StaticBouncyDebris struct {
	Body
}

func NewStaticBouncyDebris() *StaticBouncyDebris {
	return &StaticBouncyDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(100,0,0,1)",
				strokeStyle: "rgba(100,0,0,1)",
			},
			sticky:  false,
			bouncy:  true,
			breaks:  false,
			absorbs: false,
		},
	}
}
