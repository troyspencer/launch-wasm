package bodies

type LaunchBlock struct {
	*StaticDebris
}

func NewLaunchBlock() *LaunchBlock {
	return &LaunchBlock{
		StaticDebris: NewStaticDebris(),
	}
}

func (b *LaunchBlock) FillStyle() string {
	return b.StaticDebris.FillStyle()
}

func (b *LaunchBlock) StrokeStyle() string {
	return b.StaticDebris.StrokeStyle()
}

func (b *LaunchBlock) Sticky() bool {
	return b.StaticDebris.Sticky()
}

func (b *LaunchBlock) Bouncy() bool {
	return b.StaticDebris.Bouncy()
}
