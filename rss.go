package feed

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	indent = "  "
)

// RSS generates an RSS feed for the provided feed with optional options.
func RSS(feed Feed, opts ...Options) (string, error) {
	rss, err := newRSS(feed)
	if err != nil {
		return "", err
	}

	// populate all the feed options
	for _, opt := range opts {
		if err := opt.populate(feed, rss); err != nil {
			return "", fmt.Errorf("%s: %v", opt.name(), err)
		}
	}

	var out []byte
	if rss.minimize {
		out, err = xml.Marshal(rss)
	} else {
		out, err = xml.MarshalIndent(rss, "", indent)
	}
	if err != nil {
		return "", fmt.Errorf("xml marshal: %v", err)
	}

	sb := new(strings.Builder)
	sb.WriteString(xml.Header)
	sb.Write(out)

	return sb.String(), nil
}

// newRSS creates a new RSS <rss> element for a Feed.
// Also creates the <channel> and <item> elements.
func newRSS(f Feed) (*rssXML, error) {
	if f.Updated.IsZero() {
		f.Updated = time.Now().UTC()
	}

	ch := &rssChannel{
		Title:          f.Title,
		Link:           f.Link.URL,
		Description:    f.Description,
		Language:       f.Language,
		Copyright:      f.Copyright,
		Generator:      f.Generator,
		ManagingEditor: formatRSSAuthor(f.Author),
		PubDate:        f.Updated.Format(time.RFC1123Z),
		LastBuildDate:  f.Updated.Format(time.RFC1123Z),
	}

	for _, item := range f.Items {
		i, err := newRSSItem(item)
		if err != nil {
			return nil, err
		}
		ch.Items = append(ch.Items, i)
	}

	return &rssXML{
		Version:   "2.0",
		ContentNS: "http://purl.org/rss/1.0/modules/content/",
		Channel:   ch,
	}, nil
}

// newRSSItem creates a new RSS <item> element from an Item.
func newRSSItem(i *Item) (*rssItem, error) {
	v := &rssItem{
		Title:       i.Title,
		Link:        i.Link.URL,
		Description: i.Description,
		GUID:        i.ID,
		PubDate:     i.Updated.Format(time.RFC1123Z),
	}
	if i.Author != nil {
		v.Author = formatRSSAuthor(i.Author)
	}
	if i.Enclosure != nil {
		v.Enclosure = &rssEnclosure{
			URL:    i.Enclosure.URL,
			Length: i.Enclosure.Length,
			Type:   i.Enclosure.Type,
		}
	}
	return v, nil
}

// formatRSSAuthor formats an author for an RSS feed.
// If both the Author's Email and Name are provided it is
// formated as: "jappleseed@example.com (Johnny Appleseed)",
// Otherwise only the email address is returned.
func formatRSSAuthor(a *Author) string {
	if a == nil || a.Email == "" {
		return ""
	}
	if a.Name == "" {
		return a.Email
	}
	return fmt.Sprintf("%s (%s)", a.Email, a.Name)
}

type rssXML struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	ContentNS    string   `xml:"xmlns:content,attr"`
	ItunesNS     string   `xml:"xmlns:itunes,attr,omitempty"`
	GooglePlayNS string   `xml:"xmlns:googleplay,attr,omitempty"`
	PodcastNS    string   `xml:"xmlns:podcast,attr,omitempty"`
	Channel      *rssChannel

	minimize bool
}

type rssChannel struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`              // required
	Link           string   `xml:"link"`               // required
	Description    string   `xml:"description"`        // required
	Language       string   `xml:"language,omitempty"` // itunes required
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"`
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"`
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	TTL            int      `xml:"ttl,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	Image          *rssImage
	TextInput      *rssTextInput

	rssItunes
	rssGooglePlay
	rssPodcastIndex

	Items []*rssItem `xml:"item"`
}

type rssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`       // required
	Link        string   `xml:"link"`        // required
	Description string   `xml:"description"` // required
	Content     *rssContent
	Author      string `xml:"author,omitempty"`
	Category    string `xml:"category,omitempty"`
	Comments    string `xml:"comments,omitempty"`
	Enclosure   *rssEnclosure
	GUID        string `xml:"guid,omitempty"`
	PubDate     string `xml:"pubDate,omitempty"`
	Source      string `xml:"source,omitempty"`

	rssItunesItem
}

type rssImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

type rssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

type rssContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

type rssEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type rssItunes struct {
	ItunesImage      *itunesImage      `xml:",omitempty"`                // itunes required
	ItunesCategories []*itunesCategory `xml:",omitempty"`                // itunes required
	ItunesExplicit   string            `xml:"itunes:explicit,omitempty"` // itunes required
	ItunesAuthor     string            `xml:"itunes:author,omitempty"`
	ItunesOwner      *itunesOwner      `xml:"itunes:owner,omitempty"`
	ItunesTitle      string            `xml:"itunes:title,omitempty"`
	ItunesType       string            `xml:"itunes:type,omitempty"`
	ItunesComplete   string            `xml:"itunes:complete,omitempty"`
}

type rssItunesItem struct {
	ItunesImage    *itunesImage `xml:",omitempty"`
	ItunesDuration string       `xml:"itunes:duration,omitempty"`
	ItunesExplicit string       `xml:"itunes:explicit,omitempty"`
}

type itunesImage struct {
	XMLName xml.Name `xml:"itunes:image"`
	Href    string   `xml:"href,attr"`
}

type itunesCategory struct {
	XMLName xml.Name        `xml:"itunes:category"`
	Text    string          `xml:"text,attr"`
	Sub     *itunesCategory `xml:",omitempty"`
}

type itunesOwner struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Email   string   `xml:"itunes:email,omitempty"`
	Name    string   `xml:"itunes:name,omitempty"`
}

type rssGooglePlay struct {
	GooglePlayCategories []*googlePlayCategory `xml:",omitempty"`
}

type googlePlayCategory struct {
	XMLName xml.Name `xml:"googleplay:category"`
	Text    string   `xml:"text,attr"`
}

type rssPodcastIndex struct {
	PodcastFunding *rssPodcastFunding `xml:",omitempty"`
}

type rssPodcastFunding struct {
	XMLName string `xml:"podcast:funding"`
	URL     string `xml:"url,attr"`
	Text    string `xml:",innerxml"`
}

// formatDuration formats a time.Duration in the format hh:mm:ss.
func formatDuration(d time.Duration) string {
	if d < 0 {
		return ""
	}
	d = d.Round(time.Second)
	hh := d / time.Hour
	d %= time.Hour
	mm := d / time.Minute
	d %= time.Minute
	ss := d / time.Second
	return fmt.Sprintf("%d:%02d:%02d", hh, mm, ss)
}

// ApplePodcast describes an Apple Podcasts (iTunes)-supported RSS feed.
type ApplePodcast struct {
	// Categories defines a list of (possibly nested) categories for the Apple Podcast directory.
	Categories []Category

	// Type of show: episodic or serial.
	Type string

	// Complete should be true if the podcast will never have more episodes.
	Complete bool
}

func (ApplePodcast) name() string { return "ApplePodcast" }

func (ap ApplePodcast) populate(f Feed, rss *rssXML) error {
	if rss == nil || rss.Channel == nil {
		return nil
	}

	rss.ItunesNS = "http://www.itunes.com/dtds/podcast-1.0.dtd"

	if f.Image != nil {
		rss.Channel.ItunesImage = &itunesImage{Href: f.Image.URL}
	}
	for _, cat := range ap.Categories {
		rss.Channel.ItunesCategories = append(rss.Channel.ItunesCategories, ap.newCategory(&cat))
	}
	rss.Channel.ItunesExplicit = strconv.FormatBool(f.Explicit)

	if f.Author != nil {
		rss.Channel.ItunesAuthor = f.Author.Name
	}

	if f.Owner != nil {
		rss.Channel.ItunesOwner = &itunesOwner{
			Email: f.Owner.Email,
			Name:  f.Owner.Name,
		}
	}

	for i, item := range rss.Channel.Items {
		src := f.Items[i]
		if src == nil || src.Enclosure == nil {
			return fmt.Errorf("item at %d: must contain an enclosure", i)
		}

		if d := src.Enclosure.Duration; d > 0 {
			item.ItunesDuration = formatDuration(d)
		}
		if v := src.Image; v != nil {
			item.ItunesImage = &itunesImage{
				Href: v.URL,
			}
		}
		if v := src.Explicit; v {
			item.ItunesExplicit = "true"
		}
	}

	rss.Channel.ItunesType = ap.Type

	if ap.Complete {
		rss.Channel.ItunesComplete = "Yes"
	}

	return nil
}

// newCategory creates a new <itunes:category> element. Possibily containing sub-elements.
func (ap ApplePodcast) newCategory(c *Category) *itunesCategory {
	if c == nil {
		return nil
	}
	return &itunesCategory{
		Text: c.Name,
		Sub:  ap.newCategory(c.Sub),
	}
}

// MinimizeOutput sets whether or not the generated feed is minified.
// Setting this value to true will not indent the generated feed.
type MinimizeOutput bool

func (MinimizeOutput) name() string { return "MinimizeOutput" }

func (m MinimizeOutput) populate(f Feed, rss *rssXML) error {
	rss.minimize = bool(m)
	return nil
}

// GooglePlay describes a Google Play podcast supported RSS feed.
type GooglePlay struct {
	Categories []string
}

func (GooglePlay) name() string { return "GooglePlay" }

func (gp GooglePlay) populate(f Feed, rss *rssXML) error {
	if rss == nil || rss.Channel == nil {
		return nil
	}

	rss.GooglePlayNS = "http://www.google.com/schemas/play-podcasts/1.0"

	for _, cat := range gp.Categories {
		v := &googlePlayCategory{Text: cat}
		rss.Channel.GooglePlayCategories = append(rss.Channel.GooglePlayCategories, v)
	}

	return nil
}

type PodcastIndex struct {
	Funding *Link
}

func (PodcastIndex) name() string { return "PodcastIndex" }

func (p PodcastIndex) populate(f Feed, rss *rssXML) error {
	if rss == nil || rss.Channel == nil {
		return nil
	}

	rss.PodcastNS = "https://podcastindex.org/namespace/1.0"

	if p.Funding != nil {
		rss.Channel.PodcastFunding = &rssPodcastFunding{
			URL:  p.Funding.URL,
			Text: p.Funding.Text,
		}
	}

	return nil
}
