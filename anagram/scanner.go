package anagram

// Scanner scans words
type Scanner struct {
	options  *Options
	anagram  *Anagram
	storage  *Storage
	results  [][]Part
	reporter func(string)
}

// Initialize sets up the scanner
func (s *Scanner) Initialize(anagram *Anagram, options *Options, reporter func(string)) {
	s.anagram = anagram
	s.options = options
	s.reporter = reporter
	s.results = make([][]Part, anagram.Length)
	for i := 0; i < anagram.Length; i++ {
		s.results[i] = make([]Part, 0, 2048)
	}

	s.storage = InitStorage(anagram, options)
}

func processSlice(parts []Part, results *[]Part, word *Word, completedChannel chan int, channelID int) {
	for j := 0; j < len(parts); j++ {
		part := &parts[j]
		if (part.DoNotUseMask & word.UsedMask) != 0 {
			continue
		}

		var target Part
		if word.Combine(part, &target) {
			*results = append(*results, target)
			(*results)[len(*results)-1].Remaining = target.Remaining
		}
	}

	completedChannel <- channelID
}

// ProcessWord processes a single word
func (s *Scanner) ProcessWord(text string) {
	if len(text) <= s.options.MinimumLength {
		return
	}

	for i := 0; i < len(s.results); i++ {
		s.results[i] = s.results[i][0:0]
	}

	// s.results = make([]Part, 100)
	word := s.anagram.Combine(text)
	if word == nil {
		return
	}

	minLength := s.options.MinimumLength
	completedChannel := make(chan int, 100)
	channelCount := 0
	for i := word.Length + minLength + 1; i < s.anagram.Length-minLength; i++ {
		go processSlice(s.storage.parts[i], &s.results[i], word, completedChannel, i)
		channelCount++
	}

	lengthCluster := s.storage.parts[word.Length]
	for j := 0; j < len(lengthCluster); j++ {
		if word.IsComplete(&lengthCluster[j]) {
			s.reporter(lengthCluster[j].text + " " + word.text)
		}
	}

	for ; channelCount > 0; channelCount-- {
		i := <-completedChannel
		s.storage.addResults(s.results[i])
	}

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
