package combat

type StatId uint8

const (
	Agility       = StatId(1)
	Strength      = StatId(2)
	MagicStrength = StatId(3)
	Vitality      = StatId(4)
	MagicVitality = StatId(5)
)

type DamageBasis struct {
	Damage
	ScalingStat  StatId
	ScalingMulti float32
	Variance     float32
}

type Attack struct {
	Id          uint16
	Name        string
	Damages     []DamageBasis
	SelfDamages []DamageBasis
}
