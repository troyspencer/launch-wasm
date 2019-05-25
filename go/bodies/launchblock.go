package bodies

type LaunchBlock struct {
	*StaticDebris
}

func NewLaunchBlock() *LaunchBlock {
	return &LaunchBlock{
		StaticDebris: NewStaticDebris(),
	}
}

func (lb *LaunchBlock) FillStyle() string {
	return lb.StaticDebris.FillStyle()
}

func (lb *LaunchBlock) StrokeStyle() string {
	return lb.StaticDebris.StrokeStyle()
}

func (lb *LaunchBlock) Sticky() bool {
	return lb.StaticDebris.Sticky()
}
