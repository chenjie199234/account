package util

import (
	"context"
	"math/rand"

	"github.com/chenjie199234/Corelib/util/common"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#*" //letters' length is 64 and 2^6=64

func MakeRandCode() string {
	b := make([]byte, 6)
	r := rand.Uint64()
	for i := range b {
		b[i] = letters[(r<<(i*6))>>58]
	}
	return common.Byte2str(b)
}

// TODO
func SendEmailCode(ctx context.Context, email, code, action string) error {
	return nil
}

// TODO
func SendTelCode(ctx context.Context, tel, code, action string) error {
	return nil
}
