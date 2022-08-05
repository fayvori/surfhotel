package handlers

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5HotelSearchHash(params string) string {
	hashes := md5.New()
	hashes.Write([]byte(params))
	return hex.EncodeToString(hashes.Sum(nil))
}
