package bodies

type StickyDebris struct {
	Body
}

func NewStickyDebris() *StickyDebris {
	return &StickyDebris{
		Body: Body{
			JSColors: JSColors{
				fillStyle:   "rgb(252, 241, 37)",
				strokeStyle: "rgb(252, 241, 37)",
			},
			sticky:  true,
			bouncy:  false,
			breaks:  false,
			absorbs: false,
		},
	}
}
