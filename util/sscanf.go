package util

import (
	"errors"
	"strings"
)

// Sscanf 以'$'作为format串中字符串的占位符，将str里的内容按照format的形式，将其中的$占位的内容以字符串的形式解析到args参数中
func Sscanf(str, format string, args ...*string) error {
	if len(str) < len(format) {
		return errors.New("len(str) < len(format)")
	}
	strIndex := 0
	currentArgsIndex := 0
	for i := 0; i < len(format); i++ {
		//若不是占位符，则进行单纯的字符匹配
		if format[i] != '$' {
			if strIndex < len(str) && format[i] == str[strIndex] {
				strIndex++
				continue
			}
			return errors.New("str not valid format string")
		}
		//若为占位符，则需要分两种情况：1.占位符有右边界 2.占位符无右边界
		if currentArgsIndex >= len(args) {
			return errors.New("args index overstep the boundary")
		}
		var ret string
		//情况2
		if i == len(format)-1 {
			ret = str[strIndex:]
		} else { //情况1
			tStr := str[strIndex:]
			idx := strings.IndexByte(tStr, format[i+1])
			if idx == -1 {
				return errors.New("str not fit the \"format\"")
			}
			ret = tStr[:idx]
			strIndex += idx
		}
		*args[currentArgsIndex] = ret
		currentArgsIndex++
	}
	return nil
}
