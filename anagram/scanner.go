package anagram

const threadCount = 8

// Scanner scans words
type Scanner struct {
	wordAnalyzer wordAnalyzer
	wordCombiner wordCombiner
}

type wordAnalyzer struct {
	options *Options
	anagram *Anagram
	texts   chan string
	words   chan Word
}

type wordCombiner struct {
	storage         *Storage
	options         *Options
	anagram         *Anagram
	threads         map[int]*scanThread
	stageTwoThreads map[int]*scanThread
	processedWordID int
	stageOne        chan int
	words           chan Word
	endChannel      chan bool
	endRequested    bool
	reporter        func(string)
}

// NewScanner initalizes a new scanner
func NewScanner(anagram *Anagram, options *Options, reporter func(string)) (s *Scanner) {
	s = &Scanner{
		wordCombiner: wordCombiner{
			anagram:  anagram,
			options:  options,
			reporter: reporter,
		},
		wordAnalyzer: wordAnalyzer{
			anagram: anagram,
			options: options,
		},
	}

	s.wordCombiner.storage = InitStorage(anagram, options)
	s.wordCombiner.endChannel = make(chan bool, 1)

	s.wordAnalyzer.texts = make(chan string, 20)
	s.wordAnalyzer.words = make(chan Word, 20)

	s.wordCombiner.words = s.wordAnalyzer.words

	go s.wordCombiner.combineAcceptedWords()
	go s.wordAnalyzer.scanWords()
	return
}

func (w *wordAnalyzer) scanWords() {
	for text := range w.texts {
		if len(text) > w.options.MinimumLength {
			if word := w.anagram.Combine(text); word != nil {
				w.words <- *word
			}
		}
	}

	close(w.words)
}

func (w *wordCombiner) combineAcceptedWords() {
	wordID := 0
	w.processedWordID = 1
	generateID := func() int {
		wordID++
		return wordID
	}

	w.stageOne = make(chan int, 1)
	w.threads = make(map[int]*scanThread)
	w.stageTwoThreads = make(map[int]*scanThread)
	wordChannel := w.words

	for {
		select {
		case id := <-w.stageOne:
			w.stageTwoThreads[id] = w.threads[id]

		case word, ok := <-wordChannel:
			if ok {
				w.launchNewThread(&word, generateID())
			} else {
				w.words = nil
				wordChannel = nil
			}
		}

		for {
			if t, ok := w.stageTwoThreads[w.processedWordID]; ok {
				w.runStageTwo(t, w.processedWordID)
				w.processedWordID++
			} else {
				break
			}
		}

		if len(w.threads) > 2000 {
			wordChannel = nil
		} else {
			wordChannel = w.words
		}

		if w.endRequested && len(w.threads) == 0 && w.words == nil {
			w.endChannel <- true
			break
		}
	}
}

func (w *wordCombiner) launchNewThread(word *Word, id int) {
	t := newScanThread(w.storage, word)
	w.threads[id] = t
	t.StageOne()
	go t.Scan(id, w.stageOne)
}

func (w *wordCombiner) runStageTwo(t *scanThread, id int) {
	delete(w.stageTwoThreads, id)
	delete(w.threads, id)
	t.StageTwo()

	// lots to do - start scanning for each thread in completed stage
	if len(t.complete.parts) > 1000 {
		w.relaunchStageOneForAllThreads()
	}

	t.Scan2(id)
	t.HandleCompletionInSynchronizationThread(w.reporter)
}

func (w *wordCombiner) relaunchStageOneForAllThreads() {
	oldMap := w.stageTwoThreads
	w.stageTwoThreads = make(map[int]*scanThread)
	for id, t := range oldMap {
		t.StageTwo()
		go t.Scan(id, w.stageOne)
	}
}

// ProcessWord processes a single word
func (s *Scanner) ProcessWord(text string) {
	s.wordAnalyzer.texts <- text
}

// Final closed the scanner, call after the last word
func (s *Scanner) Final() {
	close(s.wordAnalyzer.texts)
	s.wordCombiner.endRequested = true
	<-s.wordCombiner.endChannel
}
