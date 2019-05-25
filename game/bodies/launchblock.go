package bodies

type LaunchBlock struct {
	*StaticDebris
}

func NewLaunchBlock() *LaunchBlock {
	return &LaunchBlock{
		StaticDebris: NewStaticDebris(),
	}
}
