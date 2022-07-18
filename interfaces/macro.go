package interfaces

type SlaveRequestMacroType string

const (
	MacroSpeed SlaveRequestMacroType = "SPEED"

	MacroPing   SlaveRequestMacroType = "PING"
	MacroUDP    SlaveRequestMacroType = "UDP"
	MacroScript SlaveRequestMacroType = "SCRIPT"
	MacroGeo    SlaveRequestMacroType = "GEO"

	MacroInvalid SlaveRequestMacroType = "INVALID"
)

// Macro is the atom runner for a job. Since some matrices
// could be combined, e.g. HTTPPing / RTTPing, so instead of
// triggering two similar jobs, we only run a macro job once
// and return attributes for multiple matrices
type SlaveRequestMacro interface {
	// define the macro type to match
	Type() SlaveRequestMacroType

	// define the task for the macro,
	Run(proxy Vendor, request *SlaveRequest) error
}
