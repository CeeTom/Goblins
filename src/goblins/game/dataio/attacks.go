package dataio

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"goblins/game/combat"
	"io"
)

// "ATTACK" 0x00 0x00
const attackHeader = 0x0000756765848465

func ReadAttack(r io.Reader) (*combat.Attack, error) {
	ret := new(Attack)

	head, err := readU64(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack header")
	}
	if head != attackHeader {
		msg := fmt.Sprintf("Invalid attack header: %0#16x", head)
		return nil, errors.New(msg)
	}

	err = binary.Read(r, binary.LittleEndian, &ret.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack id")
	}

	dmgCount, err := readU16(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack damage count")
	}

	ret.Damages = make([]combat.DamageBasis, int(dmgCount))

	for i := 0; i < int(dmgCount); i++ {
		err = readDamageBasis(r, &ret.Damages[i])
		if err != nil {
			msg := fmt.Sprintf("Couldn't read %d-th attack damage", i)
			return nil, errors.Wrap(err, msg)
		}
	}

	selfDmgCount, err := readU16(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack self damage count")
	}

	ret.SelfDamages = make([]combat.DamageBasis, int(selfDmgCount))

	for i := 0; i < int(selfDmgCount); i++ {
		err = readDamageBasis(r, &ret.SelfDamages[i])
		if err != nil {
			msg := fmt.Sprintf("Couldn't read %d-th attack self damage", i)
			return nil, errors.Wrap(err, msg)
		}
	}
}

func readDamageBasis(r io.Reader, basis *combat.DamageBasis) error {

}
