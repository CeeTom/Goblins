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
