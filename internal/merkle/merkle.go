package merkle

import (
	"crypto/sha256"
	"fmt"
	"math"
)

func GetSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func IsRelativeOnTheRight(pieceIndex, curTreeHeight int) bool {
	twoPowOfHeight := int(math.Pow(2, float64(curTreeHeight)))
	return (pieceIndex % twoPowOfHeight) < twoPowOfHeight/2
}

func ConcatByteSlices(first, second []byte) []byte {
	return []byte(fmt.Sprintf("%s%s", first, second))
}

func GetSiblingIndex(pieceIndex int) int {
	if IsRelativeOnTheRight(pieceIndex, 1) {
		return pieceIndex + 1
	}
	return pieceIndex - 1
}
