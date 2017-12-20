package dataio

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
)

var little = binary.LittleEndian

func readU64(r io.Reader) (v uint64, err error) {
	err = binary.Read(r, little, &v)
	return
}

func readU32(r io.Reader) (v uint32, err error) {
	err = binary.Read(r, little, &v)
	return
}

func readU16(r io.Reader) (v uint16, err error) {
	err = binary.Read(r, little, &v)
	return
}

func readS64(r io.Reader) (v int64, err error) {
	err = binary.Read(r, little, &v)
	return
}

func readS32(r io.Reader) (v int32, err error) {
	err = binary.Read(r, little, &v)
	return
}

func readS16(r io.Reader) (v int16, err error) {
	err = binary.Read(r, little, &v)
	return
}

func readString(r io.Reader) (string, error) {
	var l uint16
	if err := binary.Read(r, little, &l); err != nil {
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

func writeString(w io.Writer, s string) error {
	bytes := []byte(s)

	err := binary.Write(w, little, uint16(len(bytes)))
	if err != nil {
		return errors.Wrap(err, "Couldn't write string length")
	}
	_, err = w.Write(bytes)
	if err != nil {
		return errors.Wrap(err, "Couldn't write string")
	}
	return nil
}
