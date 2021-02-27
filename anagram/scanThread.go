package anagram

type scanThread struct {
	partial  partsScan
	complete completionScan
	storage  *Storage
	word     *Word
}

type partsScan struct {
	parts      [][]Part
	upperBound []int
	newParts   []Part
}

type completionScan struct {
	parts      []Part
	upperBound int
	results    []string
}

func newScanThread(storage *Storage, word *Word) (t *scanThread) {
	t = &scanThread{
		storage: storage,
		word:    word,
	}

	t.partial.parts = make([][]Part, 0, storage.anagram.Length)
	t.partial.upperBound = make([]int, 0, storage.anagram.Length)
	t.partial.newParts = make([]Part, 0, 2048)
	t.complete.results = make([]string, 0, 256)
	return
}

func (t *scanThread) HandleCompletionInSynchronizationThread(reporter func(string)) {
	for _, result := range t.partial.newParts {
		t.storage.Add(&result)
	}

	var part Part
	t.word.ToPart(&part)
	t.storage.Add(&part)

	for _, result := range t.complete.results {
		reporter(result)
	}
}

func (t *scanThread) StageOne() {
	t.partial.parts = t.partial.parts[:0]
	minLength := t.storage.options.MinimumLength

	t.partial.newParts = t.partial.newParts[:0]
	for i := t.word.Length + minLength + 1; i < t.storage.anagram.Length-minLength; i++ {
		t.partial.parts = append(t.partial.parts, t.storage.parts[i])
		t.partial.upperBound = append(t.partial.upperBound, len(t.storage.parts[i]))
	}

	t.complete.parts = t.storage.parts[t.word.Length]
	t.complete.upperBound = len(t.storage.parts[t.word.Length])
	t.complete.results = t.complete.results[:0]
}

// scan the unscanned parts (start at the last upper bound)
func (t *scanThread) StageTwo() {
	minLength := t.storage.options.MinimumLength

	offset := t.word.Length + minLength + 1
	for i := range t.partial.parts {
		t.partial.parts[i] = t.storage.parts[i+offset][t.partial.upperBound[i]:]
		t.partial.upperBound[i] += len(t.partial.parts[i])
	}

	t.complete.parts = t.storage.parts[t.word.Length][t.complete.upperBound:]
	t.complete.upperBound += len(t.complete.parts)
}

func (t *scanThread) Scan(id int, scanComplete chan int) {
	t.partial.Scan(t.word)
	t.complete.Scan(t.word)
	scanComplete <- id
}

func (t *scanThread) Scan2(id int) {
	total := 0
	for _, part := range t.partial.parts {
		total += len(part)
	}

	t.partial.Scan(t.word)
	t.complete.Scan(t.word)
}

func (t *partsScan) Scan(word *Word) {
	for _, parts := range t.parts {
		for i := 0; i < len(parts); i++ {
			part := &parts[i]
			if (part.DoNotUseMask & word.UsedMask) != 0 {
				continue
			}

			var target Part
			if word.Combine(part, &target) {
				t.newParts = append(t.newParts, target)
				t.newParts[len(t.newParts)-1].Remaining = target.Remaining
			}
		}
	}
}

func (t *completionScan) Scan(word *Word) {
	for i := 0; i < len(t.parts); i++ {
		part := &t.parts[i]
		if (part.DoNotUseMask & word.UsedMask) != 0 {
			continue
		}

		if word.IsComplete(part) {
			t.results = append(t.results, part.text+" "+word.text)
		}
	}
}
