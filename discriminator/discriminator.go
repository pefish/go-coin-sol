package discriminator

import (
	"crypto/sha256"
)

func GetDiscriminator(namespace string, funcName string) []byte {

	preimage := namespace + ":" + funcName
	hash := sha256.Sum256([]byte(preimage))

	var sighash [8]byte
	copy(sighash[:], hash[:8]) // 复制前 8 个字节

	return sighash[:]

}
