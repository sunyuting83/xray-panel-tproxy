package datafactory

import "encoding/base64"

// decodeBytes decode Bytes
func DecodeBytes(a string) (b string) {
	bb, _ := base64.RawURLEncoding.DecodeString(a)
	return string(bb)
}

// deCodeBytes
func DeCodeBytes(a string) (b string) {
	var str []byte = []byte(a)
	decodeBytes := make([]byte, base64.StdEncoding.DecodedLen(len(str))) // 计算解码后的长度
	base64.StdEncoding.Decode(decodeBytes, str)
	return string(decodeBytes)
}

// DeBase decode base
func DeBase(s string) (k string) {
	decoded, _ := base64.StdEncoding.DecodeString(s)
	k = string(decoded)
	return
}
