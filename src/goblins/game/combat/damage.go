package combat

type DamageTypeId uint8

const (
	PhysDamage = DamageTypeId(1 + iota)
	MagicDamage
)

var AllDamageTypes = [...]DamageTypeId{
	PhysDamage,
	MagicDamage,
}

func (dmg DamageTypeId) AsU64() uint64 {
	return uint64(dmg)
}

func (dmg DamageTypeId) Name() string {
	switch dmg {
	case PhysDamage:
		return "PhysDamage"
	case MagicDamage:
		return "MagicDamage"
	default:
		return "[Unknown Damage Type Id]"
	}
}

type Damage struct {
	Type     DamageTypeId
	Amount   int32
	Pierce   float32
	Statuses []struct {
		Status      StatusId
		Probability float32
	}
}
