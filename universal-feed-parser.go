package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"
)

// Universal feed structure that handles both RSS and Atom
type UniversalFeed struct {
    Title   string
    Link    string
    Items   []UniversalItem
}

type UniversalItem struct {
    Title       string
    Link        string
    Description string
    Content     string
    PubDate     time.Time
    GUID        string
}

// RSS 2.0 Structure
type RSS struct {
    XMLName xml.Name `xml:"rss"`
    Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
    Title       string    `xml:"title"`
    Link        string    `xml:"link"`
    Description string    `xml:"description"`
    Items       []RSSItem `xml:"item"`
}

type RSSItem struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
    GUID        string `xml:"guid"`
}

// Atom 1.0 Structure
type Atom struct {
    XMLName xml.Name `xml:"feed"`
    Title   string   `xml:"title"`
    Link    []AtomLink `xml:"link"`
    Entries []AtomEntry `xml:"entry"`
}

type AtomLink struct {
    Href string `xml:"href,attr"`
    Rel  string `xml:"rel,attr"`
}

type AtomEntry struct {
    Title   string     `xml:"title"`
    Link    []AtomLink `xml:"link"`
    Summary string     `xml:"summary"`
    Content AtomContent `xml:"content"`
    Updated string     `xml:"updated"`
    ID      string     `xml:"id"`
}

type AtomContent struct {
    Type string `xml:"type,attr"`
    Body string `xml:",chardata"`
}

// Parse any feed format into universal structure
func parseFeed(data []byte) (UniversalFeed, error) {
    // Try RSS first
    var rss RSS
    if err := xml.Unmarshal(data, &rss); err == nil && rss.Channel.Title != "" {
        return rssToUniversal(rss), nil
    }
    
    // Try Atom
    var atom Atom
    if err := xml.Unmarshal(data, &atom); err == nil && atom.Title != "" {
        return atomToUniversal(atom), nil
    }
    
    return UniversalFeed{}, fmt.Errorf("unsupported feed format")
}

func rssToUniversal(rss RSS) UniversalFeed {
    feed := UniversalFeed{
        Title: rss.Channel.Title,
        Link:  rss.Channel.Link,
        Items: make([]UniversalItem, 0, len(rss.Channel.Items)),
    }
    
    for _, item := range rss.Channel.Items {
        pubDate := parseFlexibleDate(item.PubDate)
        
        feed.Items = append(feed.Items, UniversalItem{
            Title:       item.Title,
            Link:        item.Link,
            Description: item.Description,
            PubDate:     pubDate,
            GUID:        item.GUID,
        })
    }
    
    return feed
}

func atomToUniversal(atom Atom) UniversalFeed {
    feed := UniversalFeed{
        Title: atom.Title,
        Items: make([]UniversalItem, 0, len(atom.Entries)),
    }
    
    // Find feed link
    for _, link := range atom.Link {
        if link.Rel == "alternate" || link.Rel == "" {
            feed.Link = link.Href
            break
        }
    }
    
    for _, entry := range atom.Entries {
        pubDate := parseFlexibleDate(entry.Updated)
        
        // Find entry link
        var entryLink string
        for _, link := range entry.Link {
            if link.Rel == "alternate" || link.Rel == "" {
                entryLink = link.Href
                break
            }
        }
        
        content := entry.Summary
        if entry.Content.Body != "" {
            content = entry.Content.Body
        }
        
        feed.Items = append(feed.Items, UniversalItem{
            Title:       entry.Title,
            Link:        entryLink,
            Description: entry.Summary,
            Content:     content,
            PubDate:     pubDate,
            GUID:        entry.ID,
        })
    }
    
    return feed
}

// Flexible date parser that handles multiple formats
func parseFlexibleDate(dateStr string) time.Time {
    if dateStr == "" {
        return time.Now()
    }
    
    // Try common formats
    formats := []string{
        time.RFC1123Z,      // RSS: "Mon, 02 Jan 2006 15:04:05 -0700"
        time.RFC1123,       // RSS without timezone
        time.RFC3339,       // Atom: "2006-01-02T15:04:05Z07:00"
        time.RFC3339Nano,   // Atom with nanoseconds
        "2006-01-02 15:04:05",
        "2006-01-02",
    }
    
    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t
        }
    }
    
    // If all fail, return current time
    log.Printf("Couldn't parse date: %s", dateStr)
    return time.Now()
}
