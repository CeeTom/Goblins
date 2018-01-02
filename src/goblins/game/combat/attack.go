package combat

import (
	"goblins/game"
	"math"
)

type ScalingFuncId uint8

const (
	NullAttack = 0xFFFF
)

const (
	Zero = ScalingFuncId(iota)
	Linear
	Exponential
	Logarithmic
)

var AllScalingFuncs = [...]ScalingFuncId{
	Zero,
	Linear,
	Exponential,
	Logarithmic,
}

func (scaleF ScalingFuncId) Scale(baseDmg int32, statValue uint16,
	multi float32) float32 {
	switch scaleF {
	case Zero:
		return float32(baseDmg)
	case Linear:
		return float32(baseDmg) + float32(statValue)*multi
	case Exponential:
		m64 := float64(multi)
		v64 := float64(statValue)
		sign := float64(1)
		if m64 < 0 {
			sign = -1
			m64 *= -1
		}
		return float32(baseDmg) + float32(sign*math.Pow(v64/200, m64))
	case Logarithmic:
		m64 := float64(multi)
		v64 := float64(statValue)
		sign := float64(1)
		if m64 < 0 {
			sign = -1
			m64 *= -1
		}
		return float32(baseDmg) + float32(sign*math.Pow(v64, 1/m64))
	default:
		panic("bad scaling func")
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
	ScalingStat  game.StatId
	ScalingMulti float32
	Variance     float32
}

type Attack struct {
	Id           uint16
	Name         string
	StrengthCost int32
	MagicCost    int32
	CastTime     uint16
	Damages      []DamageBasis
	SelfDamages  []DamageBasis
}
