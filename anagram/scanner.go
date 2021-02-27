package anagram

import "fmt"

const threadCount = 8

// Scanner scans words
type Scanner struct {
	options      *Options
	anagram      *Anagram
	storage      *Storage
	wordAnalyzer wordAnalyzer
	endChannel   chan bool
	endRequested bool
	reporter     func(string)
}

type wordAnalyzer struct {
	texts chan string
	words chan Word
}

// NewScanner initalizes a new scanner
func NewScanner(anagram *Anagram, options *Options, reporter func(string)) (s *Scanner) {
	s = &Scanner{
		anagram:  anagram,
		options:  options,
		reporter: reporter,
	}

	s.storage = InitStorage(anagram, options)
	s.endChannel = make(chan bool, 1)
	s.wordAnalyzer.texts = make(chan string, 20)
	s.wordAnalyzer.words = make(chan Word, 20)

	go s.processAcceptedWords()
	go s.scanWords()
	return
}

func (s *Scanner) scanWords() {
	for text := range s.wordAnalyzer.texts {
		if len(text) > s.options.MinimumLength {
			if word := s.anagram.Combine(text); word != nil {
				s.wordAnalyzer.words <- *word
			}
		}
	}

	close(s.wordAnalyzer.words)
}

func (s *Scanner) processAcceptedWords() {
	wordID := 0
	nextProcessedWordID := 1
	generateID := func() int {
		wordID++
		return wordID
	}

	stageOne := make(chan int, 1)
	stageTwo := make(chan int, 1)
	threads := make(map[int]*scanThread)
	completedStageTwo := make(map[int]*scanThread)

	for {
		select {
		case id := <-stageTwo:
			t := threads[id]
			delete(threads, id)
			t.HandleCompletionInSynchronizationThread(s.reporter)

		case id := <-stageOne:
			t := threads[id]
			completedStageTwo[id] = t

		case word, ok := <-s.wordAnalyzer.words:
			if ok {
				id := generateID()
				t := newScanThread(s.storage, &word)
				threads[id] = t
				t.StageOne()
				go t.Scan(id, stageOne)
			} else {
				s.wordAnalyzer.words = nil
			}
		}

		count := 0
		for {
			if t, ok := completedStageTwo[nextProcessedWordID]; ok {
				count++
				id := nextProcessedWordID
				nextProcessedWordID++

				delete(completedStageTwo, id)
				delete(threads, id)
				t.StageTwo()

				// lots to do - start scanning for each thread in completed stage
				if len(t.complete.parts) > 1000 {
					for nid, nt := range completedStageTwo {
						nt.StageTwo()
						go nt.Scan(nid, stageOne)
					}

					completedStageTwo = make(map[int]*scanThread)
				}

				t.Scan2(id)
				t.HandleCompletionInSynchronizationThread(s.reporter)
			} else {
				if count > 0 {
					fmt.Println(count, " ", len(threads))
				}
				break
			}
		}

		if s.endRequested && len(threads) == 0 && s.wordAnalyzer.words == nil {
			s.endChannel <- true
			break
		}
	}
}

// ProcessWord processes a single word
func (s *Scanner) ProcessWord(text string) {
	s.wordAnalyzer.texts <- text
}

// Final closed the scanner, call after the last word
func (s *Scanner) Final() {
	close(s.wordAnalyzer.texts)
	s.endRequested = true
	<-s.endChannel
}
