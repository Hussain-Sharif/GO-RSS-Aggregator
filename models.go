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

func databaseFeedFollowsToFeedFollows(allFeedFollows database.FeedFollow)Feed_Follows {
	return Feed_Follows{
		ID: allFeedFollows.ID,
		CreatedAt: allFeedFollows.CreatedAt,
		UpdatedAt: allFeedFollows.UpdatedAt,
		UserID: allFeedFollows.UserID,
		FeedID: allFeedFollows.FeedID,
	}
}
