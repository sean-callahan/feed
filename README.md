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
```

#### Output

```
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>Example Feed</title>
    <link>https://github.com/sean-callahan/feed</link>
    <description>An example feed that is generated</description>
    <pubDate>Sun, 25 Jul 2021 04:21:30 +0000</pubDate>
    <lastBuildDate>Sun, 25 Jul 2021 04:21:30 +0000</lastBuildDate>
    <item>
      <title>Example Item</title>
      <link>https://github.com/sean-callahan/feed</link>
      <description>First and only item of the feed</description>
      <pubDate>Sun, 25 Jul 2021 04:21:30 +0000</pubDate>
    </item>
  </channel>
</rss>
```