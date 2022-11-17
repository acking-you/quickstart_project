package util

type Set[T comparable] map[T]struct{}

func (s *Set[T]) Insert(v T) {
	(*s)[v] = struct{}{}
}

func (s *Set[T]) Contains(v T) bool {
	if _, ok := (*s)[v]; ok {
		return true
	}
	return false
}
