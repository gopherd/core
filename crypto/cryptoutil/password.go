package cryptoutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"strings"

	"github.com/gopherd/core/math/random"
)

func GenerateSalt(n int) string {
	return random.String(n, random.CryptoSource, random.O_DIGIT|random.O_LOWER_CHAR|random.O_UPPER_CHAR)
}

func GeneratePasswordSalt() string {
	return GenerateSalt(32)
}

func GenerateDeviceId(prefix string) string {
	if !strings.HasSuffix(prefix, ":") {
		prefix = prefix + ":"
	}
	var flags = random.O_DIGIT | random.O_LOWER_CHAR | random.O_UPPER_CHAR
	return prefix + random.String(31-len(prefix), random.CryptoSource, flags) + "@" +
		random.String(32, random.CryptoSource, flags)
}

func EncryptPassword(password, salt string) string {
	return Sha512_256(salt + password)
}

func GenerateToken() string {
	return GenerateTokenWithLength(64)
}

func GenerateTokenWithLength(length int) string {
	return random.String(length, random.CryptoSource, random.O_DIGIT|random.O_LOWER_CHAR|random.O_UPPER_CHAR)
}

func MD5(src string) string        { return fmt.Sprintf("%x", md5.Sum([]byte(src))) }
func Sha1(src string) string       { return fmt.Sprintf("%x", sha1.Sum([]byte(src))) }
func Sha256(src string) string     { return fmt.Sprintf("%x", sha256.Sum256([]byte(src))) }
func Sha512(src string) string     { return fmt.Sprintf("%x", sha512.Sum512([]byte(src))) }
func Sha512_256(src string) string { return fmt.Sprintf("%x", sha512.Sum512_256([]byte(src))) }

func MD5Upper(src string) string        { return fmt.Sprintf("%X", md5.Sum([]byte(src))) }
func Sha1Upper(src string) string       { return fmt.Sprintf("%X", sha1.Sum([]byte(src))) }
func Sha256Upper(src string) string     { return fmt.Sprintf("%X", sha256.Sum256([]byte(src))) }
func Sha512Upper(src string) string     { return fmt.Sprintf("%X", sha512.Sum512([]byte(src))) }
func Sha512_256Upper(src string) string { return fmt.Sprintf("%X", sha512.Sum512_256([]byte(src))) }
