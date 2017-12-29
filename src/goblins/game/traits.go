package game

type TraitId uint8

const (
	IronWill = TraitId(iota)
	Bookish
	Brutish
	Dextrous
	Stalwart
	Student
	Unpredictable
	Temper
	NullTrait = 0xFF
)

var AllTraits = [...]TraitId{
	IronWill,
	Bookish,
	Brutish,
	Dextrous,
	Stalwart,
	Student,
	Unpredictable,
	Temper,
	NullTrait,
}

func (t TraitId) AsU64() uint64 {
	return uint64(t)
}

func (t TraitId) Name() string {
	switch t {
	case IronWill:
		return "IronWill"
	case Bookish:
		return "Bookish"
	case Brutish:
		return "Brutish"
	case Dextrous:
		return "Dextrous"
	case Stalwart:
		return "Stalwart"
	case Student:
		return "Student"
	case Unpredictable:
		return "Unpredictable"
	case Temper:
		return "Temper"
	case NullTrait:
		return "No Trait"
	default:
		return "[Unknown Trait Id]"
	}
}
