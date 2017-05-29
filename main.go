package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/goinggo/tracelog"
)

// Forum type stores fields related to a target forum url.
type Forum struct {
	ID                       string
	ScrapeDir                string
	Title                    string
	RootURL                  string
	ForumURL                 string
	PaginationClass          string
	PageActiveClass          string
	PageInactiveClass        string
	PageCount                int
	ThreadsTableSelector     string
	ThreadTitleSelector      string
	ThreadURLSelector        string
	ThreadLastReplySelector  string
	ThreadReplyCountSelector string
	ThreadOffset             int
	Threads                  []Thread
}

var forum Forum

// ForumPage TODO: compose Forum of []ForumPage instead of []Thread
type ForumPage struct {
	Number  int
	Threads []Thread
}

// Thread type stores fields related to a forum thread.
type Thread struct {
	ThreadID                 int
	ThreadTitle              string
	ThreadURL                string
	ThreadLastReplyTimestamp string
	ThreadReplyCount         int
	ThreadCreator            string
	ThreadPathComplete       string
	PageCount                int
	Posts                    []Post
}

// Post type represents each single post within a thread
type Post struct {
	ID       int
	Date     string
	Username string
	Name     string
	Location string
	Content  string
}

// Config holds global configuration properties
type Config struct {
	ConfigFileName string
	SubPath        string
	IndexPath      string
	ThreadPath     string
	PrintThreads   bool
	Maxthreads     int
	Throttle       int
}

// Conf is a global instance of Config and holds global configuration properties
var Conf Config

// Create index and thread scrape directories, setting their paths in Conf.
func setPaths(forum *Forum) {
	t := time.Now()
	Conf.SubPath = t.Format("2006-01-02-150405")
	Conf.IndexPath = filepath.Join(forum.ScrapeDir, Conf.SubPath, "index")
	Conf.ThreadPath = filepath.Join(forum.ScrapeDir, Conf.SubPath, "threads")
	err := os.MkdirAll(Conf.IndexPath, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory")
		log.Println(err)
		return
	}
	err = os.MkdirAll(Conf.ThreadPath, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory")
		log.Println(err)
		return
	}
}

// captureFlags sets flag default values, checks for CLI flags for '-config', '-maxthreads', '-throttle', '-printall' and set the related Conf fields.
func captureFlags() {
	flag.StringVar(&Conf.ConfigFileName, "config", "garden.json", "Name of configuration file")
	flag.IntVar(&Conf.Maxthreads, "maxthreads", 3, "(For dev env) When defined, sets the max number of threads scraped") // TODO: Set maxthreads default to 0.
	flag.IntVar(&Conf.Throttle, "throttle", 1000, "When defined, adds a throttle (in milliseconds) between page scraped")
	flag.BoolVar(&Conf.PrintThreads, "`printall`", false, "When true, logs out thread stats")
	flag.Parse()
	// fmt.Printf("Running with maxthreads: %d, throttle: %d, config: %s\n\n", Conf.Maxthreads, Conf.Throttle, Conf.ConfigFileName)
}

// loadConfig looks for a loads a required job configuration file based on the '-config' flag that was passed in at launch.
func loadConfig() {
	file, err := os.Open(filepath.Join("config", Conf.ConfigFileName))
	if err != nil {
		log.Fatalf("Failed to open file config\n")
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&forum)
	file.Close()
	checkError(err)
}

func logString() string {
	return fmt.Sprintf("Scrape job starting with maxthreads: %d, throttle: %d, config: %s, HTML dir: %s", Conf.Maxthreads, Conf.Throttle, Conf.ConfigFileName, Conf.SubPath)
}

func main() {
	tracelog.StartFile(tracelog.LevelTrace, "log", 100)

	captureFlags()
	loadConfig()
	setPaths(&forum)
	tracelog.Trace("BBSpider", forum.ID, logString())

	ForumScrape(&forum)
	ThreadScrape(&forum)

	tracelog.Stop()
}
