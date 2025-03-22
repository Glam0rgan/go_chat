package hash

import "hash/fnv"

func EHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	sum := h.Sum32()

	return uint32(sum)
}