# BBSpider

This purpose of this application is to manage a system of automated jobs and scheduled tasks that run against various 'bulletin board' style websites in order to scrape their content, organize the scraped content, and then provide desired stats to BBSpider users or other content analysis software.


## Changelog

*    **2017-05-15 (v0.0.1)** : Added scrape and capture of all thread pages.  Added logging system that writes to file with each job run.  (Needs improvement) 
*    **2017-05-14 (v0.0.1)** : Added throttle, maxthreads, config command-line parameters (all with defaults). Add config file for garden.com. Note: maxthreads default is set to 3.  Explicitly set to 0 in CLI for override (for all scrape threads/pages will be run)
*    **2017-05-13 (v0.0.1)** : Functioning display forum threads in console for 'All Things Gardening forum' forum on garden.com.
*    **2017-05-12 (v0.0.1)** : First commit.


## TODO Next
* Set up database to store indexed data from the scraped HTML for future analysis (including key phrases of interest).
* Add a sister process which scans the scraped HTML files and organizes content in a relational manner (users table, posts table, threads table, etc.).
* Consider pros and cons of indexing complete word counts associated with all users, posts, threads, or forum as this should help establish user post patterns, even with words that are not usually of interest.


## Pending Requirements

* Manage website targets by URL, pagination profile (pagination URL logic), thread logic (css selectors), post logic (css selectors), and username logic.
* Manage content scraping scheduled tasks that spider the task's target site looking for and storing any new content since the last scrape job ran.
* Manage content analysis scheduled tasks that scan scraped content and then organize post content, poster, timestamp, embedded URLs and post stats in a relational database based on criteria of interest provider by BBSpider users.
* Manage content analysis criteria in a manner that allows BBSpider users to update key phrases or URLs of interest.
* Provide desired stats and post material to BBSpider users with a ReactJS user interface that consumes a RESTful API.  The RESTful API will be designed with a conventional approach that allows for easy integration with other existing content analysis applications.
* Manage BBSpider users and administrators along with associated account parent organizations.


## Future Features for Consideration
* Add ability for Users to manage their own "tags" to associate to and categorize any thread or post of interest.  The tags can be added, edited or deleted, made "public" (for other users) or "private" (allowing use of tag only to tag creator).  When reviewing posts via the BBSpiderUI, allow the user to tag any content with any public or private tag.  Then, allow the user to filter content by their associated tags.
* Add ability for Users to enter keywords or phrases that get flagged for future data analysis.  Alternatively, run complete word count indexing in the db on all posts during HTML processing.  Create a UI approach to weight/sort posts and users by any word's count.


## Installation

To run the app locally, download the project, open terminal and navigate to the project directory.  

Then, install all dependencies:
$ `go get ./...`

Then, build the app:
$ `go build`

There are default values for each flag, but it’s intended to be run with flags as so a user or another program would launch it with a command like so:
$ `./bbspider -config=garden.json -throttle-1000 -printall=true -maxthreads=0`


## Usage

The following flags are used in the command line when launching the app:

* **-config** specifies a required JSON config file which will be unique for each forum.  This JSON object will change significantly as the application evolves.
* **-throttle** specifies the time in milliseconds to sleep between each HTTP request.
* **-printall** designates whether or not you want some basic thread stats to get printed to the console for each parsed thread.
* **-maxthreads** is another optional flag that, when used, limits the number of threads and/or pages the crawler scrapes which I found helpful for development purposes.  For production, this field will not be required.  Set to 0 to essentially disable the field.  It is set to 3 by default for development purposes.


## Config files

The **-config** JSON files reside in the 'config' directory.  The JSON content of each file should have configuration properties defined similarly to those defined in config/garden.json:

```
{
    "id": "gardening",
    "scrapeDir": "scrapes/garden/gardening",
    "rootURL": "https://garden.org",
    "forumURL": "/forums/view/gardening/",
    "paginationClass": "page_chunk",
    "pageActiveClass": "PageActive",
    "pageInactiveClass": "PageInactive",
    "threadsTableSelector": ".pretty-table tbody tr",
    "threadTitleSelector": "td:nth-child(2)",
    "threadURLSelector": "td a",
    "threadLastReplySelector": "td:nth-child(3)",
    "threadReplyCountSelector": "td:nth-child(4)",
    "threadOffset": 20
}
```

### Config file properties

A number of the config properties are based on CSS selectors as they are handled by Goquery, a Go library that brings jQuery-like features to the Go CSS Selector library Cascadia.  The following properties are required:

* **id** represents the unique identifier for the forum to be scraped.
* **scrapeDir** represents the location that BBSpider will save the scraped HTML.  Please note 'garden' parent of 'gardening', where 'garden' represents the website (parent) of the target forum.
* **rootURL** defines the website root.
* **forumURL** defines the location of the forum within the website root URL.
* **paginationClass** (not yet utilized) represents the CSS selector for the class within which forum pagination links reside.
* **pageActiveClass** (not yet utilized) represents the CSS selector for the class associated with the active pagination page link.
* **pageInactiveClass** represents the CSS selector for the class of pagination page links (that are inactive).  It's used to find the last page of a forum or thread to get the respective page count. (TODO: Replace this with a new pageInactiveSelector property)
* **threadsTableSelector** represents the CSS selector for the class of the top level of a forum index page where each row will represent a forum thread.
* **threadTitleSelector** represents the CSS selector for the class of the element where the thread title resides.
* **threadURLSelector** represents the CSS selector for the class of the element where the thread page 1 URL resides.
* **threadLastReplySelector** represents the CSS selector for the class of the element where the thread's last reply date resides.
* **threadReplyCountSelector** represents the CSS selector for the class of the element where the thread's reply count resides.
* **threadOffset** drives the pagination logic for how many threads (or posts) exist in a single page.


## Dependencies

The following libraries are required for the app and can be installed in your environment one at a time with a go get command (e.g. $ `go get github.com/PuerkitoBio/goquery`) or all together using $ `go get ./...`:

* "github.com/PuerkitoBio/goquery"
* "github.com/goinggo/tracelog"


## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
