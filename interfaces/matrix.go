package interfaces

type SlaveRequestMatrixType string

const (
	MatrixAverageSpeed   SlaveRequestMatrixType = "SPEED_AVERAGE"
	MatrixMaxSpeed       SlaveRequestMatrixType = "SPEED_MAX"
	MatrixPerSecondSpeed SlaveRequestMatrixType = "SPEED_PER_SECOND"

	MatrixUDPType SlaveRequestMatrixType = "UDP_TYPE"

	MatrixInboundGeoIP  SlaveRequestMatrixType = "GEOIP_INBOUND"
	MatrixOutboundGeoIP SlaveRequestMatrixType = "GEOIP_OUTBOUND"

	MatrixScriptTest SlaveRequestMatrixType = "TEST_SCRIPT"
	MatrixHTTPPing   SlaveRequestMatrixType = "TEST_PING_CONN"
	MatrixRTTPing    SlaveRequestMatrixType = "TEST_PING_RTT"

	MatrixInvalid SlaveRequestMatrixType = "INVALID"
)

func (srmt *SlaveRequestMatrixType) Valid() bool {
	if srmt == nil {
		return false
	}

	switch *srmt {
	case MatrixAverageSpeed, MatrixMaxSpeed, MatrixPerSecondSpeed,
		MatrixUDPType,
		MatrixInboundGeoIP, MatrixOutboundGeoIP,
		MatrixScriptTest, MatrixHTTPPing, MatrixRTTPing:
		return true
	}

	return false
}

// Matrix is the the atom attribute for a job
// e.g. to fetch the RTTPing of a node,
// it calls RTTPing matrix, which would initiate
// a ping macro and return the RTTPing attribute
type SlaveRequestMatrix interface {
	// define the matrix type to match
	Type() SlaveRequestMatrixType

	// define which macro job to run
	MacroJob() SlaveRequestMacroType

	// define the function to extract attribute
	// from macro result
	Extract(SlaveRequestMatrixEntry, SlaveRequestMacro)
}

type MatrixResponse struct {
	Type    SlaveRequestMatrixType
	Payload string
}
