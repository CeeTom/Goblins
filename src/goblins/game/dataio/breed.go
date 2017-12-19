package dataio

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"goblins/game"
	"io"
)

// "BREED " 0x00 0x00
const breedHeader = 0x0000326869698266

func ReadBreed(r io.Reader) (*game.Breed, error) {
	ret := new(game.Breed)

	head, err := readU64(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read breed header")
	}
	if head != breedHeader {
		msg := fmt.Sprintf("Invalid breed header: %0#16x", head)
		return nil, errors.New(msg)
	}

	err = binary.Read(r, binary.LittleEndian, &ret.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read breed id")
	}

	err = binary.Read(r, binary.LittleEndian, &ret.ProclivityTable)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read proclivity table")
	}

	err = binary.Read(r, binary.LittleEndian, &ret.AttackTable)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack table")
	}

	err = binary.Read(r, binary.LittleEndian, &ret.TraitTable)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read trait table")
	}

	ret.Name, err = readString(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read breed name")
	}

	return ret, nil
}
