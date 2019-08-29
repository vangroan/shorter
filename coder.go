package main

import "strings"

const (
	alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	base     uint64 = uint64(len(alphabet))
)

// EncodeNumber takes a base 10 number and
// converts it to a shortened representation.
func EncodeNumber(num uint64) string {
	sb := make([]string, 0)

	for num > 0 {
		sb = append(sb, string(alphabet[num%base]))
		num /= base
	}

	return strings.Join(sb, "")
}
