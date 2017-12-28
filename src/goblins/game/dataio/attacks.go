package dataio

import (
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"goblins/game/combat"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// "ATTACK" 0x00 0x00
const attackHeader = 0x00004B4341545441

func ReadAllAttacks(dirname string) ([]*combat.Attack, error) {
	ret := make([]*combat.Attack, 0, 64)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack directory entries")
	}

	for _, fileInfo := range files {
		// skip directories
		if !fileInfo.IsDir() {
			name := path.Join(dirname, fileInfo.Name())
			if strings.HasSuffix(name, ".atk") {
				file, err := os.Open(name)
				if err != nil {
					msg := fmt.Sprintf("Couldn't open attack file: %s", name)
					return nil, errors.Wrap(err, msg)
				}
				atk, err := ReadAttack(file)
				if err != nil {
					msg := fmt.Sprintf("Couldn't read attack file: %s", name)
					return nil, errors.Wrap(err, msg)
				}
				idx := int(atk.Id)
				if idx >= len(ret) {
					if idx >= cap(ret) {
						nret := make([]*combat.Attack, len(ret), cap(ret)*2)
						copy(nret, ret)
						ret = nret
					}
					ret = ret[:idx+1]
				}
				if ret[idx] != nil {
					msg := fmt.Sprintf("Attacks %s and %s have same id %d",
						ret[idx].Name, atk.Name, idx)
					return nil, errors.New(msg)
				}
				ret[idx] = atk
			}
		}
	}
	return ret, nil
}

func ReadAttack(r io.Reader) (*combat.Attack, error) {
	ret := new(combat.Attack)

	head, err := readU64(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack header")
	}
	if head != attackHeader {
		msg := fmt.Sprintf("Invalid attack header: %0#16x", head)
		return nil, errors.New(msg)
	}

	err = binary.Read(r, little, &ret.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack id")
	}

	err = binary.Read(r, little, &ret.StrengthCost)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack strength cost")
	}

	err = binary.Read(r, little, &ret.MagicCost)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack magic cost")
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

	ret.Name, err = readString(r)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't read attack name")
	}

	return ret, nil
}

func readDamageBasis(r io.Reader, basis *combat.DamageBasis) error {
	err := binary.Read(r, little, &basis.Type)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage type")
	}

	err = binary.Read(r, little, &basis.Amount)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage amount")
	}

	err = binary.Read(r, little, &basis.Pierce)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage pierce")
	}

	err = binary.Read(r, little, &basis.ScalingFunc)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage scaling function")
	}

	err = binary.Read(r, little, &basis.ScalingStat)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage scaling stat")
	}

	err = binary.Read(r, little, &basis.ScalingMulti)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage scaling multi")
	}

	err = binary.Read(r, little, &basis.Variance)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage variance")
	}

	var statusCount uint16
	err = binary.Read(r, little, &statusCount)
	if err != nil {
		return errors.Wrap(err, "Couldn't read damage status count")
	}

	basis.Statuses =
		make([]struct {
			Status      combat.StatusId
			Probability float32
		}, int(statusCount))

	for i := 0; i < int(statusCount); i++ {
		err = binary.Read(r, little, &basis.Statuses[i])
		if err != nil {
			return errors.Wrap(err, "Couldn't read damage status")
		}
	}

	return nil
}

func WriteAttack(w io.Writer, attack *combat.Attack) error {
	err := binary.Write(w, little, uint64(attackHeader))
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack header")
	}

	err = binary.Write(w, little, attack.Id)
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack id")
	}

	err = binary.Write(w, little, attack.StrengthCost)
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack strength cost")
	}

	err = binary.Write(w, little, attack.MagicCost)
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack magic cost")
	}

	err = binary.Write(w, little, uint16(len(attack.Damages)))
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack damages count")
	}

	for i := range attack.Damages {
		err = writeDamageBasis(w, &attack.Damages[i])
		if err != nil {
			msg := fmt.Sprintf("Couldn't write %d-th attack damage", i)
			return errors.Wrap(err, msg)
		}
	}

	err = binary.Write(w, little, uint16(len(attack.SelfDamages)))
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack self damages count")
	}

	for i := range attack.SelfDamages {
		err = writeDamageBasis(w, &attack.SelfDamages[i])
		if err != nil {
			msg := fmt.Sprintf("Couldn't write %d-th self attack damage", i)
			return errors.Wrap(err, msg)
		}
	}

	err = writeString(w, attack.Name)
	if err != nil {
		return errors.Wrap(err, "Couldn't write attack name")
	}
	return nil
}

func writeDamageBasis(w io.Writer, basis *combat.DamageBasis) error {
	err := binary.Write(w, little, basis.Type)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage type")
	}

	err = binary.Write(w, little, basis.Amount)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage amount")
	}

	err = binary.Write(w, little, basis.Pierce)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage pierce")
	}

	err = binary.Write(w, little, basis.ScalingFunc)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage scaling function")
	}

	err = binary.Write(w, little, basis.ScalingStat)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage scaling stat")
	}

	err = binary.Write(w, little, basis.ScalingMulti)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage scaling multi")
	}

	err = binary.Write(w, little, basis.Variance)
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage variance")
	}

	err = binary.Write(w, little, uint16(len(basis.Statuses)))
	if err != nil {
		return errors.Wrap(err, "Couldn't write damage status count")
	}

	for i := range basis.Statuses {
		err = binary.Write(w, little, &basis.Statuses[i])
		if err != nil {
			msg := fmt.Sprintf("Couldn't write %d-th damage status", i)
			return errors.Wrap(err, msg)
		}
	}

	return nil
}
