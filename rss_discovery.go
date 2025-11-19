package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type FeedDiscoveryResult struct {
    OriginalURL string   `json:"original_url"`
    FeedURLs    []FeedURL `json:"feed_urls"`
    Error       string   `json:"error,omitempty"`
}

type FeedURL struct {
    URL   string `json:"url"`
    Type  string `json:"type"`   // "rss", "atom"
    Title string `json:"title"`
}

// Discovers RSS/Atom feeds from a webpage
func discoverFeeds(pageURL string) (FeedDiscoveryResult, error) {
    result := FeedDiscoveryResult{
        OriginalURL: pageURL,
        FeedURLs:    []FeedURL{},
    }
    
    // 1. Try common RSS URL patterns first
    commonPatterns := []string{
        "/feed",
        "/rss",
        "/feed.xml",
        "/rss.xml",
        "/index.xml",
        "/atom.xml",
        "/feed/",
        "/rss/",
    }
    
    baseURL := strings.TrimSuffix(pageURL, "/")
    
    for _, pattern := range commonPatterns {
        testURL := baseURL + pattern
        if isValidFeed(testURL) {
            feedType := getFeedType(testURL)
            result.FeedURLs = append(result.FeedURLs, FeedURL{
                URL:   testURL,
                Type:  feedType,
                Title: fmt.Sprintf("%s feed", feedType),
            })
        }
    }
    
    // 2. If patterns don't work, parse HTML for auto-discovery
    if len(result.FeedURLs) == 0 {
        htmlFeeds, err := parseHTMLForFeeds(pageURL)
        if err == nil {
            result.FeedURLs = append(result.FeedURLs, htmlFeeds...)
        }
    }
    
    if len(result.FeedURLs) == 0 {
        result.Error = "No RSS/Atom feeds found"
    }
    
    return result, nil
}

// Checks if URL is a valid RSS/Atom feed
func isValidFeed(url string) bool {
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil || resp.StatusCode != 200 {
        return false
    }
    defer resp.Body.Close()
    
    // Read first 512 bytes to check if it's XML
    buffer := make([]byte, 512)
    n, _ := resp.Body.Read(buffer)
    content := string(buffer[:n])
    
    // Check for RSS/Atom markers
    return strings.Contains(content, "<rss") || 
           strings.Contains(content, "<feed") ||
           strings.Contains(content, "<?xml")
}

// Determines feed type (RSS vs Atom)
func getFeedType(url string) string {
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return "unknown"
    }
    defer resp.Body.Close()
    
    buffer := make([]byte, 1024)
    n, _ := resp.Body.Read(buffer)
    content := string(buffer[:n])
    
    if strings.Contains(content, "<rss") {
        return "rss"
    } else if strings.Contains(content, "<feed") {
        return "atom"
    }
    return "unknown"
}

// Parses HTML for RSS auto-discovery tags
func parseHTMLForFeeds(pageURL string) ([]FeedURL, error) {
    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(pageURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    doc, err := html.Parse(resp.Body)
    if err != nil {
        return nil, err
    }
    
    var feeds []FeedURL
    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "link" {
            var rel, feedType, href, title string
            for _, attr := range n.Attr {
                switch attr.Key {
                case "rel":
                    rel = attr.Val
                case "type":
                    feedType = attr.Val
                case "href":
                    href = attr.Val
                case "title":
                    title = attr.Val
                }
            }
            
            // Check if it's an RSS/Atom feed link
            if rel == "alternate" && 
               (strings.Contains(feedType, "rss") || 
                strings.Contains(feedType, "atom")) {
                
                // Convert relative URLs to absolute
                if !strings.HasPrefix(href, "http") {
                    href = resolveURL(pageURL, href)
                }
                
                feedFormat := "rss"
                if strings.Contains(feedType, "atom") {
                    feedFormat = "atom"
                }
                
                feeds = append(feeds, FeedURL{
                    URL:   href,
                    Type:  feedFormat,
                    Title: title,
                })
            }
        }
        
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)
    
    return feeds, nil
}

// Resolves relative URLs to absolute
func resolveURL(baseURL, relativeURL string) string {
    if strings.HasPrefix(relativeURL, "http") {
        return relativeURL
    }
    
    base := strings.TrimSuffix(baseURL, "/")
    if strings.HasPrefix(relativeURL, "/") {
        return base + relativeURL
    }
    return base + "/" + relativeURL
}
