package anagramtests

import (
	"anagram"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the deep serialization of anagram.Part over a channel
func TestChannelCommunication(t *testing.T) {
	a := anagram.InitizalizeAnagram("Best Secret")
	part := createPart(&a, "test")
	c := make(chan anagram.Part, 2)

	c <- part
	part = createPart(&a, "best")
	c <- part

	j := <-c
	assert.NotEqual(t, part.DoNotUseMask, j.DoNotUseMask)

	k := <-c
	assert.Equal(t, part.DoNotUseMask, k.DoNotUseMask)
	assert.Equal(t, part.Remaining[0], k.Remaining[0])
	assert.Equal(t, part.Remaining[1], k.Remaining[1])
	assert.Equal(t, part.Remaining[2], k.Remaining[2])
	assert.Equal(t, part.Remaining[3], k.Remaining[3])
}

func TestAnalysis(t *testing.T) {
	a := anagram.InitizalizeAnagram("Best Secret")

	assert.Equal(t, "Best Secret", a.Text)
	assert.Equal(t, 10, a.Length)
}

func TestFirstWord(t *testing.T) {
	a := anagram.InitizalizeAnagram("Best Secret")
	wordBest := a.Combine("Best")

	assert.Equal(t, 4, wordBest.Length)
	assert.Equal(t, 6, wordBest.RestLength)
}

func TestWordCombination(t *testing.T) {
	a := anagram.InitizalizeAnagram("Best Secret")
	partBest := createPart(&a, "Best")
	wordSecret := a.Combine("Secret")

	assert.True(t, wordSecret.IsComplete(&partBest))
}

func TestWordMask(t *testing.T) {
	a := anagram.InitizalizeAnagram("Best Secret")
	partBest := createPart(&a, "Best")
	wordSecret := a.Combine("Secret")

	assert.Equal(t, uint32(0), partBest.DoNotUseMask&wordSecret.UsedMask)
}

func TestWordInvalidMask(t *testing.T) {
	a := anagram.InitizalizeAnagram("Best Secret")
	partBest := createPart(&a, "Best")
	wordSecret := a.Combine("Becret")

	assert.NotEqual(t, uint32(0), partBest.DoNotUseMask&wordSecret.UsedMask)
}

func TestScannerBestSecretDouble(t *testing.T) {
	scanner := createTestScanner("Best Secret")
	scanner.scanner.ProcessWord("Best")
	scanner.scanner.ProcessWord("Secret")

	assert.Equal(t, 1, len(scanner.results))
}

func TestScannerBestSecretTriple(t *testing.T) {
	scanner := createTestScanner("Best Secret")
	scanner.scanner.ProcessWord("bet")
	scanner.scanner.ProcessWord("erst")
	scanner.scanner.ProcessWord("sec")

	assert.Equal(t, 1, len(scanner.results))
}

func TestScannerBestSecretAll(t *testing.T) {
	scanner := createTestScanner("Best Secret")
	stream := scanner.scanner
	stream.ProcessWord("beet")
	stream.ProcessWord("crests")
	stream.ProcessWord("beets")
	stream.ProcessWord("beret")
	stream.ProcessWord("berets")
	stream.ProcessWord("beset")
	stream.ProcessWord("best")
	stream.ProcessWord("bests")
	stream.ProcessWord("bet")
	stream.ProcessWord("bets")
	stream.ProcessWord("better")
	stream.ProcessWord("betters")
	stream.ProcessWord("cess")
	stream.ProcessWord("crest")
	stream.ProcessWord("crete")
	stream.ProcessWord("erect")
	stream.ProcessWord("erects")
	stream.ProcessWord("erst")
	stream.ProcessWord("rest")
	stream.ProcessWord("sec")
	stream.ProcessWord("secret")
	stream.ProcessWord("secrets")
	stream.ProcessWord("sect")
	stream.ProcessWord("sects")

	assert.Equal(t, 16, len(scanner.results))
}

func createPart(a *anagram.Anagram, text string) anagram.Part {
	var part anagram.Part
	word := a.Combine(text)
	word.ToPart(&part)
	return part
}

type testScanner struct {
	scanner anagram.Scanner
	results []string
}

func createTestScanner(anagramText string) *testScanner {
	a := anagram.InitizalizeAnagram(anagramText)
	o := anagram.Options{MinimumLength: 2}
	test := testScanner{scanner: anagram.Scanner{}, results: make([]string, 0)}

	reporter := func(text string) {
		test.results = append(test.results, text)
	}

	test.scanner.Initialize(&a, &o, reporter)
	return &test
}
