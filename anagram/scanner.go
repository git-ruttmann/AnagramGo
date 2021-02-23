package anagram

// Scanner scans words
type Scanner struct {
	options  *Options
	anagram  *Anagram
	storage  *Storage
	results  []Part
	reporter func(string)
}

// Initialize sets up the scanner
func (s *Scanner) Initialize(anagram *Anagram, options *Options, reporter func(string)) {
	s.anagram = anagram
	s.options = options
	s.reporter = reporter
	s.results = make([]Part, 0, 160000)
	s.storage = InitStorage(anagram, options)
}

// ProcessWord processes a single word
func (s *Scanner) ProcessWord(text string) {
	if len(text) <= s.options.MinimumLength {
		return
	}

	s.results = s.results[0:0]
	// s.results = make([]Part, 100)
	word := s.anagram.Combine(text)
	if word == nil {
		return
	}

	minLength := s.options.MinimumLength
	for i := word.Length + minLength + 1; i < s.anagram.Length-minLength; i++ {
		// for i := word.Length + 1; i < s.anagram.Length; i++ {
		lengthCluster := s.storage.parts[i]
		for j := 0; j < len(lengthCluster); j++ {
			var target Part
			if word.Combine(&lengthCluster[j], &target) {
				s.results = append(s.results, target)
				s.results[len(s.results)-1].Remaining = target.Remaining
			}
		}
	}

	lengthCluster := s.storage.parts[word.Length]
	for j := 0; j < len(lengthCluster); j++ {
		if word.IsComplete(&lengthCluster[j]) {
			s.reporter(lengthCluster[j].text + " " + word.text)
		}
	}

	s.storage.addResults(s.results)
	//	fmt.Println(text, ": ", len(s.results))
	var part Part
	word.ToPart(&part)
	s.storage.Add(&part)
}

func (s *Storage) addResults(results []Part) {
	for _, result := range results {
		s.Add(&result)
	}
}
