package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

//sha1 密码加密
func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum(nil))
}

//MD5 加密
func Md5(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	return hex.EncodeToString(md5.Sum(nil))
}
