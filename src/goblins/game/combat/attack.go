package combat

type StatId uint8
type ScalingFuncId uint8

const (
	Agility = StatId(iota + 1)
	Strength
	MagicStrength
	Vitality
	MagicVitality
)

const (
	Zero = ScalingFuncId(iota)
	Linear
	Exponential
	Logarithmic
)

var AllStats = [...]StatId{
	Agility,
	Strength,
	MagicStrength,
	Vitality,
	MagicVitality,
}

var AllScalingFuncs = [...]ScalingFuncId{
	Zero,
	Linear,
	Exponential,
	Logarithmic,
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

func (scaleF ScalingFuncId) AsU64() uint64 {
	return uint64(scaleF)
}

func (scaleF ScalingFuncId) Name() string {
	switch scaleF {
	case Zero:
		return "Zero"
	case Linear:
		return "Linear"
	case Exponential:
		return "Exponential"
	case Logarithmic:
		return "Logarithmic"
	default:
		return "[Unknown Scaling Func Id]"
	}
}

type DamageBasis struct {
	Damage
	ScalingFunc  ScalingFuncId
	ScalingStat  StatId
	ScalingMulti float32
	Variance     float32
}

type Attack struct {
	Id           uint16
	Name         string
	StrengthCost int32
	MagicCost    int32
	Damages      []DamageBasis
	SelfDamages  []DamageBasis
}
