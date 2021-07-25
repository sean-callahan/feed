package feed

import "time"

type Feed struct {
	ID          string
	Link        *Link
	Created     time.Time
	Updated     time.Time
	Title       string
	Subtitle    string
	Description string
	Language    string
	Copyright   string

	Explicit bool
	Author   *Author
	Image    *Image
	Items    []*Item
}

type Author struct {
	Name  string
	Email string
}

type Image struct {
	URL    string
	Title  string
	Link   string
	Width  int
	Height int
}

type Link struct {
	URL string
}

type Category struct {
	Name string
	Sub  *Category
}

type Item struct {
	ID          string
	Link        *Link
	Created     time.Time
	Updated     time.Time
	Title       string
	Description string
	Author      *Author
	Enclosure   *Enclosure
	Duration    time.Duration
}

type Enclosure struct {
	URL    string
	Length string
	Type   string
}

// Options provide additional options when generating feeds.
type Options interface {
	name() string
	populate(f Feed, rss *rssXML) error
}
