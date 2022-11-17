package util

import "strings"

func IsUpper(ch uint8) bool {
	return ch&32 == 0
}

func ToUpper(ch uint8) uint8 {
	return uint8(int8(ch) & int8(-33))
}

func ToLower(ch uint8) uint8 {
	return ch | 32
}

func SnackCase2PascalCase(name string) string {
	ret := SnackCase2CamelCase(name)
	if len(ret) == 0 {
		return ret
	}
	return strings.ToUpper(ret[0:1]) + ret[1:]
}

func SnackCase2CamelCase(name string) string {
	var ret []uint8
	flag := false
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			flag = true
			continue
		}
		ch := name[i]
		if flag {
			ch = ToUpper(ch) //字符转大写
			flag = false
		}
		ret = append(ret, ch)
	}
	if len(ret) == 0 {
		return ""
	}
	return string(ret)
}
func PascalCase2SnackCase(code string) string {
	return CamelCase2SnackCase(strings.ToLower(code[0:1]) + code[1:])
}

func CamelCase2SnackCase(code string) string {
	var ret []uint8
	for i := 0; i < len(code); i++ {
		if IsUpper(code[i]) {
			ret = append(ret, '_')
		}
		ret = append(ret, ToLower(code[i]))
	}
	if len(ret) == 0 {
		return ""
	}
	return string(ret)
}
