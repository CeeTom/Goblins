package game

import (
	"math/rand"
)

const (
	MaxMonsterRunes = 9

	MinFatigue = -1000
	MaxFatigue = 1000
)

type AttackId uint16
type Seed [MaxMonsterRunes]uint8

type MonsterMeasure struct {
	Time  Time
	Value int16
}

type MonsterAttr struct {
	Value      uint16
	Proclivity float32
}

type Monster struct {
	Breed     BreedId
	Seed      Seed
	Name      string
	Traits    []TraitId
	Attacks   []AttackId
	Happiness int16

	Attrs struct {
		Agility       MonsterAttr
		Strength      MonsterAttr
		MagicStrength MonsterAttr
		Vitality      MonsterAttr
		MagicVitality MonsterAttr
	}

	measures struct {
		Fatigue MonsterMeasure
		Age     MonsterMeasure
	}
}

func (m *Monster) HasTrait(trait TraitId) bool {
	for _, t := range m.Traits {
		if t == trait {
			return true
		}
	}
	return false
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

func (m *Monster) Fatigue(clock Time) int16 {
	t := clock - m.measures.Fatigue.Time
	newFatigue := int64(m.measures.Fatigue.Value) - int64(t/3600)
	if newFatigue < MinFatigue {
		newFatigue = MinFatigue
	}
	return int16(newFatigue)
}

func (m *Monster) AddFatigue(amount int16, clock Time) {
	newFatigue := int32(m.Fatigue(clock)) + int32(amount)
	if newFatigue < MinFatigue {
		newFatigue = MinFatigue
	} else if newFatigue > MaxFatigue {
		newFatigue = MaxFatigue
	}
	m.measures.Fatigue.Value = int16(newFatigue)
	m.measures.Fatigue.Time = clock
}

func (m *Monster) Age(clock Time) int16 {
	ageTicks := Time(m.measures.Age.Value) + clock - m.measures.Age.Time
	return int16(ageTicks / 4096)
}

// calls global rand, may need to change
func randProcMulti(p float32) float32 {
	return 1 + rand.Float32()/(p+1) - 1/(p+1)
}

func (attr *MonsterAttr) gainStat(gain int8, traitAdd float32) {
	multi := randProcMulti(attr.Proclivity) + (traitAdd)
	if gain < 0 {
		multi = 1 / multi
	}
	newStat := int32(float32(gain)*multi) + int32(attr.Value)
	if newStat < 1 {
		newStat = 1
	}
	if newStat > 1000 {
		newStat = 1000
	}
	attr.Value = uint16(newStat)
}

func (m *Monster) Exercise(training *Training, clock Time) {
	m.AddFatigue(training.EnergyCost, clock)
	m.Happiness += training.Gains.Happiness

	strMulti := float32(0.0)
	magMulti := float32(0.0)
	agiMulti := float32(0.0)
	defMulti := float32(0.0)

	for _, trait := range m.Traits {
		switch trait {
		case Bookish:
			strMulti -= 0.1
			magMulti += 0.1
		case Brutish:
			strMulti += 0.1
			magMulti -= 0.1
		case Dextrous:
			defMulti -= 0.1
			agiMulti += 0.1
		case Stalwart:
			defMulti += 0.1
			agiMulti -= 0.1
		}
	}

	m.Attrs.Strength.gainStat(training.Gains.Strength, strMulti)
	m.Attrs.MagicStrength.gainStat(training.Gains.MagicStrength, magMulti)
	m.Attrs.Vitality.gainStat(training.Gains.Vitality, strMulti+defMulti)
	m.Attrs.MagicVitality.gainStat(training.Gains.MagicVitality,
		magMulti+defMulti)
	m.Attrs.Agility.gainStat(training.Gains.Agility, agiMulti)
}
