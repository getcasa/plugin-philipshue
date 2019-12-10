package devices

type state struct {
	ON        bool
	Bri       int
	Hue       int
	Sat       int
	Effect    string
	XY        []int
	CT        int
	Alert     string
	ColorMode string
	Mode      string
	Reachable bool
}

type swUpdate struct {
	State       string
	LastInstall interface{}
}

type ct struct {
	Min int
	Max int
}

type control struct {
	MinDimLevel    int
	MaxLumen       int
	ColorGamutType string
	CT             ct
}

type streaming struct {
	Renderer bool
	proxy    bool
}

type capabilities struct {
	Certified bool
	Control   control
	Streaming streaming
}

type config struct {
	Archetype string
	Function  string
	Direction string
}

// Hue define a philips hue light bulb
type Hue struct {
	State            state
	SWUpdate         swUpdate
	Type             string
	Name             string
	ModelID          string
	ManufacturerName string
	ProductName      string
	Capabilities     capabilities
	Config           config
	UniqueID         string
	SWVersion        string
	SWConfigID       string
	ProductID        string
}
