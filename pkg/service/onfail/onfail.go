package onfail

type Action int

const (
	Nothing Action = 0
	Restart Action = 1
)

func (a Action) String() string {
	actions := [...]string{
		"Nothing",
		"Restart",
	}

	if a < Nothing || a > Restart {
		return "Nothing"
	}

	return actions[a]
}
