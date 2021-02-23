package anagram

// Storage stores
type Storage struct {
	parts   [][]Part
	options *Options
	anagram *Anagram
}

// InitStorage initializes a storage
func InitStorage(a *Anagram, options *Options) *Storage {
	var s Storage
	s.options = options
	s.parts = make([][]Part, a.Length)

	for index := range s.parts {
		s.parts[index] = make([]Part, 0, 100)
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
