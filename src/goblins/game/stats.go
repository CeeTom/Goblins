package game

type StatId uint8

const (
	Agility = StatId(iota + 1)
	Strength
	MagicStrength
	Vitality
	MagicVitality
)

var AllStats = [...]StatId{
	Agility,
	Strength,
	MagicStrength,
	Vitality,
	MagicVitality,
}

func (stat StatId) AsU64() uint64 {
	return uint64(stat)
}

func (stat StatId) Name() string {
	switch stat {
	case Agility:
		return "Agility"
	case Strength:
		return "Strength"
	case MagicStrength:
		return "MagicStrength"
	case Vitality:
		return "Vitality"
	case MagicVitality:
		return "MagicVitality"
	default:
		return "[Unknown Stat Id]"
	}
}
