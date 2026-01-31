package utils

import "math/big"

func EncodeBase62(s string) string {
	i := new(big.Int)
	i.SetBytes([]byte(s))
	return i.Text(62)
}

func DecodeBase62(s string) string {
	i := new(big.Int)
	i, ok := i.SetString(s, 62)
	if !ok || i == nil {
		return ""
	}
	return string(i.Bytes())
}
