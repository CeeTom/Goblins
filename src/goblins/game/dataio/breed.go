package dataio

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"goblins/game"
	"io"
)

// "BREED " 0x00 0x00
const breedHeader = 0x0000204445455242

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

	err = binary.Read(r, little, &ret.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read breed id")
	}

	err = binary.Read(r, little, &ret.ProclivityTable)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read proclivity table")
	}

	err = binary.Read(r, little, &ret.AttackTable)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack table")
	}

	err = binary.Read(r, little, &ret.TraitTable)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read trait table")
	}

	ret.Name, err = readString(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read breed name")
	}

	return ret, nil
}

func WriteBreed(w io.Writer, breed *game.Breed) error {
	err := binary.Write(w, little, uint64(breedHeader))
	if err != nil {
		return errors.Wrap(err, "Couldn't write breed header")
	}

	err = binary.Write(w, little, breed.Id)
	if err != nil {
		return errors.Wrap(err, "Couldn't write breed id")
	}

	err = binary.Write(w, little, &breed.ProclivityTable)
	if err != nil {
		return errors.Wrap(err, "Couldn't write breed proclivity table")
	}

	err = binary.Write(w, little, &breed.AttackTable)
	if err != nil {
		return errors.Wrap(err, "Couldn't write breed attack table")
	}

	err = binary.Write(w, little, &breed.TraitTable)
	if err != nil {
		return errors.Wrap(err, "Couldn't write breed trait table")
	}

	err = writeString(w, breed.Name)
	if err != nil {
		return errors.Wrap(err, "Couldn't write breed name")
	}

	return nil
}
