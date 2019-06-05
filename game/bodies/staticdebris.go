package bodies

type StaticDebris struct {
	Body
}

func NewStaticDebris() *StaticDebris {
	return &StaticDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgba(50,50,50,1)",
				strokeStyle: "rgba(50,50,50,1)",
			},
			sticky:  false,
			bouncy:  false,
			breaks:  false,
			absorbs: false,
		},
	}
}
