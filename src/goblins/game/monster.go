package game

const (
	MaxMonsterRunes = 9
)

type MonsterMood int16
type AttackId uint16
type Seed [MaxMonsterRunes]uint8

type MonsterAttr struct {
	Value      uint16
	Proclivity float32
}

type Monster struct {
	Breed   BreedId
	Seed    Seed
	Name    string
	Traits  []TraitId
	Attacks []AttackId

	Attrs struct {
		Agility       MonsterAttr
		Strength      MonsterAttr
		MagicStrength MonsterAttr
		Vitality      MonsterAttr
		MagicVitality MonsterAttr
	}

	Moods struct {
		Happiness MonsterMood
		Fatigue   MonsterMood
	}
}

func (m *Monster) Stat(stat StatId) *MonsterAttr {
	switch stat {
	case Agility:
		return &m.Attrs.Agility
	case Strength:
		return &m.Attrs.Strength
	case MagicStrength:
		return &m.Attrs.MagicStrength
	case Vitality:
		return &m.Attrs.Vitality
	case MagicVitality:
		return &m.Attrs.MagicVitality
	default:
		return nil
	}
}
