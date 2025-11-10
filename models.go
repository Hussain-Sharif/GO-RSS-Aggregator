package main

import (
	"time"

	"github.com/Hussain-Sharif/GO-RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url    string    `json:"url"`
	UserID  uuid.UUID 	`json:"user_id"`
}

func databaseFeedToFeed(dbFeed database.Feed) Feed{
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:    dbFeed.Url,
		UserID: dbFeed.UserID,
	}
}

func databaseFeedsToFeeds(allDBFeeds []database.Feed) []Feed{
	feeds:= []Feed{}
	for _,dbFeed:=range allDBFeeds{
		feeds=append(feeds, databaseFeedToFeed(dbFeed))
	}
	return feeds
}


type Feed_Follows struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID  uuid.UUID 	`json:"user_id"`
	FeedID uuid.UUID `json:"feed_id"`
}

func databaseFeedFollowsToFeedFollows(feedFollow database.FeedFollow)Feed_Follows {
	return Feed_Follows{
		ID: feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID: feedFollow.UserID,
		FeedID: feedFollow.FeedID,
	}
}

func databaseGetFeedFollowsToGetFeedFollows(getAllFeedFollows []database.FeedFollow) []Feed_Follows{
	allFeedFollows:=[]Feed_Follows{}
	for _,eachFeedFollow :=range getAllFeedFollows{
		allFeedFollows=append(allFeedFollows, databaseFeedFollowsToFeedFollows(eachFeedFollow))
	} 
	return allFeedFollows
}


type Post struct{
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string `json:"title"`
	Description *string `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databasePostToPost(dbPost database.Post) Post{
	var description *string 
	if dbPost.Description.Valid{
		description=&dbPost.Description.String
	}

	return Post{
		ID 	: dbPost.ID,
		CreatedAt: dbPost.CreatedAt,
		UpdatedAt: dbPost.UpdatedAt,
		Title: dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url: dbPost.Url,
		FeedID: dbPost.FeedID,
	}
}


func databasePostsTOPosts(dbPosts []database.Post) []Post{
	posts:=[]Post{}

	for _,eachPost := range dbPosts{
		posts=append(posts, databasePostToPost(eachPost))
	}

	return  posts

}