package main

import (
	"github.com/mmcdole/gofeed"
	"github.com/gorilla/feeds"
	"regexp"
	"fmt"
	"flag"
	"log"
	"time"
	"os"
)

type gorillaFeed = feeds.Feed
type gorillaImage = feeds.Image
type gorillaItem = feeds.Item
type gorillaLink = feeds.Link
type gorillaAuthor = feeds.Author
type gorillaEnclosure = feeds.Enclosure

// var log = logging.MustGetLogger("filter-rss")
// var format = logging.MustStringFormatter(
// 	`%{color}%{time:15:04:05.000} %{shortfunc} %{level} %{color:reset} %{message}`,
// )
var logger = log.Default()

func populate_gorillaItem(now time.Time, item *gofeed.Item) gorillaItem {
	log.Printf("Populating %v", item.Title)

	names, emails := "", ""
	if len(item.Authors) != 0 {
		log.Printf("\tGetting Names and Emails")
		names, emails = func(people []*gofeed.Person) (string, string) {
			names, emails := "", ""
			for i, person := range people {
				names += person.Name
				emails += person.Email
				if i < len(people)-1 {
					names += " "
					emails += " "
				}
			}
			return names, emails
		}(item.Authors)
	} else {
		log.Printf("\tNo Names and Emails")
	}

	log.Printf("\tnames: %v", names)
	log.Printf("\temails: %v", emails)
	enclosure_url, enclosure_length, enclosure_type := "", "", ""
	if len(item.Enclosures) != 0 {
		log.Printf("\tGetting Enclosures")
		enclosure_url, enclosure_length, enclosure_type = func(enclosures []*gofeed.Enclosure) (string, string, string) {
			enclosure_url, enclosure_length, enclosure_type := "", "", ""
			for i, e := range enclosures {
				enclosure_url += e.URL
				enclosure_length += e.Length
				enclosure_type += e.Type
				if i < len(enclosures)-1 {
					enclosure_url += " "
					enclosure_length += " "
					enclosure_type += " "
				}
			}
			return enclosure_url, enclosure_length, enclosure_type
		}(item.Enclosures)
	} else {
		log.Printf("\tNo Enclosures")
	}

	log.Printf("\tenclosure_url: %v", enclosure_url)
	log.Printf("\tenclosure_length: %v", enclosure_length)
	log.Printf("\tenclosure_type: %v", enclosure_type)

	log.Printf("Loading gorillaItem")
	log.Printf("\tTitle: %v", item.Title)
	log.Printf("\tLink: %v", item.Link)
	log.Printf("\tAuthor: %v", item.Authors)
	log.Printf("\tDescription: %v", item.Description)
	log.Printf("\tId: %v", item.GUID)
	log.Printf("\tUpdated: %v", item.Published)
	log.Printf("\tCreated: %v", now)
	log.Printf("\tEnclosure: %v", item.Enclosures)
	log.Printf("\tContent: %v", item.Content)

	var updated time.Time
	if item.UpdatedParsed != nil {
		updated = *(item.UpdatedParsed)
	}

	return gorillaItem{
		Title: item.Title,
		Link: &gorillaLink{Href: item.Link},
		Author: &gorillaAuthor{Name: names, Email: emails},
		Description: item.Description,
		Id: item.GUID,
		Updated: updated,
		Created: now,
		Enclosure: &gorillaEnclosure{Url: enclosure_url, Length: enclosure_length, Type: enclosure_type},
		Content: item.Content,
	}
}

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "%s [url]\n", os.Args[0])
		flag.PrintDefaults()
	}

	pattern := flag.String("pattern", "", "a regex pattern")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// logging.SetLevel(logging.WARNING, "filter-rss")
	// logging.SetFormatter(format)
	// verbosity
	if *verbose {
		// logging.SetLevel(logging.DEBUG, "filter-rss")
		logger.SetPrefix("DEBUG ")
		logger.SetOutput(os.Stdout)
	} else {
		null, _ := os.Open(os.DevNull)
		logger.SetOutput(null)
	}

	// log.Debugf("flag_verbose: %v", *verbose)
	// log.Debugf("flag_pattern: %v", *pattern)
	log.Printf("flag_verbose: %v", *verbose)
	log.Printf("flag_pattern: %v", *pattern)

	urls := flag.Args()
	// log.Debugf("args: %v", urls)
	log.Printf("args: %v", *pattern)

	fp := gofeed.NewParser()
	go_feed, err := fp.ParseURL(urls[0])
	if err != nil {
		// log.Errorf("%v", err)
		// os.Exit(1)
		log.Fatal(err)
	}

	now := time.Now()
	var matchedItems []*gorillaItem
	for _, item := range go_feed.Items {
		for _, category := range item.Categories {
			matched, err := regexp.MatchString(*pattern, category)
			if err != nil {
				// log.Errorf("%v", err)
				// os.Exit(1)
				log.Fatal(err)
			}
			if matched {
				// log.Debugf("%v", item.Title)
				log.Printf("Found %v", item.Title)
				matchedItem := populate_gorillaItem(now, item)
				matchedItems = append(matchedItems, &matchedItem)
				break
			}
		}
	}

	names, emails := "", ""
	if len(go_feed.Authors) != 0 {
		names, emails = func(people []*gofeed.Person) (string, string) {
			names, emails := "", ""
			for i, person := range people {
				names += person.Name
				emails += person.Email
				if i < len(people)-1 {
					names += " "
					emails += " "
				}
			}
			return names, emails
		}(go_feed.Authors)
	}

	log.Printf("Creating Feed")
	log.Printf("\tTitle: %v", go_feed.Title)
	log.Printf("\tLink: %v", go_feed.Link)
	log.Printf("\tDescription: %v", go_feed.Description)
	log.Printf("\tAuthor: %v", go_feed.Authors)
	log.Printf("\tUpdated: %v", go_feed.Published)
	log.Printf("\tCreated: %v", now)
	log.Printf("\tCopyright: %v", go_feed.Copyright)
	log.Printf("\tImage: %v", go_feed.Image)

	var updated time.Time
	if go_feed.UpdatedParsed != nil {
		updated = *(go_feed.UpdatedParsed)
	}
	var image gorillaImage
	if go_feed.Image != nil {
		image = gorillaImage{
			Url: go_feed.Image.URL,
			Title: go_feed.Image.Title,
		}
	}

	matchedFeed := &gorillaFeed {
			Title: go_feed.Title,
			Link: &gorillaLink{Href: go_feed.Link},
			Description: go_feed.Description,
			Author: &gorillaAuthor{Name: names, Email: emails},
			Updated: updated,
			Created: now,
			Copyright: go_feed.Copyright,
			Image: &image,
		}

	log.Printf("gorillaItems are being loaded")
	matchedFeed.Items = matchedItems

	log.Printf("RSS conversion")
	rss, err := matchedFeed.ToRss()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(rss)
}
