package util

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"sort"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

const (
	AppKey = "4409e2ce8ffd12b8"
	AppSec = "59b43e04ad6965f34319062b478f83dd"
)

func Signature(params *map[string]string) {
	var keys []string
	(*params)["appkey"] = AppKey
	for k := range *params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var query string
	for _, k := range keys {
		query += k + "=" + url.QueryEscape((*params)[k]) + "&"
	}
	query = query[:len(query)-1] + AppSec
	hash := md5.New()
	hash.Write([]byte(query))
	(*params)["sign"] = hex.EncodeToString(hash.Sum(nil))
}

func ClientSign(data string) string {
	h1 := sha512.New()
	h2 := sha3.New512()
	h3 := sha512.New384()
	h4 := sha3.New384()
	h5, _ := blake2b.New512(nil)

	h1.Write([]byte(data))
	h2.Write([]byte(hex.EncodeToString(h1.Sum(nil))))
	h3.Write([]byte(hex.EncodeToString(h2.Sum(nil))))
	h4.Write([]byte(hex.EncodeToString(h3.Sum(nil))))
	h5.Write([]byte(hex.EncodeToString(h4.Sum(nil))))
	return hex.EncodeToString(h5.Sum(nil))
}
