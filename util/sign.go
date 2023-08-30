package util

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"

	"github.com/chenjie199234/account/ecode"
)

func SignCheck(secret, HEXnoncesign string) (e error) {
	noncesign, e := hex.DecodeString(HEXnoncesign)
	if e != nil {
		return ecode.ErrDataBroken
	}
	if len(noncesign) < 8+64 {
		return ecode.ErrDataBroken
	}
	oldsign := make([]byte, 64)
	copy(oldsign, noncesign[len(noncesign)-64:])
	newsign := sha512.Sum512(append(noncesign[:len(noncesign)-64], secret...))
	if !bytes.Equal(oldsign, newsign[:]) {
		return ecode.ErrSignCheckFailed
	}
	return nil
}
func SignMake(secret string, nonce []byte) (HEXnoncesign string) {
	tmp := make([]byte, 8+len(nonce))
	binary.BigEndian.PutUint64(tmp, uint64(len(nonce)))
	copy(tmp[8:], nonce)
	newsign := sha512.Sum512(append(tmp, secret...))
	return hex.EncodeToString(append(tmp, newsign[:]...))
}
