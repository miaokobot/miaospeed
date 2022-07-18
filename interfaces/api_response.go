package interfaces

type SlaveEntrySlot struct {
	Grouping       string
	ProxyInfo      ProxyInfo
	InvokeDuration int64
	Matrices       []MatrixResponse
}

func (ses *SlaveEntrySlot) Get(idx int) *MatrixResponse {
	if idx < len(ses.Matrices) {
		return &ses.Matrices[idx]
	}
	return nil
}

type SlaveTask struct {
	Request SlaveRequest
	Results []SlaveEntrySlot
}

type SlaveProgress struct {
	Index   int
	Record  SlaveEntrySlot
	Queuing int
}

type SlaveResponse struct {
	ID               string
	MiaoSpeedVersion string

	Error    string
	Result   *SlaveTask
	Progress *SlaveProgress
}
