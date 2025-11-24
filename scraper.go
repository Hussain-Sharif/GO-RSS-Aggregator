package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

func startScrapping(
	db *database.Queries,
	concurrency int, // how many different goroutines we want to do scrapping on
	timeBetweenRequest time.Duration,
){
	log.Printf("Scrapping on %v goroutines every %s duration",concurrency,timeBetweenRequest)
	ticker:=time.NewTicker(timeBetweenRequest) // as we are using NewTicker we continue to tick around the specified time
	for ; ; <-ticker.C{ // the empty ; ; helps to execute the for loop immedieately then it waits as per the ticker.C interval 
		feeds,err:=db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)// Till here we are trying to get the Feeds which are in the DB and we are considering those limit of feeds which are actually new to the DB(1st priority) and the least recent   
		//  
		if err!=nil{
			log.Println("error while getting the saved feeds on DB:",err)
			continue
		}

		wg:=&sync.WaitGroup{}
		for _,feed:=range feeds{
			wg.Add(1)
			go scrapeFeed(db,wg,feed)
		}
		wg.Wait()

	}
}

// func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed){
// 	defer wg.Done() 

// 	_,err:=db.MarkFeedAsFetched(context.Background(),feed.ID)
// 	if err!=nil{
// 		log.Println("Error Occured when marking as feed fetched",err)
// 		return 
// 	}

// 	rssFeed,err:=urlToFeed(feed.Url)
// 	if err!=nil{
// 		log.Println("Error Occured when fetching feed",err)
// 		return 
// 	}

// 	for _,item:=range rssFeed.Channel.Item{
// 		description := sql.NullString{}
// 		if item.Description!=""{
// 			description.String=item.Description
// 			description.Valid=true
// 		}

// 		pub_at,err:=time.Parse(time.RFC1123Z,item.PubDate)

// 		if err!=nil{
// 			log.Printf("Couldn't able to parse the Published Date %v with error: %v",item.PubDate,err)
// 			continue
// 		}

// 		_,err=db.CreatePost(context.Background(),database.CreatePostParams{
// 			ID: uuid.New(),
// 			CreatedAt: time.Now().UTC(),
// 			UpdatedAt: time.Now().UTC(),
// 			Title: item.Title,
// 			Description: description,
// 			PublishedAt: pub_at,
// 			Url: item.Link,
// 			FeedID: feed.ID,
// 		})
// 		if(err!=nil){
// 			if strings.Contains(err.Error(),"duplicate key"){
// 				continue
// 			}
// 			log.Printf("Failed to Create a post with error: %v",err)
// 		}

// 	}
// 	log.Printf("Feed %s collected %v posts found",feed.Name,len(rssFeed.Channel.Item))


// }
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
    defer wg.Done()
    
    // Add retry logic with exponential backoff
    maxRetries := 3
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := attemptScrapeFeed(db, feed)
        if err == nil {
            return // Success!
        }
        
        if attempt < maxRetries-1 {
            backoff := time.Duration(attempt+1) * time.Second
            log.Printf("Retry %d for feed %s after %v", attempt+1, feed.Name, backoff)
            time.Sleep(backoff)
        } else {
            log.Printf("Failed to scrape feed %s after %d attempts: %v", 
                feed.Name, maxRetries, err)
        }
    }
}

func attemptScrapeFeed(db *database.Queries, feed database.Feed) error {
    // Mark as fetched FIRST (prevents hammering broken feeds)
    _, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
    if err != nil {
        return fmt.Errorf("marking feed: %w", err)
    }
    
    // Fetch with timeout
    client := http.Client{
        Timeout: 10 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if len(via) >= 5 {
                return fmt.Errorf("too many redirects")
            }
            return nil
        },
    }
    
    resp, err := client.Get(feed.Url)
    if err != nil {
        return fmt.Errorf("fetching: %w", err)
    }
    defer resp.Body.Close()
    
    // Check status code
    if resp.StatusCode != 200 {
        return fmt.Errorf("bad status: %d", resp.StatusCode)
    }
    
    // Limit response size (prevent memory attacks)
    const maxSize = 10 * 1024 * 1024 // 10MB
    limitedReader := io.LimitReader(resp.Body, maxSize)
    data, err := io.ReadAll(limitedReader)
    if err != nil {
        return fmt.Errorf("reading body: %w", err)
    }
    
    // Parse feed
    universalFeed, err := parseFeed(data)
    if err != nil {
        return fmt.Errorf("parsing: %w", err)
    }
    
    // Save posts with deduplication
    savedCount := 0
    for _, item := range universalFeed.Items {
        if item.Link == "" || item.Title == "" {
            continue // Skip invalid items
        }
        
        description := sql.NullString{}
        if item.Description != "" {
            description.String = item.Description
            description.Valid = true
        }
        
        _, err = db.CreatePost(context.Background(), database.CreatePostParams{
            ID:          uuid.New(),
            CreatedAt:   time.Now().UTC(),
            UpdatedAt:   time.Now().UTC(),
            Title:       item.Title,
            Description: description,
            PublishedAt: item.PubDate,
            Url:         item.Link,
            FeedID:      feed.ID,
        })
        
        if err == nil {
            savedCount++
        } else if !strings.Contains(err.Error(), "duplicate key") {
            log.Printf("Error saving post: %v", err)
        }
    }
    
    log.Printf("Feed %s: saved %d/%d posts", feed.Name, savedCount, len(universalFeed.Items))
    return nil
}
