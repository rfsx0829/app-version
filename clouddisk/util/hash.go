package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// HashValue sum hash
func HashValue(data []byte) string {
	return sumMD5(data)
}

func sumMD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func sumSHA256(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
