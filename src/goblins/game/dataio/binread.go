package dataio

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
)

func readU64(r io.Reader) (v uint64, err error) {
	err = binary.Read(r, binary.LitteEndian, &v)
}

func readU32(r io.Reader) (v uint32, err error) {
	err = binary.Read(r, binary.LitteEndian, &v)
}

func readU16(r io.Reader) (v uint16, err error) {
	err = binary.Read(r, binary.LitteEndian, &v)
}

func readS64(r io.Reader) (v int64, err error) {
	err = binary.Read(r, binary.LitteEndian, &v)
}

func readS32(r io.Reader) (v int32, err error) {
	err = binary.Read(r, binary.LitteEndian, &v)
}

func readS16(r io.Reader) (v int16, err error) {
	err = binary.Read(r, binary.LitteEndian, &v)
}

func readString(r io.Reader) (string, error) {
	var l uint16
	if err := binary.Read(r, binary.LittleEndian, &l); err != nil {
		return "", errors.Wrap(err, "Couldn't read string length")
	}
	bytes := make([]byte, int(l))
	read := 0
	for read < int(l) {
		n, err := r.Read(bytes[read:])
		if err != nil {
			return "", errors.Wrap(err, "Couldn't read string")
		}
		read += n
	}
	return string(bytes), nil
}
