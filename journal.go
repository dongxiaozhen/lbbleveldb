package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
)

func log(fileName string) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		checksum := binary.LittleEndian.Uint32(data[:4])
		length := binary.LittleEndian.Uint16(data[4:])
		chunkType := data[6]
		fmt.Println("log checksum,length,chunkType", checksum, length, chunkType)

		seq, lth, err := decodeBatchHeader(data[7:])
		if err != nil {
			fmt.Println("log decode batch header error", err)
			return
		}
		fmt.Println("log decode header seq,lth", seq, lth)
		peed := 0
		data = data[19:]
		for i := 0; i < lth; i++ {
			peed, err = decodeBatch(data)
			if err != nil {
				fmt.Println("decode batch err", err)
				return
			}
			data = data[peed:]
		}
	}

}

var batchHeaderLen = 12

func decodeBatchHeader(data []byte) (seq uint64, batchLen int, err error) {
	if len(data) < batchHeaderLen {
		return 0, 0, errors.New("len err")
	}

	seq = binary.LittleEndian.Uint64(data)
	batchLen = int(binary.LittleEndian.Uint32(data[8:]))
	if batchLen < 0 {
		return 0, 0, errors.New("len err2")
	}
	return
}

func decodeBatch(data []byte) (int, error) {
	o := 0
	k1, v1 := "", ""
	// Key type.
	skeyType := keyType(data[o])
	if skeyType > keyTypeVal {
		return 0, errors.New("len err3")
	}
	o++

	// Key.
	x, n := binary.Uvarint(data[o:])
	o += n
	if n <= 0 || o+int(x) > len(data) {
		return 0, errors.New("len err4")
	}
	k1 = string(data[o : o+int(x)])
	o += int(x)

	// Value.
	if skeyType == keyTypeVal {
		x, n = binary.Uvarint(data[o:])
		o += n
		if n <= 0 || o+int(x) > len(data) {
			return 0, errors.New("len err5")
		}
		v1 = string(data[o : o+int(x)])
		o += int(x)
	} else {
	}
	fmt.Printf("logbatch type:%d,k:%s,v:%s\n", skeyType, string(k1), string(v1))
	return o, nil
}

type keyType uint

func (kt keyType) String() string {
	switch kt {
	case keyTypeDel:
		return "d"
	case keyTypeVal:
		return "v"
	}
	return fmt.Sprintf("<invalid:%#x>", uint(kt))
}

// Value types encoded as the last component of internal keys.
// Don't modify; this value are saved to disk.
const (
	keyTypeDel = keyType(0)
	keyTypeVal = keyType(1)
)

type batchIndex struct {
	keyType            keyType
	keyPos, keyLen     int
	valuePos, valueLen int
}

func (index batchIndex) k(data []byte) []byte {
	return data[index.keyPos : index.keyPos+index.keyLen]
}

func (index batchIndex) v(data []byte) []byte {
	if index.valueLen != 0 {
		return data[index.valuePos : index.valuePos+index.valueLen]
	}
	return nil
}

func (index batchIndex) kv(data []byte) (key, value []byte) {
	return index.k(data), index.v(data)
}
