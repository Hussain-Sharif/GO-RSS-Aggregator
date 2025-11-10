package main

import (
	"context"
	"database/sql"
	"log"
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

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed){
	defer wg.Done() 

	_,err:=db.MarkFeedAsFetched(context.Background(),feed.ID)
	if err!=nil{
		log.Println("Error Occured when marking as feed fetched",err)
		return 
	}

	rssFeed,err:=urlToFeed(feed.Url)
	if err!=nil{
		log.Println("Error Occured when fetching feed",err)
		return 
	}

	for _,item:=range rssFeed.Channel.Item{
		description := sql.NullString{}
		if item.Description!=""{
			description.String=item.Description
			description.Valid=true
		}

		pub_at,err:=time.Parse(time.RFC1123Z,item.PubDate)

		if err!=nil{
			log.Printf("Couldn't able to parse the Published Date %v with error: %v",item.PubDate,err)
			continue
		}

		_,err=db.CreatePost(context.Background(),database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Description: description,
			PublishedAt: pub_at,
			Url: item.Link,
			FeedID: feed.ID,
		})
		if(err!=nil){
			if strings.Contains(err.Error(),"duplicate key"){
				continue
			}
			log.Printf("Failed to Create a post with error: %v",err)
		}

	}
	log.Printf("Feed %s collected %v posts found",feed.Name,len(rssFeed.Channel.Item))


}