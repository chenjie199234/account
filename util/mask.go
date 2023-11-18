package util

import (
	"github.com/chenjie199234/Corelib/util/common"
)

func MaskTel(origin string) string {
	if len(origin) == 0 {
		return ""
	}
	tmp := make([]byte, len(origin))
	for i := range origin {
		if i%2 == 0 {
			tmp[i] = origin[i]
		} else {
			tmp[i] = '*'
		}
	}
	return common.BTS(tmp)
}
func MaskEmail(origin string) string {
	if len(origin) == 0 {
		return ""
	}
	tmp := make([]byte, len(origin))
	index := -1
	for i := range origin {
		if origin[i] == '@' {
			index = i
			break
		}
		if i%2 == 0 {
			tmp[i] = origin[i]
		} else {
			tmp[i] = '*'
		}
	}
	if index == -1 || index == 0 {
		return ""
	}
	for i := index; i < len(origin); i++ {
		tmp[i] = origin[i]
	}
	return common.BTS(tmp)
}
func MaskIDCard(origin string) string {
	if len(origin) == 0 {
		return ""
	}
	tmp := make([]byte, len(origin))
	for i := range origin {
		if i%2 == 0 {
			tmp[i] = origin[i]
		} else {
			tmp[i] = '*'
		}
	}
	return common.BTS(tmp)
}
