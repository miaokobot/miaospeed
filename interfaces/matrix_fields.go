package interfaces

type InvalidDS struct{}

type HTTPPingDS struct {
	Value uint16
}

type RTTPingDS struct {
	Value uint16
}

type AverageSpeedDS struct {
	Value uint64
}

type MaxSpeedDS struct {
	Value uint64
}

type PerSecondSpeedDS struct {
	Max     uint64
	Average uint64
	Speeds  []uint64
}

type UDPTypeDS struct {
	Value string
}

type ScriptTestDS struct {
	Key string
	ScriptResult
}

type InboundGeoIPDS struct {
	MultiStacks
}

type OutboundGeoIPDS struct {
	MultiStacks
}
