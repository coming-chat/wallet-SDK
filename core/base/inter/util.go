package inter

import "strings"

// IsValidHexString
func IsValidHexString(str string) (valid bool, length int) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:] // remove 0x prefix
	}
	for _, ch := range []byte(str) {
		valid := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
		if !valid {
			return false, 0
		}
	}
	return true, len(str)
}
