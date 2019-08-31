package main

const (
	alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	base     uint64 = uint64(len(alphabet))
)

func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func indexOf(str string, char byte) int {
	for i := 0; i < len(str); i++ {
		if str[i] == char {
			return i
		}
	}

	return -1
}

// EncodeNumber takes a base 10 number and
// converts it to a shortened representation.
func EncodeNumber(num uint64) string {
	// Zero case will not be covered by loop
	if num == 0 {
		return string(alphabet[0])
	}

	runes := make([]rune, 0)

	for num > 0 {
		runes = append(runes, rune(alphabet[num%base]))
		num /= base
	}

	return reverse(string(runes))
}

func DecodeNumber(str string) uint64 {
	var num uint64
	for i := 0; i < len(str); i++ {
		num = num*base + uint64(indexOf(alphabet, str[i]))
	}
	return num
}
