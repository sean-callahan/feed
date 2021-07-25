## feed

feed is a feed generator library that can generate RSS feeds in Go

### Usage 

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sean-callahan/feed"
)

func main() {
	f := feed.Feed{
		Title:       "Example Feed",
		Link:        &feed.Link{URL: "https://github.com/sean-callahan/feed"},
		Description: "An example feed that is generated",
		Author:      &feed.Author{Name: "Sean Callahan"},
		Updated:     time.Now().UTC(),
		Items: []*feed.Item{
			{
				Title:       "Example Item",
				Link:        &feed.Link{URL: "https://github.com/sean-callahan/feed"},
				Description: "First and only item of the feed",
				Updated:     time.Now().UTC(),
			},
		},
	}

	rss, err := feed.RSS(f)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(rss)
}
