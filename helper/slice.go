package helper

func SliceCopy(b []byte) []byte {
	var copyB = make([]byte, len(b))
	copy(copyB, b)
	return copyB
}
