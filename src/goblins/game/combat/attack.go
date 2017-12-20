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
