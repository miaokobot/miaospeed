package interfaces

type SlaveRequestMatrixEntry struct {
	Type   SlaveRequestMatrixType
	Params string
}

type SlaveRequestOptions struct {
	Filter   string
	Matrices []SlaveRequestMatrixEntry
}

func (sro *SlaveRequestOptions) Clone() *SlaveRequestOptions {
	return &SlaveRequestOptions{
		Filter:   sro.Filter,
		Matrices: cloneSlice(sro.Matrices),
	}
}

type SlaveRequestBasics struct {
	ID        string
	Slave     string
	SlaveName string
	Invoker   string
	Version   string
}

func (srb *SlaveRequestBasics) Clone() *SlaveRequestBasics {
	return &SlaveRequestBasics{
		ID:        srb.ID,
		Slave:     srb.Slave,
		SlaveName: srb.SlaveName,
		Invoker:   srb.Invoker,
		Version:   srb.Version,
	}
}

type SlaveRequestNode struct {
	Name    string
	Payload string
}

func (srn *SlaveRequestNode) Clone() *SlaveRequestNode {
	return &SlaveRequestNode{
		Name:    srn.Name,
		Payload: srn.Payload,
	}
}

type SlaveRequest struct {
	Basics  SlaveRequestBasics
	Options SlaveRequestOptions
	Configs SlaveRequestConfigs

	Vendor VendorType
	Nodes  []SlaveRequestNode

	Challenge string
}

func (sr *SlaveRequest) Clone() *SlaveRequest {
	return &SlaveRequest{
		Basics:    *sr.Basics.Clone(),
		Options:   *sr.Options.Clone(),
		Configs:   *sr.Configs.Clone(),
		Nodes:     cloneSlice(sr.Nodes),
		Challenge: sr.Challenge,
	}
}
