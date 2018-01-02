package combat

import (
	"goblins/game"
	"math"
	"math/rand"
	"time"
)

type FightStatus uint8

const (
	FightWaiting = FightStatus(iota)
	FightRunning
	FightLeftWin
	FightRightWin
	FightDraw

	FightHz       = 40
	FightWaitTime = 7 * FightHz
	FightTime     = 90 * FightHz
)

type Action struct {
	Attack   *Attack
	TimeCast uint16
}

type pendingDamage struct {
	Amount   int32
	Pierce   float32
	Statuses []StatusId
}

type Fighter struct {
	Monster       *game.Monster
	Health        int32
	MaxHealth     int32
	Energy        int32
	MagicEnergy   int32
	CurrentAction Action
	Statuses      map[StatusId]uint16
}

type Fight struct {
	Status FightStatus
	Left   *Fighter
	Right  *Fighter
	Timer  uint16
	rand   *rand.Rand
}

func NewFight(left, right *Fighter) *Fight {
	return &Fight{
		Status: FightWaiting,
		Left:   left,
		Right:  right,
		Timer:  FightWaitTime,
		rand:   rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func NewFighter(monster *game.Monster) *Fighter {
	statuses := make(map[StatusId]uint16)
	hp := 2*float64(monster.Attrs.Vitality.Value) +
		4.75*math.Sqrt(float64(monster.Attrs.MagicVitality.Value))
	// placeholder
	return &Fighter{
		Monster:     monster,
		Health:      int32(hp),
		MaxHealth:   int32(hp),
		Energy:      100,
		MagicEnergy: 100,
		Statuses:    statuses,
	}
}

func (f *Fighter) Attack(attack *Attack, clock uint16) bool {
	if f.CurrentAction.Attack != nil ||
		attack.StrengthCost > f.Energy ||
		attack.MagicCost > f.MagicEnergy {
		return false
	}
	f.Energy -= int32(attack.StrengthCost)
	f.MagicEnergy -= int32(attack.MagicCost)
	agi := f.Monster.Attrs.Agility.Value
	realCastTime :=
		uint16((1 - 1.265e-5*math.Pow(float64(agi), 1.5)) *
			float64(attack.CastTime))
	f.CurrentAction = Action{
		Attack:   attack,
		TimeCast: clock - realCastTime,
	}
	return true
}

func (f *Fighter) stepAttack(clock uint16) bool {
	if f.CurrentAction.Attack != nil {
		return f.CurrentAction.TimeCast == clock
	}
	return false
}

func (f *Fighter) stepEnergies(clock uint16) {
	if (clock & 31) == 31 {
		f.Energy++
		f.MagicEnergy++
	}
}

func (f *Fighter) hasStatus(status StatusId, clock uint16) bool {
	expire, ok := f.Statuses[status]
	return ok && expire < clock
}

func (self *Fighter) useAttack(tgt *Fighter, clock uint16, r *rand.Rand) {
	atk := self.CurrentAction.Attack
	self.CurrentAction = Action{}
	var tgtStatuses, selfStatuses []*StatusEffect
	for n, _ := range atk.Damages {
		dmg := &atk.Damages[n]
		baseDmg := dmg.ScalingFunc.Scale(dmg.Amount,
			self.Monster.Stat(dmg.ScalingStat).Value, dmg.ScalingMulti)
		baseDmg *= 1.0 + rand.Float32()*dmg.Variance
		multi := float32(1.0)
		if self.hasStatus(Angry, clock) {
			multi += 0.1
		}
		if self.hasStatus(Stoic, clock) {
			multi -= 0.1
		}
		if self.hasStatus(Charmed, clock) {
			multi -= 0.1
		}
		tgt.takeDamage(baseDmg, multi, dmg.Type, clock)
		for k, _ := range dmg.Statuses {
			sts := &dmg.Statuses[k]
			if rand.Float32() < sts.Probability {
				tgtStatuses = append(tgtStatuses, sts)
			}
		}
	}
	if tgtStatuses != nil {
		tgt.gainStatuses(tgtStatuses, clock)
	}

	for n, _ := range atk.SelfDamages {
		dmg := &atk.SelfDamages[n]
		baseDmg := dmg.ScalingFunc.Scale(dmg.Amount,
			self.Monster.Stat(dmg.ScalingStat).Value, dmg.ScalingMulti)
		baseDmg *= 1.0 + rand.Float32()*dmg.Variance
		multi := float32(1.0)
		self.takeDamage(baseDmg, multi, dmg.Type, clock)
		for k, _ := range dmg.Statuses {
			sts := &dmg.Statuses[k]
			if rand.Float32() < sts.Probability {
				selfStatuses = append(selfStatuses, sts)
			}
		}
	}
	if selfStatuses != nil {
		self.gainStatuses(selfStatuses, clock)
	}
}

func (f *Fighter) takeDamage(amount float32, multi float32,
	dmgType DamageTypeId, clock uint16) {
	if f.hasStatus(Angry, clock) {
		multi += 0.1
	}
	if f.hasStatus(Stoic, clock) {
		multi -= 0.1
	}
	if amount > 0 {
		if dmgType == PhysDamage {
			vit := f.Monster.Attrs.Vitality.Value
			// ~ 0.25/sqrt(1000) * sqrt(vit)
			multi -= float32(0.008 * math.Sqrt(float64(vit)))
		} else if dmgType == MagicDamage {
			mvit := f.Monster.Attrs.MagicVitality.Value
			multi -= 0.65 / 1000 * float32(mvit)
		}
		if multi < 0 {
			// don't heal from attacks
			multi = 0
		}
	}
	// not overflow protected
	f.Health -= int32(amount * multi)
	if f.Health < 0 {
		f.Health = 0
	} else if f.Health > f.MaxHealth {
		f.Health = f.MaxHealth
	}
}

func (f *Fighter) gainStatuses(statuses []*StatusEffect, clock uint16) {
	for _, sts := range statuses {
		expire := uint16(0)
		if sts.Duration < clock {
			expire = clock - sts.Duration
		}
		f.Statuses[sts.Status] = expire
	}
}

func (f *Fight) Step() {
	switch f.Status {
	case FightWaiting:
		if f.Timer == 0 {
			f.Status = FightRunning
			f.Timer = FightTime
		} else {
			f.Timer--
		}
	case FightRunning:
		if f.Left.stepAttack(f.Timer) {
			f.Left.useAttack(f.Right, f.Timer, f.rand)
		}
		if f.Right.stepAttack(f.Timer) {
			f.Right.useAttack(f.Left, f.Timer, f.rand)
		}

		if f.Left.Health <= 0 {
			if f.Right.Health <= 0 {
				f.Status = FightDraw
			}
			f.Status = FightRightWin
		} else if f.Right.Health <= 0 {
			f.Status = FightLeftWin
		} else if f.Timer == 0 {
			f.Status = FightDraw
		} else {
			f.Right.stepEnergies(f.Timer)
			f.Left.stepEnergies(f.Timer)
			f.Timer--
		}
	default:
		// do nothing, status is already win or draw
	}
}
