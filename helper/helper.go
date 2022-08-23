package helper

import "encoding/hex"

func GetAddress(s string) [32]byte {
	decodeString, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}
	}

	var result [32]byte

	for i := 0; i < 32; i++ {
		result[i] = decodeString[i]
	}

	return result
}
