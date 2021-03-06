package combat

type StatusId uint8

const (
	Sleep = StatusId(1 + iota)
	Angry
	Charmed
	Stoic
	Poised
)

var AllStatuses = [...]StatusId{
	Sleep,
	Angry,
	Charmed,
	Stoic,
	Poised,
}

func (status StatusId) AsU64() uint64 {
	return uint64(status)
}

func (status StatusId) Name() string {
	switch status {
	case Sleep:
		return "Sleep"
	case Angry:
		return "Angry"
	case Charmed:
		return "Charmed"
	case Stoic:
		return "Stoic"
	case Poised:
		return "Poised"
	default:
		return "[Unknown Status Id]"
	}
}
