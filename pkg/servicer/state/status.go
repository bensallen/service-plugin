package state

type Value int

const (
	Unknown  Value = 0
	Active   Value = 1
	Exited   Value = 2
	Stopped  Value = 3
	Starting Value = 4
	Stopping Value = 5
	Enabled  Value = 6
	Disabled Value = 7
)

func (v Value) String() string {
	states := [...]string{
		"Unknown",
		"Active",
		"Exited",
		"Stopped",
		"Starting",
		"Stopping",
		"Enabled",
		"Disabled",
	}

	if v < Unknown || v > Disabled {
		return "Unknown"
	}

	return states[v]
}
