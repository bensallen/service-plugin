package state

// Value is the int representation of a Service state
type Value int

const (
	// Unknown state is a service's initial state
	Unknown Value = 0
	// Active state is a service that has successfully started
	Active Value = 1
	// Exited state is a service that exited unexpectedly
	Exited Value = 2
	// Stopped state is a service that has successfully stopped
	Stopped Value = 3
	// Starting state is a service that has initated Start()
	Starting Value = 4
	// Stopping state is a service that has initated Stop()
	Stopping Value = 5
	// FailedRequire state is a service that has specified a Required service,
	// and on Start() one or more of those services were not in the Active state.
	FailedRequire Value = 6
)

// String returns the string representation a state. Returns "Unknown" if
// given state is undefined.
func (v Value) String() string {
	states := [...]string{
		"Unknown",
		"Active",
		"Exited",
		"Stopped",
		"Starting",
		"Stopping",
		"Failed Require",
	}

	if v < Unknown || v > FailedRequire {
		return "Unknown"
	}

	return states[v]
}
