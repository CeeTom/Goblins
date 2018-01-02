package game

type TrainingId uint8

type Training struct {
	Id         TrainingId
	Name       string
	EnergyCost int16
	Gains      struct {
		Happiness     int16
		Agility       int8
		Strength      int8
		MagicStrength int8
		Vitality      int8
		MagicVitality int8
	}
}
