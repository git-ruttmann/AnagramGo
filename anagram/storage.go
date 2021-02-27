package anagram

// Storage stores
type Storage struct {
	parts   [][]Part
	options *Options
	anagram *Anagram
}

func estimateCapacity(index int) int {
	switch {
	case index < 3:
		return 0
	case index < 9:
		return 1000000
	case index < 12:
		return 50000
	default:
		return 10000
	}
}

// InitStorage initializes a storage
func InitStorage(a *Anagram, options *Options) *Storage {
	var s Storage
	s.options = options
	s.anagram = a
	s.parts = make([][]Part, a.Length)

	for index := range s.parts {
		s.parts[index] = make([]Part, 0, estimateCapacity(index))
	}

	return &s
}

// Add adds the part to the storage
func (s *Storage) Add(part *Part) {
	a := s.parts[part.RestLength]
	a = append(a, *part)
	a[len(a)-1].Remaining = part.Remaining
	s.parts[part.RestLength] = a
}
