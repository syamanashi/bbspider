package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

// ThreadScrape scrapes a forum on this style bulletin board (bb type t.b.d.)
// TODO: Determiine if ThreadOffset needs to be a property of each thread - or if it can continue to be used for both forums and threads.
func ThreadScrape(forum *Forum) {

	var maxThreadCount int
	if Conf.Maxthreads > 0 {
		maxThreadCount = Conf.Maxthreads
	} else {
		maxThreadCount = forum.PageCount + 1
	}

	// Loop through each thread (or up to maxThreadCount)
	// for i, thread := range forum.Threads {
	for i := 0; i < maxThreadCount; i++ {
		thread := forum.Threads[i]

		// A thread is always at least one page long
		thread.PageCount = 1

		fmt.Printf("\n\n%d) ThreadID: %d\n", i+1, thread.ThreadID)
		// fmt.Printf("#%v\n", thread)

		// Get the root thread page to determine thread pageCount
		threadCompleteURL := forum.RootURL + thread.ThreadURL
		fmt.Printf("threadCompleteURL: %s\n", threadCompleteURL)

		response, err := http.Get(threadCompleteURL)
		checkError(err)
		defer response.Body.Close()

		doc, err := goquery.NewDocumentFromReader(io.Reader(response.Body))
		checkError(err)

		// PageCount defined by last page link if it exists.  If not present, leave thread.PageCount to be 1.
		threadPageCount, err := strconv.Atoi(doc.Find("." + forum.PageInactiveClass).Last().Text())
		if threadPageCount > 0 {
			thread.PageCount = threadPageCount
		}

		// Output threadPageCount
		fmt.Printf("Thread page count: %d\n", thread.PageCount)

		// Create threadID folder where the HTML for each page within the thread gets stored
		thread.ThreadPathComplete = filepath.Join(Conf.ThreadPath, strconv.Itoa(thread.ThreadID))
		err = os.MkdirAll(thread.ThreadPathComplete, os.ModePerm)
		if err != nil {
			log.Println("Error creating thread directory for ThreadID:", thread.ThreadID)
			log.Println(err)
			return
		}

		// Loop through all thread pages and store each in a seperate HTML file within the thread's directory
		// TODO: Move the following repeated logic (forum and thread) to another func:
		var maxPageCount int
		if Conf.Maxthreads > 0 {
			maxPageCount = Conf.Maxthreads + 1
		} else {
			maxPageCount = forum.PageCount + 1
		}

		amt := time.Duration(Conf.Throttle)

		// Loop the the thread's pages
		for i := 1; i < maxPageCount; i++ {
			threadPageScrape(&thread, i)
			time.Sleep(time.Millisecond * amt)
		}

	}
}

func threadPageScrape(thread *Thread, pageNum int) {

	offset := (pageNum - 1) * forum.ThreadOffset
	url := forum.RootURL + thread.ThreadURL + "?offset=" + strconv.Itoa(offset)

	fmt.Printf("Scraping thread page URL: %s\n", url)

	response, err := http.Get(url)
	checkError(err)

	filePath := thread.ThreadPathComplete + "/offset_" + strconv.Itoa(offset) + ".html"

	file, err := os.Create(filePath)
	checkError(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	checkError(err)
	defer response.Body.Close()

}
