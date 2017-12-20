package combat

type DamageType uint8

const (
	PhysDamage  = DamageType(1)
	MagicDamage = DamageType(2)
)

type Damage struct {
	Type     DamageType
	Amount   int32
	Pierce   float32
	Statuses []struct {
		Status      StatusId
		Probability float32
	}
}
