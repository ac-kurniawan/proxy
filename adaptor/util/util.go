package util

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/ac-kurniawan/proxy/core"
	"github.com/ac-kurniawan/proxy/library"
	"golang.org/x/crypto/pbkdf2"
)

type Util struct {
	library.AppTrace
	library.AppLog
	Secret string
}

// EncryptPassword implements core.IUtil.
func (u *Util) EncryptPassword(password string) string {
	encrypted := pbkdf2.Key([]byte(password), []byte(u.Secret), 1000, 32, sha1.New)
	return hex.EncodeToString(encrypted)
}

// GenerateJWT implements core.IUtil.
func (u *Util) GenerateJWT(attribute map[string]interface{}, exp int64) string {
	panic("unimplemented")
}

// GetExpFromToken implements core.IUtil.
func (u *Util) GetExpFromToken(token string) int64 {
	panic("unimplemented")
}

func NewUtil(module Util) core.IUtil {
	return &module
}
