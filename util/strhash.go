package util

func StringHash(s string) (hash uint16) {
	for _, c := range s {
		ch := uint16(c)
		hash = hash + ((hash) << 5) + ch + (ch << 7)
	}
	return
}
