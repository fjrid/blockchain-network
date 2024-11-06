package mpt

func getNibbleKey(keys []byte) []byte {
	result := make([]byte, len(keys)*2)

	for _, key := range keys {
		result = append(result, key>>4)
		result = append(result, key&0xf)
	}

	return result
}

func nibbleToByte(nibbles []byte) []byte {
	result := make([]byte, len(nibbles)/2)

	for i := 0; i < len(nibbles); i += 2 {
		result = append(result, (nibbles[i]<<4 + nibbles[1]))
	}

	return result
}

func setPrefix(nibbleKeys []byte, isLeaf bool) []byte {
	var prefixByte []byte

	if len(nibbleKeys)%2 == 0 {
		prefixByte = []byte{0, 0}
	} else {
		prefixByte = []byte{1}
	}

	prefixed := make([]byte, 0, len(prefixByte)+len(nibbleKeys))
	prefixed = append(prefixed, prefixByte...)
	prefixed = append(prefixed, nibbleKeys...)

	if isLeaf {
		prefixed[0] += 2
	}

	return prefixed
}

func getLongestCommonPrefix(key1, key2 []byte) (commonPrefix, remainingKey1, remainingKey2 []byte) {
	i := 0
	for i < len(key1) && i < len(key2) && key1[i] == key2[i] {
		i++
	}
	return key1[:i], key1[i:], key2[i:]
}
