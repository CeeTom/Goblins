package dataio

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"goblins/game"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	// "BREED " 0x00 0x00
	breedHeader  = 0x0000204445455242
	BreedFileExt = ".brd"
)

func ReadAllBreeds(dirname string) ([]*game.Breed, error) {
	ret := make([]*game.Breed, 0, 64)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read breed directory entries")
	}

	for _, fileInfo := range files {
		// skip directories
		if !fileInfo.IsDir() {
			name := path.Join(dirname, fileInfo.Name())
			if strings.HasSuffix(name, BreedFileExt) {
				file, err := os.Open(name)
				if err != nil {
					msg := fmt.Sprintf("Couldn't open breed file: %s", name)
					return nil, errors.Wrap(err, msg)
				}
				brd, err := ReadBreed(file)
				if err != nil {
					msg := fmt.Sprintf("Couldn't read breed file: %s", name)
					return nil, errors.Wrap(err, msg)
				}
				idx := int(brd.Id)
				if idx >= len(ret) {
					if idx >= cap(ret) {
						newCap := cap(ret) * 2
						if idx >= newCap {
							newCap = idx + 1
						}
						nret := make([]*game.Breed, len(ret), newCap)
						copy(nret, ret)
						ret = nret
					}
					ret = ret[:idx+1]
				}
				if ret[idx] != nil {
					msg := fmt.Sprintf("Breeds %s and %s have same id %d",
						ret[idx].Name, brd.Name, idx)
					return nil, errors.New(msg)
				}
				ret[idx] = brd
			}
		}
	}
	return ret, nil
}

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
