package util

import "strings"

type StrHandleByChain struct {
	Str string
}

func (s *StrHandleByChain) ReplaceAll(old, new string) *StrHandleByChain {
	s.Str = strings.ReplaceAll(s.Str, old, new)
	return s
}
