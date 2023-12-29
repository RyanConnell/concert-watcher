package set

type Set[T comparable] struct {
	values map[T]struct{}
}

func NewSet[T comparable](data ...T) *Set[T] {
	values := make(map[T]struct{}, len(data))
	for _, d := range data {
		values[d] = struct{}{}
	}
	return &Set[T]{values: values}
}

func (s *Set[T]) Add(value T) {
	s.values[value] = struct{}{}
}

func (s *Set[T]) Values() []T {
	var values []T
	for v := range s.values {
		values = append(values, v)
	}
	return values
}
