package main

import (
	"fmt"
	"regexp"
)

// WordSlice returns a slice of words contained in passed-in string.
func WordSlice(text string) []string {
	words := regexp.MustCompile("\\w+")
	return words.FindAllString(text, -1)
}

// PrintThreads prints the thread title, last reply from threads in a single forum.
func PrintThreads(s []Thread, pageNum int, offset int) {
	terminalThreadSpacer := "  "

	fmt.Printf("Thread count: %d\n", len(s))
	fmt.Printf("PageNum: %d\n", pageNum)

	for i, t := range s {
		fmt.Printf("\n\n%d) %s\n", (i+1)+((pageNum-1)*offset), t.ThreadTitle)
		fmt.Printf("%s Thread Creator: %s\n", terminalThreadSpacer, t.ThreadCreator)
		fmt.Printf("%s Last Reply: %s\n", terminalThreadSpacer, t.ThreadLastReplyTimestamp)
		fmt.Printf("%s Reply Count: %d\n", terminalThreadSpacer, t.ThreadReplyCount)
		fmt.Printf("%s Thread URL: %s\n", terminalThreadSpacer, t.ThreadURL)
		fmt.Printf("%s ThreadID: %d\n", terminalThreadSpacer, t.ThreadID)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
