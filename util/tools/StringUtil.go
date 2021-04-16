package tools

import (
	"strings"
)

func StrGetMiddle(data, index, tail string) string {
	iIndex := strings.Index(data, index)
	if iIndex != -1 {
		iIndex += len(index)
	} else {
		return ""
	}
	data = string([]byte(data)[iIndex:])

	iEnd := strings.Index(data, tail)
	data = string([]byte(data)[:iEnd])
	return data
}
