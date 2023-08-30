package util

import (
	"strings"

	"github.com/chenjie199234/Corelib/util/common"
)

func MaskTel(origin string) string {
	tmp := make([]byte, len(origin))
	for i := range origin {
		if i%2 == 0 {
			tmp[i] = origin[i]
		} else {
			tmp[i] = '*'
		}
	}
	return common.Byte2str(tmp)
}
func MaskEmail(origin string) string {
	index := strings.Index(origin, "@")
	return strings.Repeat("*", index) + origin[index:]
}
func MaskIDCard(origin string) string {
	tmp := make([]byte, len(origin))
	for i := range origin {
		if i%2 == 0 {
			tmp[i] = origin[i]
		} else {
			tmp[i] = '*'
		}
	}
	return common.Byte2str(tmp)
}
