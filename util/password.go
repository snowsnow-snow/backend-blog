package util

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
)

func GenerateNewPassword(enteredPassword string) (string, string) {
	salt := uuid.NewString()
	encryption := Encryption(enteredPassword, salt)
	return encryption, salt
}
func Encryption(password string, salt string) string {
	hash := md5.New()
	for i := 0; i < 3; i++ {
		hash.Write([]byte(password + salt))
		password = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return password
}
