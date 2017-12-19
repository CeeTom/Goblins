package game

const (
	MaxMonsterRunes = 9
)

type MonsterAttr uint16
type MonsterMood int16
type AttackId uint16
type Seed [MaxMonsterRunes]uint8

type Monster struct {
	Breed  BreedId
	Seed   Seed
	Name   string
	Traits []TraitId

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
