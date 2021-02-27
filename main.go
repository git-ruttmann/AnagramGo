package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"anagram"
)

func main() {
	a := anagram.InitizalizeAnagram("Best Secret Aschheim")
	var options anagram.Options

	options.MinimumLength = 2
	options.PrintEntries = false

	resultCount := uint64(0)
	processor := anagram.NewScanner(&a, &options, func(text string) {
		resultCount++
		if options.PrintEntries {
			fmt.Println(text)
		}
	})

	file, err := os.Open("../wordlist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		processor.ProcessWord(scanner.Text())
	}
	processor.Final()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Count: ", resultCount)
}
