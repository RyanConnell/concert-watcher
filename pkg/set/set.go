package set

type Set[T comparable] struct {
	values map[T]struct{}
}

func New[T comparable](data ...T) *Set[T] {
	values := make(map[T]struct{}, len(data))
	for _, d := range data {
		values[d] = struct{}{}
	}
	return &Set[T]{values: values}
}

func (s *Set[T]) Add(values ...T) {
	for _, value := range values {
		s.values[value] = struct{}{}
	}
}

func (s *Set[T]) Values() []T {
	var values []T
	for v := range s.values {
		values = append(values, v)
	}
	return values
}

func (s *Set[T]) Contains(value T) bool {
	_, ok := s.values[value]
	return ok
}
