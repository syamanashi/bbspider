package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ForumScrape scrapes a forum on this style bulletin board (bb type t.b.d.)
func ForumScrape(forum *Forum) {

	response, err := http.Get(forum.RootURL + forum.ForumURL)
	checkError(err)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(io.Reader(response.Body))
	checkError(err)

	// Forum PageCount 0 vs 1 logic to forum.PageCount for when the pagination links are missing.
	forum.PageCount = 1

	forum.Title = doc.Find("h1").First().Text()
	forumPageCount, err := strconv.Atoi(doc.Find("." + forum.PageInactiveClass).Last().Text())
	checkError(err)

	if forumPageCount > 0 {
		forum.PageCount = forumPageCount
	}

	// Loop through and scrape forum pages.
	var maxPageCount int
	if Conf.Maxthreads > 0 {
		maxPageCount = Conf.Maxthreads + 1
	} else {
		maxPageCount = forum.PageCount + 1
	}

	amt := time.Duration(Conf.Throttle)

	for i := 1; i < maxPageCount; i++ {
		indexPageScrape(forum, i, Conf.IndexPath)
		time.Sleep(time.Millisecond * amt)
	}
}

func indexPageScrape(forum *Forum, pageNum int, path string) {

	offset := (pageNum - 1) * forum.ThreadOffset
	url := forum.RootURL + forum.ForumURL + "?offset=" + strconv.Itoa(offset)

	println("\n\nFORUM: ", forum.Title, "page", pageNum, "of", forum.PageCount, "pages (offset:", offset, ")")
	fmt.Printf("Scraping forum index page URL: %s\n", url)

	response, err := http.Get(url)
	checkError(err)

	filePath := path + "/offset_" + strconv.Itoa(offset) + ".html"

	file, err := os.Create(filePath)
	checkError(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	checkError(err)
	defer response.Body.Close()

	// Get Thread Stats for each forum index page.  Open the recently created HTML file to read from it with goquery.
	file, err = os.Open(filePath)
	checkError(err)
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	checkError(err)

	var threads []Thread

	// Find the forum page threads
	doc.Find(forum.ThreadsTableSelector).Each(func(i int, selection *goquery.Selection) {
		// For each item found, get the thread title, author, URL, lastReply, replyCount
		threadTitle := selection.Find(forum.ThreadURLSelector).Text()
		threadTitleWithAuthor := selection.Find(forum.ThreadTitleSelector).Text()
		threadURL, ok := selection.Find(forum.ThreadURLSelector).Attr("href")
		if ok == false {
			threadURL = "empty"
		}
		threadLastReplyTimestamp := selection.Find(forum.ThreadLastReplySelector).Text()
		threadReplyCountString := selection.Find(forum.ThreadReplyCountSelector).Text()
		threadReplyCount, err := strconv.Atoi(threadReplyCountString)
		if err != nil {
			threadReplyCount = 0
		}

		// Get the thread creator, (though this may be better suited with logic extracting the username from the first post in thread).
		threadTitleSlice := WordSlice(threadTitleWithAuthor)
		threadCreator := threadTitleSlice[len(threadTitleSlice)-1]

		threadID, err := strconv.Atoi(WordSlice(threadURL)[2])
		if err != nil {
			threadID = 0
		}

		thread := Thread{
			ThreadID:                 threadID,
			ThreadTitle:              threadTitle,
			ThreadURL:                threadURL,
			ThreadLastReplyTimestamp: threadLastReplyTimestamp,
			ThreadReplyCount:         threadReplyCount,
			ThreadCreator:            threadCreator,
		}

		threads = append(threads, thread)
		forum.Threads = append(forum.Threads, thread)
	})

	if Conf.PrintThreads == true {
		PrintThreads(threads, pageNum, forum.ThreadOffset)
	}
}
