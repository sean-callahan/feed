package feed

import (
	"fmt"
	"time"
)

// Feed represents a collection of items.
// In an RSS context, this is analogous to a <channel> element.
type Feed struct {
	// ID is a unique identifier for the feed.
	// This is commonly either a valid UUIDv4 or any unique string.
	ID string

	// Link holds a URL to the website that corresponds to the feed.
	Link *Link

	// Created is the timestamp when the feed was first created.
	Created time.Time

	// Updated is the timestamp of the last time the feed was generated.
	Updated time.Time

	// Title is the name of the feed.
	Title string

	// Subtitle is an optional subtitle for the feed.
	Subtitle string

	// Description is a short summary describing the feed.
	Description string

	// Langauge is the language the feed is written in.
	// This must be a valid RFC1766 language code.
	Language string

	// Copyright notice for the feed.
	Copyright string

	// Generator is the program that generated the feed.
	Generator string

	// Explicit is true when the feed may have explicit content.
	Explicit bool

	// Author is the name and email of the primary author of the feed.
	Author *Author

	// Owner is the name and/or company name of the owner of the feed.
	Owner *Author

	// Image is the image that is displayed with the feed.
	Image *Image

	// Items are all the items belonging to the feed.
	Items []*Item
}

// Author represents a feed or item author.
type Author struct {
	// Name is the full name of the author.
	Name string

	// Email is the primary email for the author.
	Email string
}

// String representation of an Author.
func (a Author) String() string {
	if a.Email == "" {
		return a.Name
	}
	if a.Name == "" {
		return a.Email
	}
	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}

// Image contains a link to an image that is displayed with a feed or item.
type Image struct {
	// URL to the image.
	URL string

	// Title (or caption) is a description of the image.
	Title string

	// Link is a URL to a website.
	Link string

	// Width is the width of the image in pixels.
	Width int

	// Height is the height of the image in pixels.
	Height int
}

// Link is a URL to webpage.
type Link struct {
	// URL to the webpage.
	URL string

	// Text describing the link.
	Text string
}

// String representation of the Link.
func (l Link) String() string {
	if l.URL == "" {
		return ""
	}
	if l.Text == "" {
		return l.URL
	}
	return fmt.Sprintf("%s (%s)", l.Text, l.URL)
}

// Category identifies a categorization taxonomy.
type Category struct {
	// Name of the category.
	Name string

	// Sub is an optional nested category.
	Sub *Category
}

// Item is an element of a feed.
type Item struct {
	// ID is a unique identifier for the item.
	ID string

	// Link to the webpage describing the item.
	Link *Link

	// Created is the time the item was first published.
	Created time.Time

	// Updated is the time the item has last been modified.
	Updated time.Time

	// Title is the name of the item.
	Title string

	// Description is a long summary of the item and its content.
	Description string

	// Image is the image that is displayed with the item.
	Image *Image

	// Author is the primary author of the item.
	Author *Author

	// Enclosure is the content contained with this item.
	Enclosure *Enclosure

	// Explicit is true when the content of the item is explicit.
	Explicit bool
}

// Enclosure encloses a file in an Item.
type Enclosure struct {
	// URL to a file.
	URL string

	// Length is a length of the file in bytes.
	Length string

	// Type is the MIME type of the enclosed file.
	Type string

	// Duration is how long the media file is.
	Duration time.Duration
}

// Options provide additional options when generating feeds.
type Options interface {
	name() string
	populate(f Feed, rss *rssXML) error
}
