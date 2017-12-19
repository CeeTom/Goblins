package game

const (
	MaxBreedTraits         = 3
	NumberOfStats          = 5
	ProclivityTableColumns = 32
	TraitTableColumns      = 32
)

type BreedId uint16

type Breed struct {
	Id              BreedId
	Name            string
	ProclivityTable [NumberOfStats][ProclivityTableColumns]float32
	AttackTable     []AttackId
	TraitTable      [MaxBreedTraits][TraitTableColumns]TraitId
}

func SeedToParts(seed Seed) (BreedId, [NumberOfStats]uint8, [MaxBreedTraits]uint8) {
	breedSeed := uint16(seed[0] & 0x07)
	breedSeed = breedSeed | (uint16(seed[1]&0x03) << 3)
	breedSeed = breedSeed | (uint16(seed[2]&0x01) << 5)
	breedSeed = breedSeed | (uint16(seed[3]&0x02) << 6)
	breedSeed = breedSeed | (uint16(seed[4]&0x04) << 7)
	breedSeed = breedSeed | (uint16(seed[5]&0x08) << 8)
	breedSeed = breedSeed | (uint16(seed[6]&0x10) << 9)
	breedSeed = breedSeed | (uint16(seed[7]&0x20) << 10)
	breedSeed = breedSeed | (uint16(seed[8]&0x01) << 11)

	statSeed := uint32(seed[0] & 0x28)
	statSeed = statSeed | (uint32(seed[1]&0x0C) << 2)
	statSeed = statSeed | (uint32(seed[2]&0x0E) << 4)
	statSeed = statSeed | (uint32(seed[3]&0x0D) << 7)
	statSeed = statSeed | (uint32(seed[4]&0x38) << 10)
	statSeed = statSeed | (uint32(seed[5]&0x16) << 13)
	statSeed = statSeed | (uint32(seed[6]&0x23) << 16)
	statSeed = statSeed | (uint32(seed[7]&0x07) << 19)
	statSeed = statSeed | (uint32(seed[8]&0x0E) << 22)

	statArray := [...]uint8{
		uint8(statSeed & 0x1F),
		uint8((statSeed >> 5) & 0x1F),
		uint8((statSeed >> 10) & 0x1F),
		uint8((statSeed >> 15) & 0x1F),
		uint8((statSeed >> 20) & 0x1F),
	}

	traitSeed := uint16(seed[0] & 0x10)
	traitSeed = traitSeed | (uint16(seed[1]&0x30) << 1)
	traitSeed = traitSeed | (uint16(seed[2]&0x30) << 3)
	traitSeed = traitSeed | (uint16(seed[3]&0x30) << 5)
	traitSeed = traitSeed | (uint16(seed[4]&0x03) << 7)
	traitSeed = traitSeed | (uint16(seed[5]&0x21) << 9)
	traitSeed = traitSeed | (uint16(seed[6]&0x0C) << 11)
	traitSeed = traitSeed | (uint16(seed[7]&0x18) << 13)
	traitSeed = traitSeed | (uint16(seed[8]&0x10) << 15)

	traitArray := [...]uint8{
		uint8(traitSeed & 0x1F),
		uint8((traitSeed >> 5) & 0x1F),
		uint8((traitSeed >> 10) & 0x1F),
	}

	return BreedId(breedSeed), statArray, traitArray
}
